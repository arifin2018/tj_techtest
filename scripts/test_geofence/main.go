package main

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare the exchange
	err = ch.ExchangeDeclare(
		"vehicle.locations", // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	// Create test location data (inside geofence area)
	location := LocationMessage{
		VehicleID: "3",
		Latitude:  -6.2088,  // Koordinat dalam radius geofence
		Longitude: 106.8456, // Koordinat dalam radius geofence
		Timestamp: time.Now().Unix(),
	}

	body, err := json.Marshal(location)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}

	// Publish the message
	err = ch.Publish(
		"vehicle.locations", // exchange
		"vehicle.location",  // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	log.Printf("Sent location update: %s", body)
}
