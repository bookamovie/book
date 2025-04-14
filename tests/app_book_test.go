package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	broker "github.com/xoticdsign/book/internal/broker/kafka"
	"github.com/xoticdsign/book/internal/lib/logger"
	storage "github.com/xoticdsign/book/internal/storage/sqlite"
	"github.com/xoticdsign/book/internal/utils"
	"github.com/xoticdsign/book/tests/suite"
	bookrcp "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestBook_Functional(t *testing.T) {
	cases := []struct {
		name        string
		in          *bookrcp.BookRequest
		expectedErr error
	}{
		{
			name: "happy test",
			in: &bookrcp.BookRequest{
				Cinema: &bookrcp.Cinema{
					Name:     "cinema",
					Location: "location",
				},
				Movie: &bookrcp.Movie{
					Title: "title",
				},
				Session: &bookrcp.Session{
					Screen: 1,
					Seat:   1,
					Date:   timestamppb.New(time.Now()),
				},
			},
			expectedErr: nil,
		},
	}

	os.Setenv(suite.CpEnvName, "config/test.yaml")
	os.Setenv(suite.LmEnvName, "silent")

	cfg, err := utils.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	storage, err := storage.New(cfg, log)
	if err != nil {
		panic(err)
	}

	suite := suite.New(t, cfg, log, storage, &broker.UnimplementedBroker{})

	go suite.App.Run()
	defer suite.App.Shutdown()

	time.Sleep(time.Second)

	for _, cs := range cases {
		suite.T.Run(cs.name, func(t *testing.T) {
			resp, err := suite.Client.Book(context.Background(), cs.in)
			assert.Equal(t, cs.expectedErr, err)
			assert.NotEmpty(t, resp.GetOrder().GetTicket())
		})
	}
}
