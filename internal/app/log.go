package app

import (
	"log/slog"

	"github.com/WhiCu/stgorders/internal/config"
	"github.com/WhiCu/stgorders/pkg/logger"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getLogger(cfg *config.LoggerConfig) *slog.Logger {
	h := make([]slog.Handler, 0, 2)

	if cfg.Path != "" {
		logFile := &lumberjack.Logger{
			Filename:  cfg.Path,
			MaxSize:   cfg.Size,
			LocalTime: true,
			Compress:  cfg.Compress,
		}

		h = append(h, slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	h = append(h, logger.MustInitLogger(cfg.Level))

	return slog.New(slogmulti.Fanout(h...))

}
