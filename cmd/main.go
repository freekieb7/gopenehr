package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx := context.Background()
	args := os.Args[1:]

	if len(args) == 0 {
		println("Usage: gopenehr serve")
		return
	}

	switch args[0] {
	case "serve":
		if err := serve(ctx); err != nil {
			fmt.Println("Error:", err)
		}
	default:
		fmt.Println("Unknown command:", args[0])
	}
}

func serve(ctx context.Context) error {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	if err := app.Listen(":3000"); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	if err := app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
