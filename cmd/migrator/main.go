package main

import (
	"fmt"
	"os"

	"github.com/bookamovie/book/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	mEnvName = "MIGRATIONS"
	sEnvName = "STORAGE"
)

var (
	ErrMigrationsNotSpecified = fmt.Errorf("%s env variable must be specified", mEnvName)
	ErrStorageNotSpecified    = fmt.Errorf("%s env variable must be specified", sEnvName)
)

func main() {
	migrations := os.Getenv(mEnvName)
	if migrations == "" {
		panic(ErrMigrationsNotSpecified)
	}
	defer os.Unsetenv(mEnvName)

	storage := os.Getenv(sEnvName)
	if storage == "" {
		panic(ErrStorageNotSpecified)
	}
	defer os.Unsetenv(sEnvName)

	file, err := utils.OpenFile(storage)
	if err != nil {
		panic(err)
	}
	file.Close()

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
