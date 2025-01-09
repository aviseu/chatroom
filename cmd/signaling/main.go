package main

import (
	"context"
	"fmt"
	"github.com/aviseu/chatroom/internal/app/signaling"
	"github.com/kelseyhightower/envconfig"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type config struct {
	Log struct {
		Level slog.Level `default:"info"`
	}
	Signaling signaling.Config
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))

	if err := run(context.Background()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// ****************************************************
	// * Load environment variables
	// ****************************************************
	slog.Info("loading environment variables...")
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to process environment variables: %w", err)
	}

	// ****************************************************
	// * Setup logger
	// ****************************************************
	slog.Info("configuring logging...")
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Log.Level}))
	slog.SetDefault(log)

	// ****************************************************
	// * Setup signaling server
	// ****************************************************
	slog.Info("setting up signaling server...")
	server := signaling.SetupServer(ctx, cfg.Signaling, signaling.SetupHandler(log))
	serverError := make(chan error, 1)
	go func() {
		serverError <- server.ListenAndServe()
	}()

	// ****************************************************
	// * Shutdown
	// ****************************************************
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverError:
		return fmt.Errorf("failed to start signaling server: %w", err)
	case <-done:
		slog.Info("shutting down signaling server...")
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown signaling server: %w", err)
		}
	}

	slog.Info("all good!")
	return nil
}
