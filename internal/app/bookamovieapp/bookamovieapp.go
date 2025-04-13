package bookamovieapp

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bookamovierpc "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v2"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type App struct {
	Server *grpc.Server

	config *utils.Config
}

func New(cfg *utils.Config) *App {
	server := grpc.NewServer()

	bookamovierpc.RegisterBookaMovieServer(server, &api{})

	return &App{
		Server: server,

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
	Book(ctx context.Context, req *bookamovierpc.BookRequest) (*bookamovierpc.BookResponse, error)
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
