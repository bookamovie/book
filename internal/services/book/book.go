package book

import (
	"context"
	"errors"
	"fmt"

	"github.com/thanhpk/randstr"
	broker "github.com/xoticdsign/book/internal/broker/kafka"
	"github.com/xoticdsign/book/internal/lib/logger"
	storage "github.com/xoticdsign/book/internal/storage/sqlite"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

var (
	ErrDuplicateOrder = fmt.Errorf("order duplicated")
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
	config utils.Config
}

func New(cfg utils.Config, log *logger.Logger, storage *storage.Storage, broker *broker.Broker) *Service {
	return &Service{
		Storage: storage,

		log:    log,
		config: cfg,
	}
}

func (s *Service) Book(ctx context.Context, data *bookrpc.BookRequest) (*bookrpc.BookResponse, error) {
	ticket := randstr.Dec(12)

	err := s.Storage.Book(&storage.BookQuery{
		Ticket: ticket,
		Data:   data,
	})
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return &bookrpc.BookResponse{}, ErrDuplicateOrder
		}
		return &bookrpc.BookResponse{}, err
	}

	err = s.Broker.BookNotify(&broker.BookNotifyEvent{
		Ticket: ticket,
		Data:   data,
	})
	if err != nil {
		return &bookrpc.BookResponse{}, err
	}

	return &bookrpc.BookResponse{}, nil
}

type UnimplementedService struct{}

func (u *UnimplementedService) Book(ctx context.Context, data *bookrpc.BookRequest) (*bookrpc.BookResponse, error) {
	return &bookrpc.BookResponse{}, nil
}
