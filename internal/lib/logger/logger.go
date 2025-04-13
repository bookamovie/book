package logger

import (
	"log/slog"
	"os"

	"github.com/xoticdsign/book/internal/utils"
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

func New(logMode string) (*Logger, error) {
	var appLog *slog.Logger
	var bookLog *slog.Logger
	var storageLog *slog.Logger
	var brokerLog *slog.Logger

	var b *os.File
	var s *os.File
	var br *os.File
	var err error

	switch logMode {
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
		b, err = utils.OpenLog("book_log.json")
		if err != nil {
			return &Logger{}, err
		}
		s, err = utils.OpenLog("storage_log.json")
		if err != nil {
			return &Logger{}, err
		}
		br, err = utils.OpenLog("broker_log.json")
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
