package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	internal "github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/health"
	"github.com/freekieb7/gopenehr/internal/web"
	"github.com/freekieb7/gopenehr/internal/web/handler"
	_ "go.uber.org/automaxprocs/maxprocs"
)

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: ")
		fmt.Println("  [command]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  serve          - Start the web server")
		// fmt.Println("  migrate [cmd]  - Run database migrations (up/down)")
		return nil
	}

	switch args[0] {
	case "serve":
		return runServer(ctx)
	// case "migrate":
	// 	return runMigrate(ctx, args[1:])
	default:
		return nil
	}

}

func runServer(ctx context.Context) error {
	// Load config
	cfg := internal.Config{}
	if err := cfg.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Init logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Init server
	srv := web.NewServer()

	// Compression middleware (use in production)
	if cfg.Environment == internal.Production {
		srv.EnableCompression()
	}

	// Routes
	healthHandler := handler.Health{
		HealthChecker: &health.Checker{
			Version: cfg.Version,
		},
	}
	healthHandler.RegisterRoutes(&srv)

	// Set up signal handling for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// Create a channel for server errors
	serverErrChan := make(chan error, 1)

	// Start server
	go func() {
		if err := srv.ListenAndServe(":3000"); err != nil {
			serverErrChan <- err
		}
	}()

	// Wait for termination signal or server error
	select {
	case sig := <-stopChan:
		logger.InfoContext(ctx, "Received shutdown signal", "signal", sig.String())
	case err := <-serverErrChan:
		if err != nil {
			logger.ErrorContext(ctx, "Server error occurred", "error", err)
			return err
		}
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.ErrorContext(ctx, "Failed to shutdown server", "error", err)
	} else {
		logger.InfoContext(ctx, "Server shutdown completed")
	}

	return nil
}
