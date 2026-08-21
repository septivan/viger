package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/septivan/viger/backend/internal/adapters/outbound/memory"
	gameservices "github.com/septivan/viger/backend/internal/core/game/services"
	reviewservices "github.com/septivan/viger/backend/internal/core/review/services"
	"github.com/septivan/viger/backend/internal/platform/config"
	"github.com/septivan/viger/backend/internal/platform/httpapi"
	"github.com/septivan/viger/backend/internal/platform/observability"
	"github.com/septivan/viger/backend/internal/platform/realtime"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	configuration, err := config.Load()
	if err != nil {
		logger.Error("configuration is invalid", "error", err)
		os.Exit(1)
	}
	games, reviews, err := memory.Seed()
	if err != nil {
		logger.Error("seed data is invalid", "error", err)
		os.Exit(1)
	}
	store, err := memory.New(games, reviews)
	if err != nil {
		logger.Error("memory store could not be created", "error", err)
		os.Exit(1)
	}
	metrics := &observability.Metrics{}
	hub := realtime.NewHub(configuration.AllowedOrigins, configuration.MaximumWSConnections, metrics)
	gameService := gameservices.New(store, store)
	reviewService := reviewservices.New(store, store, hub, reviewservices.SystemClock{}, reviewservices.RandomIDGenerator{})
	router := httpapi.NewRouter(httpapi.Settings{
		Games: gameService, Reviews: reviewService, WebSocket: hub, Metrics: metrics, Logger: logger,
		AllowedOrigins: configuration.AllowedOrigins, ReviewRateLimit: configuration.ReviewRateLimit, ReviewRateWindow: configuration.ReviewRateWindow,
	})

	server := &http.Server{
		Addr: configuration.HTTPAddress, Handler: router,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20,
	}
	shutdownContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		logger.Info("Viger API started", "address", configuration.HTTPAddress, "games", len(games), "reviews", len(reviews))
		if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			logger.Error("HTTP server stopped unexpectedly", "error", listenErr)
			os.Exit(1)
		}
	}()
	<-shutdownContext.Done()
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(context); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Viger API stopped")
}
