package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		"vehicle.locations", // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Create test location message
	location := LocationMessage{
		VehicleID: "3", // sesuaikan dengan ID vehicle yang ada
		Latitude:  -6.2088,
		Longitude: 106.8456,
		Timestamp: time.Now().Unix(),
	}

	body, err := json.Marshal(location)
	failOnError(err, "Failed to marshal JSON")

	// Publish message
	err = ch.Publish(
		"vehicle.locations", // exchange
		"vehicle.location",  // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
