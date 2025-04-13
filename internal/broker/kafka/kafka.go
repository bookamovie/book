package kafka

import (
	"log/slog"

	"github.com/IBM/sarama"

	"github.com/xoticdsign/book/internal/lib/logger"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

type Broker struct {
	Producer sarama.SyncProducer

	log    *logger.Logger
	config utils.Config
}

func New(cfg utils.Config, log *logger.Logger) (*Broker, error) {
	saramaCfg := sarama.NewConfig()

	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.KafkaConfig.Addresses, saramaCfg)
	if err != nil {
		return &Broker{}, err
	}

	return &Broker{
		Producer: producer,

		config: cfg,
	}, nil
}

func (b *Broker) Shutdown() {
	b.Producer.Close()
}

type BookNotifyEvent struct {
	Ticket string
	Data   *bookrpc.BookRequest
}

func (b *Broker) BookNotify(event *BookNotifyEvent) error {
	const op = "BookNotify()"

	partition, offset, err := b.Producer.SendMessage(&sarama.ProducerMessage{
		Topic:     b.config.KafkaConfig.Topic,
		Value:     sarama.ByteEncoder(utils.MarshalJSON(event)),
		Offset:    b.config.KafkaConfig.Offset,
		Partition: b.config.KafkaConfig.Partition,
	})
	if err != nil {
		b.log.Logs.BrokerLog.Error(
			"can't produce a message",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)

		return err
	}
	b.log.Logs.BookLog.Debug(
		"message produced",
		slog.String("op", op),
		slog.Any("partition", partition),
		slog.Any("offset", offset),
		slog.Any("event", event),
	)

	return nil
}

type UnimplementedBroker struct{}

func (u *UnimplementedBroker) BookNotify(event *BookNotifyEvent) error { return nil }
