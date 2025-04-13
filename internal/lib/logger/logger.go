package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/xoticdsign/book/internal/utils"
)

const lmEnvName = "LOG_MODE"

var (
	ErrLogModeNotSpecified = fmt.Errorf("%s env variable must be specified", lmEnvName)
	ErrWrongLogger         = fmt.Errorf("specified logger does not exists")
)

type Logger struct {
	Logs     Logs
	LogFiles []*os.File
}

type Logs struct {
	AppLog     *slog.Logger
	BookLog    *slog.Logger
	StorageLog *slog.Logger
	BrokerLog  *slog.Logger
}

type silentHandler struct{}

func (s silentHandler) Enabled(_ context.Context, _ slog.Level) bool  { return false }
func (s silentHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (s silentHandler) WithAttrs(_ []slog.Attr) slog.Handler          { return s }
func (s silentHandler) WithGroup(_ string) slog.Handler               { return s }

func New() (*Logger, error) {
	logMode := os.Getenv(lmEnvName)
	if logMode == "" {
		return &Logger{}, ErrLogModeNotSpecified
	}
	defer os.Unsetenv(lmEnvName)

	var appLog *slog.Logger
	var bookLog *slog.Logger
	var storageLog *slog.Logger
	var brokerLog *slog.Logger

	var b *os.File
	var s *os.File
	var br *os.File
	var err error

	switch logMode {
	case "silent":
		bookLog = slog.New(silentHandler{})
		storageLog = slog.New(silentHandler{})
		brokerLog = slog.New(silentHandler{})

	case "local":
		bookLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		storageLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		brokerLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case "dev":
		b, err = utils.OpenLog("book_dev_log.json")
		if err != nil {
			return &Logger{}, err
		}
		s, err = utils.OpenLog("storage_dev_log.json")
		if err != nil {
			return &Logger{}, err
		}
		br, err = utils.OpenLog("broker_dev_log_dev.json")
		if err != nil {
			return &Logger{}, err
		}

		bookLog = slog.New(slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		storageLog = slog.New(slog.NewJSONHandler(s, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		brokerLog = slog.New(slog.NewJSONHandler(br, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case "prod":
		b, err = utils.OpenLog("book_prod_log.json")
		if err != nil {
			return &Logger{}, err
		}
		s, err = utils.OpenLog("storage_prod_log.json")
		if err != nil {
			return &Logger{}, err
		}
		br, err = utils.OpenLog("broker_prod_log.json")
		if err != nil {
			return &Logger{}, err
		}

		bookLog = slog.New(slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		storageLog = slog.New(slog.NewJSONHandler(s, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		brokerLog = slog.New(slog.NewJSONHandler(br, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	default:
		return &Logger{}, ErrWrongLogger
	}

	appLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	var logFiles []*os.File

	logFiles = append(logFiles, b, s, br)

	return &Logger{
		Logs: Logs{
			AppLog:     appLog,
			BookLog:    bookLog,
			StorageLog: storageLog,
			BrokerLog:  brokerLog,
		},
		LogFiles: logFiles,
	}, nil
}

func (l *Logger) Shutdown() {
	for _, file := range l.LogFiles {
		file.Close()
	}
}
