package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	bookapp "github.com/bookamovie/book/internal/app/book"
	broker "github.com/bookamovie/book/internal/broker/kafka"
	"github.com/bookamovie/book/internal/lib/logger"
	bookservice "github.com/bookamovie/book/internal/services/book"
	storage "github.com/bookamovie/book/internal/storage/sqlite"
	"github.com/bookamovie/book/internal/utils"
)

// App{} coordinates the main components of the bookamovie service.
//
// It contains the gRPC application logic, storage backend, broker, and shared logger/config.
type App struct {
	Book    *bookapp.App
	Storage bookservice.Querier
	Broker  bookservice.Brokerer
	Log     *logger.Logger
	Config  utils.Config
}

// New() initializes the App with all necessary components.
//
// It loads config, sets up logging, storage, broker, and gRPC logic. Returns a pointer to App or an error on failure.
func New() (*App, error) {
	cfg, err := utils.LoadConfig()
	if err != nil {
		return &App{}, err
	}

	log, err := logger.New()
	if err != nil {
		return &App{}, err
	}

	s, err := storage.New(cfg, log)
	if err != nil {
		return &App{}, err
	}

	br, err := broker.New(cfg, log)
	if err != nil {
		return &App{}, err
	}

	book := bookapp.New(log, cfg, s, br)

	return &App{
		Book:    book,
		Storage: s,
		Broker:  br,
		Log:     log,
		Config:  cfg,
	}, nil
}

// Run() starts the App, launching the gRPC server and listening for OS shutdown signals.
//
// It blocks until an interrupt or error occurs, then gracefully shuts everything down.
func (a *App) Run() {
	const op = "Run()"

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error, 1)

	a.Log.Logs.AppLog.Info(
		"started an app",
		slog.String("op", op),
		slog.Any("config", a.Config),
	)

	go func() {
		err := a.Book.Run()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-sigChan:
		a.Log.Logs.AppLog.Info(
			"attempting to shut down gracefully",
			slog.String("op", op),
		)

	case err := <-errChan:
		a.Log.Logs.AppLog.Error(
			"error happened, while running. attempting to shut down gracefully",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
	}

	a.Shutdown()

	a.Log.Logs.AppLog.Info(
		"shut down gracefully",
		slog.String("op", op),
	)
}

// shutdown() gracefully shuts down all services in the correct order:
//
// broker → storage → gRPC app → logger.
func (a *App) Shutdown() {
	a.Broker.Shutdown()
	a.Storage.Shutdown()
	a.Book.Shutdown()
	a.Log.Shutdown()
}
