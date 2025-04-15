package suite

import (
	"testing"

	"github.com/xoticdsign/book/internal/app"
	bookapp "github.com/xoticdsign/book/internal/app/book"
	"github.com/xoticdsign/book/internal/lib/logger"
	bookservice "github.com/xoticdsign/book/internal/services/book"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T

	App    *app.App
	Client bookrpc.BookClient
}

func New(t *testing.T, cfg utils.Config, log *logger.Logger, storage bookservice.Querier, broker bookservice.Brokerer) *Suite {
	t.Helper()
	t.Parallel()

	app := &app.App{
		Book:    bookapp.New(log, cfg, storage, broker),
		Storage: storage,
		Broker:  broker,
		Log:     log,
		Config:  cfg,
	}

	conn, err := grpc.NewClient(cfg.BookConfig.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := bookrpc.NewBookClient(conn)

	return &Suite{
		T: t,

		App:    app,
		Client: client,
	}
}
