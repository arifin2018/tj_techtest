package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tj_techtest/config"
	"tj_techtest/pkg/rabbitmq"
	"tj_techtest/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	// Initialize Database
	config.ConnectDB()

	// Initialize RabbitMQ client
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	rmq, err := rabbitmq.NewClient(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Create context that will be canceled on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming location updates
	go rmq.ConsumeLocationUpdates(ctx)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Fleet Management API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel() // Cancel context to stop RabbitMQ consumer
		app.Shutdown()
	}()

	// Start server
	log.Fatal(app.Listen(":3000"))
}
