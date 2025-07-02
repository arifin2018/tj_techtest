package seeders

import (
	"log"
	"time"
	"tj_techtest/app/models"
	"tj_techtest/config"
)

func SeedVehicles() {
	vehicles := []models.Vehicle{
		{Name: "Vehicle 1", Latitude: -6.200000, Longitude: 106.816666, LastSeen: time.Now()},
		{Name: "Vehicle 2", Latitude: -6.210000, Longitude: 106.826666, LastSeen: time.Now()},
	}

	for _, v := range vehicles {
		result := config.DB.Create(&v)
		if result.Error != nil {
			log.Printf("Failed to seed vehicle: %v", result.Error)
		} else {
			// Seed some locations for each vehicle
			locations := []models.VehicleLocation{
				{VehicleID: v.ID, Latitude: v.Latitude, Longitude: v.Longitude, Timestamp: time.Now().Add(-24 * time.Hour)},
				{VehicleID: v.ID, Latitude: v.Latitude + 0.001, Longitude: v.Longitude + 0.001, Timestamp: time.Now().Add(-12 * time.Hour)},
				{VehicleID: v.ID, Latitude: v.Latitude + 0.002, Longitude: v.Longitude + 0.002, Timestamp: time.Now()},
			}
			for _, loc := range locations {
				if err := config.DB.Create(&loc).Error; err != nil {
					log.Printf("Failed to seed vehicle location: %v", err)
				}
			}
		}
	}
}
