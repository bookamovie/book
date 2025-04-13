package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	bookamovieapp "github.com/xoticdsign/bookamovie/internal/app/bookamovie"
	broker "github.com/xoticdsign/bookamovie/internal/broker/kafka"
	"github.com/xoticdsign/bookamovie/internal/lib/logger"
	storage "github.com/xoticdsign/bookamovie/internal/storage/sqlite"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type App struct {
	BookaMovie *bookamovieapp.App
	Storage    *storage.Storage
	Broker     *broker.Broker

	log    *logger.Logger
	config *utils.Config
}

func New() (*App, error) {
	cfg := utils.LoadConfig()

	log, err := logger.New(cfg.LogMode)
	if err != nil {
		return &App{}, err
	}

	s, err := storage.New(cfg)
	if err != nil {
		return &App{}, err
	}

	br, err := broker.New(cfg)
	if err != nil {
		return &App{}, err
	}

	bookamovie := bookamovieapp.New(log, cfg, s, br)

	return &App{
		BookaMovie: bookamovie,
		Storage:    s,
		Broker:     br,

		log:    log,
		config: cfg,
	}, nil
}

func (a *App) Run() {
	const op = "Run()"

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error, 1)

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

func shutdown(log *logger.Logger, bookamovie *bookamovieapp.App, storage *storage.Storage, broker *broker.Broker) {
	broker.Shutdown()
	storage.Shutdown()
	bookamovie.Shutdown()
	log.Shutdown()
}
