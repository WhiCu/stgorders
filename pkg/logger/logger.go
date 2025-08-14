package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func MustInitLogger(level string) slog.Handler {
	const op = "logger.MustInitLogger"

	var log slog.Handler

	out := os.Stdout

	switch level {
	case "debug":
		log = slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case "info":
		log = slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case "warn":
		log = slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})
	case "error":
		log = slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: slog.LevelError,
		})
	default:
		panic(fmt.Sprintf("%s: level - {%s} not in {debug, info, warn, error}", op, level))
	}
	return log
}
