package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/logger"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	logger.Init(cfg.Env)

	slog.Info("starting application",
		slog.String("env", cfg.Env),
		slog.Any("cfg", cfg),
		slog.Int("port", cfg.GRPC.Port),
	)

	slog.Debug("debug message")

	slog.Error("error message")

	slog.Warn("warn message")

	application := app.New(cfg)

	go application.GRPCSrv.MustRun()

	// TODO: инициализировать приложение (app)

	// TODO: запустить gRPC-сервер приложения

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCSrv.Stop()

	slog.Info("application stopped")
}
