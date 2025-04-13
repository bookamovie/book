package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v2"
	"github.com/xoticdsign/bookamovie/internal/utils"
)

type Storage struct {
	DB *sql.DB
}

func New(cfg *utils.Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", cfg.SQLiteConfig.Address)
	if err != nil {
		return &Storage{}, err
	}

	return &Storage{
		DB: db,
	}, nil
}

func (s *Storage) Book(data *bookamovie.BookRequest) error {
	return nil
}
