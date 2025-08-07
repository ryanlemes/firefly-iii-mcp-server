package logger

import (
	"log/slog"
	"os"
)

func New(serviceName string) *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log = log.With("service", serviceName)
	return log
}
