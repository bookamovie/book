package sqlite

import (
	"database/sql"
	"log/slog"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xoticdsign/book/internal/lib/logger"
	"github.com/xoticdsign/book/internal/utils"
	bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"
)

type Storage struct {
	DB *sql.DB

	log    *logger.Logger
	config utils.Config
}

func New(cfg utils.Config, log *logger.Logger) (*Storage, error) {
	db, err := sql.Open("sqlite3", cfg.SQLiteConfig.Address)
	if err != nil {
		return &Storage{}, err
	}

	return &Storage{
		DB: db,

		log:    log,
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
	const op = "Book()"

	tx, err := s.DB.Begin()
	if err != nil {
		s.log.Logs.StorageLog.Error(
			"can't start a transaction",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)

		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO bookings (id, movie, screen, seat, date, cinema, location) 
	SELECT ?, ?, ?, ?, ?, ?, ? 
	WHERE (
	(SELECT COUNT(*) 
	FROM bookings 
	WHERE screen = ? AND date = ? AND cinema = ? AND location = ?) 
	< 
    (SELECT seats 
	FROM screens 
	WHERE screen = ? AND cinema = ? AND location = ?)
	);`)

	if err != nil {
		s.log.Logs.StorageLog.Error(
			"can't prepare a statement",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)

		return err
	}
	defer tx.Stmt(stmt).Close()

	res, err := tx.Stmt(stmt).Exec(
		query.Ticket,
		query.Data.Movie.Title,
		query.Data.Session.Screen,
		query.Data.Session.Seat,
		query.Data.Session.Date.AsTime(),
		query.Data.Cinema.Name,
		query.Data.Cinema.Location,
		query.Data.Session.Screen,
		query.Data.Session.Date.AsTime(),
		query.Data.Cinema.Name,
		query.Data.Cinema.Location,
		query.Data.Session.Screen,
		query.Data.Cinema.Name,
		query.Data.Cinema.Location,
	)
	if err != nil {
		s.log.Logs.StorageLog.Error(
			"can't execute a statement",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)

		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		s.log.Logs.StorageLog.Warn(
			sqlite3.ErrConstraintUnique.Error(),
			slog.String("op", op),
		)

		return sqlite3.ErrConstraintUnique
	}

	return tx.Commit()
}

type UnimplementedStorage struct{}

func (u *UnimplementedStorage) Book(query *BookQuery) error { return nil }
