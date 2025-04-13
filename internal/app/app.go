package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xoticdsign/bookamovie/internal/app/bookamovieapp"
	broker "github.com/xoticdsign/bookamovie/internal/broker/kafka"
	storage "github.com/xoticdsign/bookamovie/internal/storage/sqlite"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type App struct {
	BookaMovie *bookamovieapp.App
	Storage    *storage.Storage
	Broker     *broker.Broker
}

func New() (*App, error) {
	cfg := utils.LoadConfig()

	s, err := storage.New(cfg)
	if err != nil {
		return &App{}, err
	}

	b, err := broker.New(cfg)
	if err != nil {
		return &App{}, err
	}

	bookamovie := bookamovieapp.New(cfg, s, b)

	return &App{
		BookaMovie: bookamovie,
		Storage:    s,
		Broker:     b,
	}, nil
}

func (a *App) Run() {
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
		// LOG SIGNAL

	case err := <-errChan:
		// LOG ERROR
	}

	shutdown(a.BookaMovie, a.Storage, a.Broker)
}

func shutdown(bookamovie *bookamovieapp.App, storage *storage.Storage, broker *broker.Broker) {
	broker.Shutdown()
	storage.Shutdown()
	bookamovie.Shutdown()
}
