package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log/slog"
	"os"
)

type config struct {
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))

	if err := run(context.Background()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(_ context.Context) error {
	// ****************************************************
	// * Load environment variables
	// ****************************************************
	slog.Info("loading environment variables...")
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to process environment variables: %w", err)
	}

	slog.Info("all good!")
	return nil
}
