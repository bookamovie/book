package kafka

import (
	"github.com/IBM/sarama"

	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

type Broker struct {
	Producer sarama.SyncProducer

	config utils.Config
}

func New(cfg utils.Config) (*Broker, error) {
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
	return nil
}

type UnimplementedBroker struct{}

func (u *UnimplementedBroker) BookNotify(event *BookNotifyEvent) error { return nil }
