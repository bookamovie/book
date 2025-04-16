package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/bookamovie/book/internal/utils"
)

const lmEnvName = "LOG_MODE"

var (
	ErrLogModeNotSpecified = fmt.Errorf("%s env variable must be specified", lmEnvName)
	ErrWrongLogger         = fmt.Errorf("specified log mode does not exists")
)

// Logger{} is the central logging component for the app.
//
// It holds references to individual subsystem loggers and file handles.
type Logger struct {
	Logs     Logs
	LogFiles []*os.File
}

// Logs{} contains loggers categorized by system component.
type Logs struct {
	AppLog     *slog.Logger
	BookLog    *slog.Logger
	StorageLog *slog.Logger
	BrokerLog  *slog.Logger
}

// silentHandler{} is a no-op slog handler that discards all logs.
type silentHandler struct{}

func (s silentHandler) Enabled(_ context.Context, _ slog.Level) bool  { return false }
func (s silentHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (s silentHandler) WithAttrs(_ []slog.Attr) slog.Handler          { return s }
func (s silentHandler) WithGroup(_ string) slog.Handler               { return s }

// New() initializes and returns a configured Logger based on the LOG_MODE environment variable.
//
// It supports multiple modes: "silent", "local", "dev", and "prod".
//   - silent: disables all logs.
//   - local: prints all logs to stdout.
//   - dev: writes logs to JSON files in log/dev/.
//   - prod: writes logs to JSON files in log/.
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
		appLog = slog.New(silentHandler{})
		bookLog = slog.New(silentHandler{})
		storageLog = slog.New(silentHandler{})
		brokerLog = slog.New(silentHandler{})

	case "local":
		appLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

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
		err = os.MkdirAll("log/dev", 0777)
		if err != nil {
			return &Logger{}, err
		}

		b, err = utils.OpenFile("log/dev/book.log")
		if err != nil {
			return &Logger{}, err
		}
		s, err = utils.OpenFile("log/dev/storage.log")
		if err != nil {
			return &Logger{}, err
		}
		br, err = utils.OpenFile("log/dev/broker.log")
		if err != nil {
			return &Logger{}, err
		}

		appLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

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
		err = os.MkdirAll("log", 0777)
		if err != nil {
			return &Logger{}, err
		}

		b, err = utils.OpenFile("log/book.log")
		if err != nil {
			return &Logger{}, err
		}
		s, err = utils.OpenFile("log/storage.log")
		if err != nil {
			return &Logger{}, err
		}
		br, err = utils.OpenFile("log/broker.log")
		if err != nil {
			return &Logger{}, err
		}

		appLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

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

// Shutdown() closes any log files opened by the logger.
func (l *Logger) Shutdown() {
	for _, file := range l.LogFiles {
		file.Close()
	}
}
