package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tj_techtest/pkg/rabbitmq"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	// Get RabbitMQ URL from environment
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	// Connect to RabbitMQ
	rmq, err := rabbitmq.NewClient(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Create context that will be canceled on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming geofence alerts
	alerts, err := rmq.ConsumeGeofenceAlerts(ctx)
	if err != nil {
		log.Fatalf("Failed to start consuming geofence alerts: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	// Process geofence alerts
	log.Println("Geofence worker started. Waiting for alerts...")
	for alert := range alerts {
		log.Printf("Received geofence alert: Vehicle %s entered geofence at location [%f, %f]",
			alert.VehicleID,
			alert.Location.Latitude,
			alert.Location.Longitude,
		)
	}
}
