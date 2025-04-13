package bookamovieservice

import (
	"context"

	bookamovierpc "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v3"
	broker "github.com/xoticdsign/bookamovie/internal/broker/kafka"
	storage "github.com/xoticdsign/bookamovie/internal/storage/sqlite"
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
}

func New(storage *storage.Storage, broker *broker.Broker) *Service {
	return &Service{
		Storage: storage,
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
