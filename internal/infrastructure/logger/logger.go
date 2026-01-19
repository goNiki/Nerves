package logger

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/goNiki/Nerves/internal/config"
)

//пока не используется
// const (
// 	devEnv     string = "dev"
// 	testEnv    string = "test"
// 	stadingEnv string = "stading"
// 	prodEnv    string = "prod"
// )

const (
	jsonFormat string = "json"
	textFormat string = "text"
)

const (
	stdoutOutput string = "stdout"
	stderrOutput string = "stderr"
	fileOutput   string = "file"
)

const (
	debugLevel string = "debug"
	infoLevel  string = "info"
	warnLevel  string = "warn"
	errorLevel string = "error"
)

type Logger struct {
	Log *slog.Logger
}

func InitLogger(cfg config.Log) *Logger {
	output := parseOutput(cfg.Output(), cfg.FilePath())

	level := parseLevel(cfg.Level())

	handler := parseFormat(cfg.Format(), *output, level)

	log := slog.New(handler)

	return &Logger{
		Log: log,
	}
}

func parseLevel(level string) slog.Level {
	switch level {
	case debugLevel:
		return slog.LevelDebug
	case infoLevel:
		return slog.LevelInfo
	case warnLevel:
		return slog.LevelWarn
	case errorLevel:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

func parseFormat(format string, output os.File, level slog.Level) *slog.JSONHandler {
	switch format {
	case jsonFormat:
		return slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: level})
	case textFormat:
		return slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: level})
	default:
		return slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: level})
	}

}

func parseOutput(output string, filePath string) *os.File {
	switch output {
	case stdoutOutput:
		return os.Stdout
	case stderrOutput:
		return os.Stderr
	case fileOutput:
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("faidel to open log file: %v", err)
		}
		return file
	default:
		return os.Stdout
	}
}

// оберка для логирования ошибки с необходимыми парамеирами
func (l *Logger) Error(ctx context.Context, op string, err error) {
	Err := slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}

	l.Log.Error("Error",
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
		Err,
	)
}
