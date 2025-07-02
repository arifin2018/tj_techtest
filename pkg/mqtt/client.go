package mqtt

import (
	"encoding/json"
	"fmt"
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

type Client struct {
	client MQTT.Client
}

func NewClient(brokerURL string, clientID string) (*Client, error) {
	opts := MQTT.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(30 * time.Second)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SubscribeToLocations() error {
	topic := "/fleet/vehicle/+/location" // + is wildcard for vehicle_id
	token := c.client.Subscribe(topic, 0, func(client MQTT.Client, msg MQTT.Message) {
		var location LocationMessage
		if err := json.Unmarshal(msg.Payload(), &location); err != nil {
			log.Printf("Error decoding location message: %v", err)
			return
		}

		// Validate required fields
		if location.VehicleID == "" {
			log.Printf("Invalid message: missing vehicle_id")
			return
		}
		if location.Timestamp == 0 {
			log.Printf("Invalid message: missing timestamp")
			return
		}

		log.Printf("Received location update for vehicle %s: [%f, %f] at %v",
			location.VehicleID,
			location.Latitude,
			location.Longitude,
			time.Unix(location.Timestamp, 0),
		)

		// TODO: Save location to database
		// TODO: Check geofence
	})

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe: %v", token.Error())
	}

	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

func (c *Client) Close() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
	}
}
