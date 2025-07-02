package models

import (
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type VehicleLocation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	VehicleID uint      `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	Vehicle   Vehicle   `json:"vehicle" gorm:"foreignKey:VehicleID"`
}

type Geofence struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	Radius    float64        `json:"radius"` // dalam meter
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
