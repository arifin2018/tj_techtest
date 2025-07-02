package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tj_techtest/pkg/mqtt"
)

func main() {
	brokerURL := "tcp://localhost:1883"
	clientID := "fleet-management-subscriber"

	client, err := mqtt.NewClient(brokerURL, clientID)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}
	defer client.Close()

	err = client.SubscribeToLocations()
	if err != nil {
		log.Fatalf("Failed to subscribe to locations: %v", err)
	}

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down MQTT subscriber...")
}
