package tests

import (
	"context"
	"testing"
	"time"

	broker "github.com/bookamovie/book/internal/broker/kafka"
	"github.com/bookamovie/book/internal/lib/logger"
	storage "github.com/bookamovie/book/internal/storage/sqlite"
	"github.com/bookamovie/book/internal/utils"
	"github.com/bookamovie/book/tests/suite"
	bookrcp "github.com/bookamovie/proto/gen/go/book/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestBook_Functional() tests the functional behavior of the Book service by simulating different booking scenarios and checking the expected outcomes.
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
