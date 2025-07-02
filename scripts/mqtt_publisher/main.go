package main

import (
	"encoding/json"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	brokerURL := "tcp://localhost:1883"
	clientID := "fleet-management-publisher"

	opts := MQTT.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetCleanSession(true).
		SetKeepAlive(2 * time.Second).
		SetPingTimeout(1 * time.Second).
		SetConnectTimeout(1 * time.Second).
		SetAutoReconnect(true).
		SetOrderMatters(false)

	log.Printf("Connecting to MQTT broker at %s...", brokerURL)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	log.Printf("Successfully connected to MQTT broker!")
	defer client.Disconnect(250)

	vehicleID := "B1234XYZ"

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		location := LocationMessage{
			VehicleID: vehicleID,
			Latitude:  -6.2088,
			Longitude: 106.8456,
			Timestamp: time.Now().Unix(),
		}

		payload, err := json.Marshal(location)
		if err != nil {
			log.Printf("Failed to marshal location message: %v", err)
			continue
		}

		topic := "/fleet/vehicle/" + vehicleID + "/location"
		token := client.Publish(topic, 0, false, payload)
		token.Wait()

		if token.Error() != nil {
			log.Printf("Failed to publish message: %v", token.Error())
		} else {
			log.Printf("Published location to topic %s: %s", topic, string(payload))
		}
	}
}
