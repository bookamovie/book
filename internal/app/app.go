package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	bookapp "github.com/xoticdsign/book/internal/app/book"
	broker "github.com/xoticdsign/book/internal/broker/kafka"
	"github.com/xoticdsign/book/internal/lib/logger"
	storage "github.com/xoticdsign/book/internal/storage/sqlite"
	"github.com/xoticdsign/book/internal/utils"
)

// App{} coordinates the main components of the bookamovie service.
//
// It contains the gRPC application logic, storage backend, broker, and shared logger/config.
type App struct {
	BookaMovie *bookapp.App
	Storage    *storage.Storage
	Broker     *broker.Broker

	log    *logger.Logger
	config utils.Config
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

	bookamovie := bookapp.New(log, cfg, s, br)

	return &App{
		BookaMovie: bookamovie,
		Storage:    s,
		Broker:     br,

		log:    log,
		config: cfg,
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

	a.log.Logs.AppLog.Info(
		"starting an app",
		slog.String("op", op),
		slog.Any("config", a.config),
	)

	go func() {
		err := a.BookaMovie.Run()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-sigChan:
		a.log.Logs.AppLog.Info(
			"attempting to shut down gracefully",
			slog.String("op", op),
		)

	case err := <-errChan:
		a.log.Logs.AppLog.Error(
			"error happened, while running. attempting to shut down gracefully",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
	}

	shutdown(a.log, a.BookaMovie, a.Storage, a.Broker)

	a.log.Logs.AppLog.Info(
		"shut down gracefully",
		slog.String("op", op),
	)
}

// shutdown() gracefully shuts down all services in the correct order:
//
// broker → storage → gRPC app → logger.
func shutdown(log *logger.Logger, bookamovie *bookapp.App, storage *storage.Storage, broker *broker.Broker) {
	broker.Shutdown()
	storage.Shutdown()
	bookamovie.Shutdown()
	log.Shutdown()
}
