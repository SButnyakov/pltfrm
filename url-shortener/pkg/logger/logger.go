package logger

import (
	"log"
	"log/slog"
	"os"
)

const (
	LOCAL = "local"
	DEV   = "dev"
)

func Init(env string) {
	switch env {
	case DEV:
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
		log.Printf("slog initialized: env=dev\n")
	default:
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
		log.Printf("slog initialized: env=local\n")
	}
}
