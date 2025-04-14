package book

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xoticdsign/book/internal/lib/logger"
	bookservice "github.com/xoticdsign/book/internal/services/book"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

// App{} represents the gRPC server application for the book service.
//
// It handles configuration, logging, and startup/shutdown lifecycle.
type App struct {
	Server *grpc.Server
	Log    *logger.Logger

	config utils.Config
}

// New() initializes and returns a new instance of the book gRPC App.
//
// It wires together logging, configuration, storage, and message broker.
func New(log *logger.Logger, cfg utils.Config, storage bookservice.Querier, broker bookservice.Brokerer) *App {
	server := grpc.NewServer()

	bookrpc.RegisterBookServer(server, &Api{Service: bookservice.New(cfg, log, storage, broker)})

	return &App{
		Server: server,
		Log:    log,

		config: cfg,
	}
}

// Run() starts the gRPC server using the configured network and address.
//
// It blocks and returns any critical error if the server fails to start.
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

// Shutdown() gracefully stops the gRPC server.
func (a *App) Shutdown() {
	a.Server.GracefulStop()
}

// Servicer() defines the interface for the booking service logic.
//
// It is implemented by the internal book service layer.
type Servicer interface {
	Book(ctx context.Context, data *bookrpc.BookRequest) (*bookrpc.BookResponse, error)
}

// api{} is the gRPC handler for the Book service.
//
// It adapts incoming gRPC calls to the internal Servicer logic.
type Api struct {
	bookrpc.UnimplementedBookServer

	Service Servicer
}

// Book() handles incoming gRPC requests to book a movie ticket.
//
// It validates input and delegates to the business logic service layer. Returns appropriate gRPC errors for invalid or duplicate requests.
func (a *Api) Book(ctx context.Context, req *bookrpc.BookRequest) (*bookrpc.BookResponse, error) {
	ok := utils.ValidateBookRequest(req)
	if !ok {
		return &bookrpc.BookResponse{}, status.Error(codes.InvalidArgument, "required request arguments must be specified")
	}

	resp, err := a.Service.Book(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, bookservice.ErrDuplicate):
			return &bookrpc.BookResponse{}, status.Error(codes.AlreadyExists, bookservice.ErrDuplicate.Error())

		default:
			return &bookrpc.BookResponse{}, status.Error(codes.Internal, "internal error")
		}
	}

	return resp, nil
}
