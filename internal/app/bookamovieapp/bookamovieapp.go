package bookamovieapp

import (
	"context"
	"net"

	"google.golang.org/grpc"

	bookamoviev1 "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v1"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type App struct {
	Server *grpc.Server

	config *utils.Config
}

func New(cfg *utils.Config) *App {
	server := grpc.NewServer()

	bookamoviev1.RegisterBookaMovieServer(server, &api{})

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

type api struct {
	bookamoviev1.UnimplementedBookaMovieServer
}

func (a *api) Movies(ctx context.Context, req *bookamoviev1.MoviesRequest) (*bookamoviev1.MoviesResponse, error) {
	return &bookamoviev1.MoviesResponse{}, nil
}

func (a *api) Movie(ctx context.Context, req *bookamoviev1.MovieRequest) (*bookamoviev1.MovieResponse, error) {
	return &bookamoviev1.MovieResponse{}, nil
}

func (a *api) Cinema(ctx context.Context, req *bookamoviev1.CinemaRequest) (*bookamoviev1.CinemaResponse, error) {
	return &bookamoviev1.CinemaResponse{}, nil
}

func (a *api) Book(ctx context.Context, req *bookamoviev1.BookRequest) (*bookamoviev1.BookResponse, error) {
	return &bookamoviev1.BookResponse{}, nil
}
