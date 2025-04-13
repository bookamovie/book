package bookamovieservice

import (
	"context"

	bookamovierpc "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v2"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Book(ctx context.Context, req *bookamovierpc.BookRequest) (*bookamovierpc.BookResponse, error) {
	return &bookamovierpc.BookResponse{}, nil
}
