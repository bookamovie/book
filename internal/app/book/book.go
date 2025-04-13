package book

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	broker "github.com/xoticdsign/book/internal/broker/kafka"
	"github.com/xoticdsign/book/internal/lib/logger"
	bookservice "github.com/xoticdsign/book/internal/services/book"
	storage "github.com/xoticdsign/book/internal/storage/sqlite"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

type App struct {
	Server *grpc.Server

	log    *logger.Logger
	config utils.Config
}

func New(log *logger.Logger, cfg utils.Config, storage *storage.Storage, broker *broker.Broker) *App {
	server := grpc.NewServer()

	bookrpc.RegisterBookServer(server, &api{service: bookservice.New(cfg, log, storage, broker)})

	return &App{
		Server: server,

		log:    log,
		config: cfg,
	}
}

func (a *App) Run() error {
	listener, err := net.Listen(a.config.BookConfig.Network, a.config.BookConfig.Address)
	if err != nil {
		return err
	}

	err = a.Server.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown() {
	a.Server.GracefulStop()
}

type Servicer interface {
	Book(ctx context.Context, data *bookrpc.BookRequest) (*bookrpc.BookResponse, error)
}

type api struct {
	bookrpc.UnimplementedBookServer

	service Servicer
}

func (a *api) Book(ctx context.Context, req *bookrpc.BookRequest) (*bookrpc.BookResponse, error) {
	ok := utils.ValidateBookRequest(req)
	if !ok {
		return &bookrpc.BookResponse{}, status.Error(codes.InvalidArgument, "required request arguments must be specified")
	}

	_, err := a.service.Book(ctx, req)
	if err != nil {
		// COME UP WITH A BETTER ERROR HANDLING

		return &bookrpc.BookResponse{}, status.Error(codes.Internal, "internal error")
	}

	return &bookrpc.BookResponse{}, nil
}
