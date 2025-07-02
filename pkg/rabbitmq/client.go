package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"strconv"
	"time"
	"tj_techtest/app/models"
	"tj_techtest/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	LocationExchange = "vehicle.locations"
	GeofenceExchange = "fleet.events"
	LocationQueue    = "location.updates"
	GeofenceQueue    = "geofence_alerts"
)

type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

type GeofenceEvent struct {
	VehicleID string `json:"vehicle_id"`
	Event     string `json:"event"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Timestamp int64 `json:"timestamp"`
}

func (c *Client) ConsumeGeofenceAlerts(ctx context.Context) (<-chan GeofenceEvent, error) {
	msgs, err := c.channel.Consume(
		GeofenceQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	events := make(chan GeofenceEvent)

	go func() {
		defer close(events)
		log.Printf("Started consuming from queue: %s", GeofenceQueue)
		for msg := range msgs {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("Received message: %s", string(msg.Body))
				var event GeofenceEvent
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					log.Printf("Error decoding geofence event: %v", err)
					continue
				}
				log.Printf("Successfully decoded event: %+v", event)
				events <- event
			}
		}
	}()

	return events, nil
}

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewClient(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		channel: ch,
	}

	err = client.setupExchangesAndQueues()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) setupExchangesAndQueues() error {
	// Declare exchanges
	err := c.channel.ExchangeDeclare(
		LocationExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.channel.ExchangeDeclare(
		GeofenceExchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare queues
	_, err = c.channel.QueueDeclare(
		LocationQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.channel.QueueDeclare(
		GeofenceQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bind queues to exchanges
	err = c.channel.QueueBind(
		LocationQueue,
		"#",
		LocationExchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(
		GeofenceQueue,
		"",
		GeofenceExchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ConsumeLocationUpdates(ctx context.Context) {
	msgs, err := c.channel.Consume(
		LocationQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return
	}

	go func() {
		for msg := range msgs {
			select {
			case <-ctx.Done():
				return
			default:
				var locMsg LocationMessage
				err := json.Unmarshal(msg.Body, &locMsg)
				if err != nil {
					log.Printf("Error unmarshaling location message: %v", err)
					continue
				}

				// Convert vehicle_id string to uint
				vehicleID, err := strconv.ParseUint(locMsg.VehicleID, 10, 32)
				if err != nil {
					log.Printf("Error parsing vehicle_id: %v", err)
					continue
				}

				// Save location to database
				location := models.VehicleLocation{
					VehicleID: uint(vehicleID),
					Latitude:  locMsg.Latitude,
					Longitude: locMsg.Longitude,
					Timestamp: time.Unix(locMsg.Timestamp, 0),
				}

				if err := config.DB.Create(&location).Error; err != nil {
					log.Printf("Error saving location: %v", err)
					continue
				}

				// Check geofence
				c.checkGeofence(ctx, location)
			}
		}
	}()

	log.Printf("Location consumer started. Waiting for messages...")
}

func (c *Client) checkGeofence(ctx context.Context, location models.VehicleLocation) {
	var geofences []models.Geofence
	if err := config.DB.Find(&geofences).Error; err != nil {
		log.Printf("Error fetching geofences: %v", err)
		return
	}

	log.Printf("Checking %d geofences for vehicle %d at location [%f, %f]",
		len(geofences), location.VehicleID, location.Latitude, location.Longitude)

	for _, geofence := range geofences {
		if isInsideGeofence(location.Latitude, location.Longitude, geofence) {
			log.Printf("Vehicle %d is inside geofence %s! Publishing event...", location.VehicleID, geofence.Name)

			event := GeofenceEvent{
				VehicleID: strconv.FormatUint(uint64(location.VehicleID), 10),
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

			eventJSON, err := json.Marshal(event)
			if err != nil {
				log.Printf("Error marshaling geofence event: %v", err)
				continue
			}

			err = c.channel.PublishWithContext(ctx,
				GeofenceExchange,
				"",
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        eventJSON,
				},
			)
			if err != nil {
				log.Printf("Error publishing geofence event: %v", err)
			} else {
				log.Printf("Successfully published geofence event: %s", string(eventJSON))
			}
		}
	}
}

// Simple circular geofence check using Haversine formula
func isInsideGeofence(lat, lon float64, geofence models.Geofence) bool {
	const earthRadius = 6371000 // Earth radius in meters

	// Convert latitude and longitude to radians
	lat1 := lat * (math.Pi / 180)
	lon1 := lon * (math.Pi / 180)
	lat2 := geofence.Latitude * (math.Pi / 180)
	lon2 := geofence.Longitude * (math.Pi / 180)

	// Calculate differences
	dlat := lat2 - lat1
	dlon := lon2 - lon1

	// Calculate distance using Haversine formula
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	log.Printf("Distance to geofence %s: %.2f meters (radius: %.2f meters)",
		geofence.Name, distance, geofence.Radius)

	return distance <= geofence.Radius
}

func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) PublishGeofenceEvent(event GeofenceEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		GeofenceExchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
	return err
}
