package bookamovie

import (
	"context"

	bookamovierpc "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v3"
	broker "github.com/xoticdsign/bookamovie/internal/broker/kafka"
	"github.com/xoticdsign/bookamovie/internal/lib/logger"
	storage "github.com/xoticdsign/bookamovie/internal/storage/sqlite"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type Querier interface {
	Book(query *storage.BookQuery) error
}

type Brokerer interface {
	BookNotify(event *broker.BookNotifyEvent) error
}

type Service struct {
	Storage Querier
	Broker  Brokerer

	log    *logger.Logger
	config *utils.Config
}

func New(cfg *utils.Config, log *logger.Logger, storage *storage.Storage, broker *broker.Broker) *Service {
	return &Service{
		Storage: storage,

		log:    log,
		config: cfg,
	}
}

func (s *Service) Book(ctx context.Context, data *bookamovierpc.BookRequest) (*bookamovierpc.BookResponse, error) {
	err := s.Storage.Book(&storage.BookQuery{})
	if err != nil {
		// ERROR HANDLING
	}

	err = s.Broker.BookNotify(&broker.BookNotifyEvent{})
	if err != nil {
		// ERROR HANDLING
	}

	return &bookamovierpc.BookResponse{}, nil
}
