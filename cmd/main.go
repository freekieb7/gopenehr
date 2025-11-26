package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freekieb7/gopenehr/internal/audit"
	auditHandler "github.com/freekieb7/gopenehr/internal/audit/handler"
	"github.com/freekieb7/gopenehr/internal/cli"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/health"
	healthHandler "github.com/freekieb7/gopenehr/internal/health/handler"
	openehrHandler "github.com/freekieb7/gopenehr/internal/openehr/handler"
	openehrService "github.com/freekieb7/gopenehr/internal/openehr/service"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
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
		fmt.Println("Usage: gopenehr [command]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  serve          - Start the web server")
		fmt.Println("  migrate [cmd]  - Run database migrations (up/down)")
		fmt.Println("  healthcheck    - Check if the server is healthy")
		fmt.Println("  version        - Show version information")
		return nil
	}

	switch args[0] {
	case "serve":
		return runServer(ctx)
	case "migrate":
		return runMigrate(ctx, args[1:])
	case "healthcheck":
		return runHealthcheck(ctx)
	case "version":
		fmt.Printf("Version: %s\n", config.Version)
		return nil
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runServer(ctx context.Context) error {
	// Load config
	cfg := config.Config{}
	if err := cfg.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Init logger
	logger := telemetry.NewLogger(cfg.LogLevel)

	// Init database
	db := database.New()
	if err := db.Connect(ctx, cfg.DatabaseURL); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Init server
	srv := fiber.New(fiber.Config{
		// Reduce memory allocations
		DisableStartupMessage: true,
		// Network performance
		Network: "tcp4",

		// Increase read/write buffers
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,

		// Connection limits
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,

		// Reduce memory footprint
		ReduceMemoryUsage: true,

		// Enable SO_REUSEPORT for better load balancing across CPU cores
		EnableSplittingOnParsers: true,
	})

	// Compression middleware
	srv.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Services
	healthChecker := health.NewChecker(cfg.Version, config.TARGET_MIGRATION_VERSION, &db)
	auditService := audit.NewService(logger, &db)
	ehrService := openehrService.NewEHRService(logger, &db)
	demographicService := openehrService.NewDemographicService(logger, &db)
	queryService := openehrService.NewQueryService(logger, &db, &auditService)

	// Routes
	healthHandler := healthHandler.NewHandler(&healthChecker)
	healthHandler.RegisterRoutes(srv)

	auditHandler := auditHandler.NewHandler(logger, &auditService)
	auditHandler.RegisterRoutes(srv)

	openEHRHandler := openehrHandler.NewHandler(&cfg, logger, &ehrService, &demographicService, &queryService, &auditService)
	openEHRHandler.RegisterRoutes(srv)

	// Set up signal handling for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// Create a channel for server errors
	serverErrChan := make(chan error, 1)

	// Start server
	go func() {
		logger.InfoContext(ctx, "Starting server", "port", cfg.Port)
		if err := srv.Listen(":" + cfg.Port); err != nil {
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

	if err := srv.ShutdownWithContext(ctx); err != nil {
		logger.ErrorContext(ctx, "Failed to shutdown server", "error", err)
	} else {
		logger.InfoContext(ctx, "Server shutdown completed")
	}

	return nil
}

func runMigrate(ctx context.Context, args []string) error {
	// Load config
	cfg := config.Config{}
	if err := cfg.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Init logger
	logger := telemetry.NewLogger(cfg.LogLevel)

	// Init database
	db := database.Database{}
	if err := db.Connect(ctx, cfg.DatabaseURL); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	migrate := cli.Migrator{
		DB:            &db,
		Logger:        logger,
		MigrationsDir: "./migrations",
	}

	return migrate.Run(ctx, args)
}

func runHealthcheck(ctx context.Context) error {
	// Get port from environment variable (default: 3000)
	cfg := config.Config{}
	if err := cfg.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := telemetry.NewLogger(cfg.LogLevel)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Check health endpoint
	url := fmt.Sprintf("http://localhost:%s/health/readyz", cfg.Port)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.ErrorContext(ctx, "Failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unhealthy: status code %d", resp.StatusCode)
	}

	fmt.Println("✓ Server is healthy")
	return nil
}
