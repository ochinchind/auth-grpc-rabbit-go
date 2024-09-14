package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"example/hello/internal/app"
	"example/hello/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	application.GRPCServer.MustRun()

	go func() {
		application.GRPCServer.MustRun()
	}()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Waiting for SIGINT (pkill -2) or SIGTERM
	<-stop

	// initiate graceful shutdown
	application.GRPCServer.Stop() // Assuming GRPCServer has Stop() method for graceful shutdown
	log.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var opts slog.HandlerOptions

	switch env {
	case envLocal, envDev:
		// Enable detailed logging with timestamps and more verbosity in dev and local environments
		opts = slog.HandlerOptions{
			Level: slog.LevelDebug, // Log everything down to debug level
		}
	case envProd:
		// Use more concise logging in production
		opts = slog.HandlerOptions{
			Level: slog.LevelInfo, // Only log important info, warnings, and errors
		}
	default:
		opts = slog.HandlerOptions{
			Level: slog.LevelInfo, // Default to info level
		}
	}

	// Create a new logger that writes to stdout (or any file if needed)
	handler := slog.NewTextHandler(os.Stdout, &opts)
	logger := slog.New(handler)

	// Optionally, log the environment as an initial log entry
	logger.Info("Logger initialized", slog.String("env", env))

	return logger
}
