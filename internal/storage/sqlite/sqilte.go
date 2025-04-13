package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

type Storage struct {
	DB *sql.DB

	config utils.Config
}

func New(cfg utils.Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", cfg.SQLiteConfig.Address)
	if err != nil {
		return &Storage{}, err
	}

	return &Storage{
		DB: db,

		config: cfg,
	}, nil
}

func (s *Storage) Shutdown() {
	s.DB.Close()
}

type BookQuery struct {
	Ticket string
	Data   *bookrpc.BookRequest
}

func (s *Storage) Book(query *BookQuery) error {
	return nil
}

type UnimplementedStorage struct{}

func (u *UnimplementedStorage) Book(query *BookQuery) error { return nil }
