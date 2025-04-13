package bookamovieservice

import (
	"context"

	bookamovierpc "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v2"
)

type Querier interface {
	Book(data *bookamovierpc.BookRequest) error
}

type Service struct {
	Storage Querier
}

func New(storage Querier) *Service {
	return &Service{
		Storage: storage,
	}
}

func (s *Service) Book(ctx context.Context, data *bookamovierpc.BookRequest) (*bookamovierpc.BookResponse, error) {
	err := s.Storage.Book(data)
	if err != nil {
		// ERROR HANDLING
	}
	return &bookamovierpc.BookResponse{}, nil
}
