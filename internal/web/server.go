package web

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

type Server struct {
	Fiber *fiber.App
}

func NewServer() Server {
	return Server{
		Fiber: fiber.New(fiber.Config{
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
		}),
	}
}

func (s *Server) ListenAndServe(addr string) error {
	if err := s.Fiber.Listen(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.Fiber.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

func (s *Server) EnableCompression() {
	s.Fiber.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Faster than LevelBestCompression
	}))
}
