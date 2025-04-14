package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrations = "migrations/sqlite"
	storage    = "storage/db.sqlite"
)

func main() {
	m, err := migrate.New("file://"+migrations, "sqlite3://"+storage)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		panic(err)
	}
}
