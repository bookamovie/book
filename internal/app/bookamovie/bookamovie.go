package bookamovie

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bookamovierpc "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v3"
	broker "github.com/xoticdsign/bookamovie/internal/broker/kafka"
	"github.com/xoticdsign/bookamovie/internal/lib/logger"
	bookamovieservice "github.com/xoticdsign/bookamovie/internal/services/bookamovie"
	storage "github.com/xoticdsign/bookamovie/internal/storage/sqlite"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type App struct {
	Server *grpc.Server

	log    *logger.Logger
	config *utils.Config
}

func New(log *logger.Logger, cfg *utils.Config, storage *storage.Storage, broker *broker.Broker) *App {
	server := grpc.NewServer()

	bookamovierpc.RegisterBookaMovieServer(server, &api{service: bookamovieservice.New(cfg, log, storage, broker)})

	return &App{
		Server: server,

		log:    log,
		config: cfg,
	}
}

func (a *App) Run() error {
	listener, err := net.Listen(a.config.BookaMovieConfig.Network, a.config.BookaMovieConfig.Address)
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
	Book(ctx context.Context, data *bookamovierpc.BookRequest) (*bookamovierpc.BookResponse, error)
}

type api struct {
	bookamovierpc.UnimplementedBookaMovieServer
	service Servicer
}

func (a *api) Book(ctx context.Context, req *bookamovierpc.BookRequest) (*bookamovierpc.BookResponse, error) {
	ok := utils.ValidateBookRequest(req)
	if !ok {
		return &bookamovierpc.BookResponse{}, status.Error(codes.InvalidArgument, "required request arguments must be specified")
	}

	resp, err := a.service.Book(ctx, req)
	if err != nil {
		// COME UP WITH A BETTER ERROR HANDLING

		return &bookamovierpc.BookResponse{}, status.Error(codes.Internal, "internal error")
	}

	return &bookamovierpc.BookResponse{}, nil
}
