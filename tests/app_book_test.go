package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	broker "github.com/xoticdsign/book/internal/broker/kafka"
	"github.com/xoticdsign/book/internal/lib/logger"
	storage "github.com/xoticdsign/book/internal/storage/sqlite"
	"github.com/xoticdsign/book/internal/utils"
	"github.com/xoticdsign/book/tests/suite"
	bookrcp "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestBook_Functional(t *testing.T) {
	cases := []struct {
		name             string
		in               *bookrcp.BookRequest
		expectedCode     codes.Code
		expectedResponse bool
	}{
		{
			name: "happy case",
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
			expectedCode:     codes.OK,
			expectedResponse: true,
		},
		{
			name: "bad request",
			in: &bookrcp.BookRequest{
				Cinema: &bookrcp.Cinema{
					Name:     "cinema",
					Location: "location",
				},
			},
			expectedCode:     codes.InvalidArgument,
			expectedResponse: false,
		},
	}

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
			st, ok := status.FromError(err)
			if ok {
				assert.Equal(t, cs.expectedCode, st.Code())
			}
			if cs.expectedResponse {
				assert.NotEmpty(t, resp.GetOrder().GetTicket())
			} else {
				assert.Empty(t, resp.GetOrder().GetTicket())
			}
		})
	}
}
