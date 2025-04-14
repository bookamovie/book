package book

import (
	"context"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/thanhpk/randstr"
	broker "github.com/xoticdsign/book/internal/broker/kafka"
	"github.com/xoticdsign/book/internal/lib/logger"
	storage "github.com/xoticdsign/book/internal/storage/sqlite"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

var (
	ErrDuplicate = fmt.Errorf("this order already exists")
)

// Querier{} abstracts the interface for the storage layer's booking method.
type Querier interface {
	Book(query *storage.BookQuery) error
	Shutdown()
}

// Brokerer{} abstracts the broker (e.g., Kafka) interface for sending booking events.
type Brokerer interface {
	BookNotify(event *broker.BookNotifyEvent) error
	Shutdown()
}

// Service{} handles business logic for booking operations.
type Service struct {
	Storage Querier
	Broker  Brokerer
	Log     *logger.Logger

	config utils.Config
}

// New() creates and returns a new Service instance with dependencies injected.
func New(cfg utils.Config, log *logger.Logger, s Querier, br Brokerer) *Service {
	return &Service{
		Storage: s,
		Broker:  br,
		Log:     log,

		config: cfg,
	}
}

// Book() processes a booking request: generates a ticket, stores the data, and notifies the broker.
//
// Returns a BookResponse with the generated ticket or an error if the operation fails.
func (s *Service) Book(ctx context.Context, data *bookrpc.BookRequest) (*bookrpc.BookResponse, error) {
	ticket := randstr.Dec(12)

	err := s.Storage.Book(&storage.BookQuery{
		Ticket: ticket,
		Data:   data,
	})
	if err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return &bookrpc.BookResponse{}, ErrDuplicate
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

	return &bookrpc.BookResponse{
		Order: &bookrpc.Order{
			Ticket: ticket,
		},
	}, nil
}

// UnimplementedService{} is a placeholder implementation of the service.
//
// Useful for testing or when mocking is required.
type UnimplementedService struct{}

// Book() returns an empty BookResponse and no error.
//
// This satisfies the Servicer interface without performing any logic.
func (u *UnimplementedService) Book(ctx context.Context, data *bookrpc.BookRequest) (*bookrpc.BookResponse, error) {
	return &bookrpc.BookResponse{}, nil
}
