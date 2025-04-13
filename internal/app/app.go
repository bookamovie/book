package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xoticdsign/bookamovie/internal/app/bookamovieapp"
	storage "github.com/xoticdsign/bookamovie/internal/storage/sqlite"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type App struct {
	BookaMovie *bookamovieapp.App
}

func New() (*App, error) {
	cfg := utils.LoadConfig()

	storage, err := storage.New(cfg)
	if err != nil {
		return &App{}, err
	}

	bookamovie := bookamovieapp.New(cfg, storage)

	return &App{
		BookaMovie: bookamovie,
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

	shutdown(a.BookaMovie)
}

func shutdown(bookamovie *bookamovieapp.App) {
	bookamovie.Shutdown()
}
