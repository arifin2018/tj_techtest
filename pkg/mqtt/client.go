package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"tj_techtest/app/models"
	"tj_techtest/config"
	"tj_techtest/pkg/rabbitmq"

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

func isWithinGeofence(lat1, lon1, lat2, lon2, radius float64) bool {
	const earthRadius = 6371000 // meters

	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	lat1Rad := lat1 * (math.Pi / 180)
	lat2Rad := lat2 * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c

	return distance <= radius
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

		// Save location to database
		locationRecord := models.VehicleLocation{
			VehicleID: 0, // will set below
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
			Timestamp: time.Unix(location.Timestamp, 0),
		}

		// Get vehicle ID from database by vehicle code (vehicle_id string)
		var vehicle models.Vehicle
		if err := config.DB.Where("name = ?", location.VehicleID).First(&vehicle).Error; err != nil {
			log.Printf("Vehicle not found for code %s: %v", location.VehicleID, err)
			return
		}
		locationRecord.VehicleID = vehicle.ID

		if err := config.DB.Create(&locationRecord).Error; err != nil {
			log.Printf("Failed to save location: %v", err)
			return
		}

		log.Printf("Saved location for vehicle %s to database", location.VehicleID)

		// Check geofence
		var geofences []models.Geofence
		if err := config.DB.Find(&geofences).Error; err != nil {
			log.Printf("Failed to load geofences: %v", err)
			return
		}

		for _, geofence := range geofences {
			if isWithinGeofence(location.Latitude, location.Longitude, geofence.Latitude, geofence.Longitude, geofence.Radius) {
				log.Printf("Vehicle %s entered geofence %s", location.VehicleID, geofence.Name)

				// Publish event to RabbitMQ
				rabbitURL := "amqp://guest:guest@localhost:5672/"
				if envURL := os.Getenv("RABBITMQ_URL"); envURL != "" {
					rabbitURL = envURL
				}

				rabbitClient, err := rabbitmq.NewClient(rabbitURL)
				if err != nil {
					log.Printf("Failed to connect to RabbitMQ: %v", err)
					continue
				}
				defer rabbitClient.Close()

				event := rabbitmq.GeofenceEvent{
					VehicleID: location.VehicleID,
					Event:     "geofence_entry",
					Location: struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					}{
						Latitude:  location.Latitude,
						Longitude: location.Longitude,
					},
					Timestamp: time.Now().Unix(),
				}

				err = rabbitClient.PublishGeofenceEvent(event)
				if err != nil {
					log.Printf("Failed to publish geofence event: %v", err)
				} else {
					log.Printf("Published geofence event for vehicle %s", location.VehicleID)
				}
			}
		}
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
