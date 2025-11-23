package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ReviewerAssignmentService/internal/config"
	"ReviewerAssignmentService/internal/database"
	"ReviewerAssignmentService/internal/handler"
	"ReviewerAssignmentService/internal/repository/postgres"
	"ReviewerAssignmentService/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.New()
	slog.Info("Starting service", "port", cfg.ServerPort, "env", "dev")

	dbPool, err := database.NewDBPool(cfg)
	if err != nil {
		slog.Error("Failed to init DB", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	teamRepo := postgres.NewTeamRepository(dbPool)
	userRepo := postgres.NewUserRepository(dbPool)
	prRepo := postgres.NewPrRepository(dbPool)

	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo, prRepo)
	prService := service.NewPRService(prRepo, userRepo)

	httpHandler := handler.New(teamService, userService, prService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: httpHandler.InitRoutes(),
	}

	go func() {
		slog.Info("Server started", "address", fmt.Sprintf(":%d", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited properly")
}
