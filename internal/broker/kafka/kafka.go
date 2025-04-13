package kafka

import (
	"github.com/IBM/sarama"

	bookamovierpc "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v3"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type Broker struct {
	Producer sarama.SyncProducer
}

func New(cfg *utils.Config) (*Broker, error) {
	producer, err := sarama.NewSyncProducer(cfg.KafkaConfig.Addresses, sarama.NewConfig())
	if err != nil {
		return &Broker{}, err
	}

	return &Broker{
		Producer: producer,
	}, nil
}

func (b *Broker) Shutdown() {
	b.Producer.Close()
}

type BookNotifyEvent struct {
	Ticket string
	Data   *bookamovierpc.BookRequest
}

func (b *Broker) BookNotify(event *BookNotifyEvent) error {
	return nil
}
