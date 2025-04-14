package kafka

import (
	"log/slog"

	"github.com/IBM/sarama"

	"github.com/xoticdsign/book/internal/lib/logger"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

// Broker{} represents a Kafka message broker that handles producing booking events to a Kafka topic.
type Broker struct {
	Producer sarama.SyncProducer
	Log      *logger.Logger

	config utils.Config
}

// New() initializes and returns a new Kafka Broker with the given configuration.
//
// It creates a new synchronous Kafka producer using Sarama.
func New(cfg utils.Config, log *logger.Logger) (*Broker, error) {
	saramaCfg := sarama.NewConfig()

	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.KafkaConfig.Addresses, saramaCfg)
	if err != nil {
		return &Broker{}, err
	}

	return &Broker{
		Producer: producer,
		Log:      log,

		config: cfg,
	}, nil
}

// Shutdown() gracefully closes the Kafka producer connection.
func (b *Broker) Shutdown() {
	b.Producer.Close()
}

// BookNotifyEvent{} represents the data structure of a booking event that will be published to the Kafka topic.
type BookNotifyEvent struct {
	Ticket string
	Data   *bookrpc.BookRequest
}

// BookNotify() sends a BookNotifyEvent to the configured Kafka topic.
//
// It serializes the event to JSON and logs success or failure.
func (b *Broker) BookNotify(event *BookNotifyEvent) error {
	const op = "BookNotify()"

	partition, offset, err := b.Producer.SendMessage(&sarama.ProducerMessage{
		Topic:     b.config.KafkaConfig.Topic,
		Value:     sarama.ByteEncoder(utils.MarshalJSON(event)),
		Offset:    b.config.KafkaConfig.Offset,
		Partition: b.config.KafkaConfig.Partition,
	})
	if err != nil {
		b.Log.Logs.BrokerLog.Error(
			"can't produce a message",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)

		return err
	}
	b.Log.Logs.BookLog.Debug(
		"message produced",
		slog.String("op", op),
		slog.Any("partition", partition),
		slog.Any("offset", offset),
		slog.Any("event", event),
	)

	return nil
}

// UnimplementedBroker{} is a stub that implements the Broker interface
//
// but does nothing. Useful for testing or placeholder functionality.
type UnimplementedBroker struct{}

// BookNotify() is the no-op implementation for the BookNotify method.
func (u *UnimplementedBroker) BookNotify(event *BookNotifyEvent) error { return nil }

// Shutdown() is the no-op implementation for the Shutdown method.
func (u *UnimplementedBroker) Shutdown() {}
