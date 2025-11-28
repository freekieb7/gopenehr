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
	"github.com/freekieb7/gopenehr/internal/cli"
	"github.com/freekieb7/gopenehr/internal/config"
	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/health"
	"github.com/freekieb7/gopenehr/internal/oauth"
	"github.com/freekieb7/gopenehr/internal/openehr"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/internal/webhook"
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
	settings := config.NewSettings()
	if err := settings.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Init logger
	logger := telemetry.NewLogger(settings.LogLevel)

	// Init database
	db := database.New()
	if err := db.Connect(ctx, settings.DatabaseURL); err != nil {
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

	// Audit logger
	auditLogger := audit.NewLogger(&db)
	webhookSaver := webhook.NewSaver(logger, &db)
	webhookSender := webhook.NewSender(logger, &db, &http.Client{
		Timeout: 10 * time.Second,
	})

	// Services
	healthChecker := health.NewChecker(settings.Version, config.TARGET_MIGRATION_VERSION, &db)
	auditService := audit.NewService(logger, &db)
	webhookService := webhook.NewService(logger, &db)
	oauthService := oauth.NewService(logger, settings.OAuthTrustedIssuers, settings.OAuthAudience)
	openEHRService := openehr.NewService(logger, &db)

	// Routes
	healthHandler := health.NewHandler(&healthChecker)
	healthHandler.RegisterRoutes(srv)

	auditHandler := audit.NewHandler(&settings, logger, &auditService, &oauthService, &auditLogger)
	auditHandler.RegisterRoutes(srv)

	webhookHandler := webhook.NewHandler(&settings, logger, &auditLogger, &oauthService, &webhookService)
	webhookHandler.RegisterRoutes(srv)

	openEHRHandler := openehr.NewHandler(&settings, logger, &openEHRService, &auditService, &webhookService, &auditLogger, &webhookSaver)
	openEHRHandler.RegisterRoutes(srv)

	// Set up signal handling for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// Create a channel for server errors
	serverErrChan := make(chan error, 1)

	// Warmup cache
	go func() {
		logger.InfoContext(ctx, "Warming up cache")
		if err := oauthService.WarmupCache(ctx); err != nil {
			logger.ErrorContext(ctx, "Cache warmup failed", "error", err)
		} else {
			logger.InfoContext(ctx, "Cache warmup completed")
		}
	}()

	// Start audit logger
	go func() {
		auditLogger.Start(ctx)
	}()

	// Start server
	go func() {
		logger.InfoContext(ctx, "Starting server", "port", settings.Port)
		if err := srv.Listen(":" + settings.Port); err != nil {
			serverErrChan <- err
		}
	}()

	webhookSaver.Start(ctx)

	// Start webhook delivery worker
	go func() {

		logger.InfoContext(ctx, "Starting webhook delivery worker")
		err := webhookSender.Run(ctx)
		if err != nil {
			logger.ErrorContext(ctx, "Webhook delivery worker error", "error", err)
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
	settings := config.NewSettings()
	if err := settings.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Init logger
	logger := telemetry.NewLogger(settings.LogLevel)

	// Init database
	db := database.New()
	if err := db.Connect(ctx, settings.DatabaseURL); err != nil {
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
	settings := config.NewSettings()
	if err := settings.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := telemetry.NewLogger(settings.LogLevel)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Check health endpoint
	url := fmt.Sprintf("http://localhost:%s/health/readyz", settings.Port)
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
