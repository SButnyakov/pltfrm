package app

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/handlers"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
	"url-shortener/pkg/config"
	"url-shortener/pkg/database/postgres"
	"url-shortener/pkg/logger"
)

func Run(env string) {
	logger.Init(env)
	slog.Info("Logger initialized")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config:", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Config loaded")

	pg, err := postgres.Connect(cfg)
	if err != nil {
		slog.Error("failed to connect to postgres:", slog.Any("error", err))
		os.Exit(1)
	}
	defer pg.Close()
	slog.Info("Connected to postgres")

	urlHandler := handlers.NewUrlHandler(service.NewUrlService(repository.NewUrlRepository(pg)))

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/urls", func(r chi.Router) {
		r.Post("/", urlHandler.CreateURL)
	})
	r.Get("/{url}", urlHandler.Redirect)
	r.Delete("/{url}", urlHandler.DeleteByURL)

	server := http.Server{
		Addr:         cfg.HTTP.Port,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		Handler:      r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("Starting server", slog.String("port", cfg.HTTP.Port))
	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error:", slog.Any("msg", err))
			os.Exit(1)
		}
	}()

	<-stop
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server:", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Server gracefully stopped")
	os.Exit(0)
}
