package book

import (
	"context"

	broker "github.com/xoticdsign/book/internal/broker/kafka"
	"github.com/xoticdsign/book/internal/lib/logger"
	storage "github.com/xoticdsign/book/internal/storage/sqlite"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
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

func (s *Service) Book(ctx context.Context, data *bookrpc.BookRequest) (*bookrpc.BookResponse, error) {
	err := s.Storage.Book(&storage.BookQuery{
		Ticket:/* TICKET GEN */ "",
		Data: data,
	})
	if err != nil {
		// ERROR HANDLING
	}

	err = s.Broker.BookNotify(&broker.BookNotifyEvent{
		Ticket:/* TICKET GEN */ "",
		Data: data,
	})
	if err != nil {
		// ERROR HANDLING
	}

	return &bookrpc.BookResponse{}, nil
}
