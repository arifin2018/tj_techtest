package controllers

import (
	"strconv"
	"tj_techtest/app/models"
	"tj_techtest/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type VehicleController struct{}

type CreateVehicleRequest struct {
	Name string `json:"name" validate:"required"`
}

var vehicleValidator = validator.New()

// GetVehicles returns all vehicles
func (c *VehicleController) GetVehicles(ctx *fiber.Ctx) error {
	var vehicles []models.Vehicle
	result := config.DB.Find(&vehicles)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting vehicles",
			"error":   result.Error.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Vehicles retrieved successfully",
		"data":    vehicles,
	})
}

// CreateVehicle creates a new vehicle
func (c *VehicleController) CreateVehicle(ctx *fiber.Ctx) error {
	var req CreateVehicleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate request
	if err := vehicleValidator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	vehicle := models.Vehicle{
		Name: req.Name,
	}

	result := config.DB.Create(&vehicle)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating vehicle",
			"error":   result.Error.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Vehicle created successfully",
		"data":    vehicle,
	})
}

// GetVehicle returns a specific vehicle
func (c *VehicleController) GetVehicle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	vehicleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid vehicle ID",
		})
	}

	var vehicle models.Vehicle
	result := config.DB.First(&vehicle, vehicleID)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Vehicle not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Vehicle retrieved successfully",
		"data":    vehicle,
	})
}

// GetVehicleLocations returns location history for a vehicle within a time range
func (c *VehicleController) GetVehicleLocations(ctx *fiber.Ctx) error {
	vehicleID := ctx.Params("id")

	// Parse query parameters
	startTimestamp, err := strconv.ParseInt(ctx.Query("start", "0"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start timestamp",
		})
	}

	endTimestamp, err := strconv.ParseInt(ctx.Query("end", "0"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end timestamp",
		})
	}

	// Query builder
	query := config.DB.Where("vehicle_id = ?", vehicleID)

	// Filter berdasarkan rentang waktu
	if startTimestamp > 0 {
		query = query.Where("EXTRACT(EPOCH FROM timestamp) >= ?", startTimestamp)
	}
	if endTimestamp > 0 {
		query = query.Where("EXTRACT(EPOCH FROM timestamp) <= ?", endTimestamp)
	}

	// Query dengan sorting berdasarkan timestamp
	var locations []models.VehicleLocation
	result := query.Order("timestamp ASC").Find(&locations)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting location history",
		})
	}

	// Transform response
	var response []fiber.Map
	for _, loc := range locations {
		response = append(response, fiber.Map{
			"vehicle_id": vehicleID,
			"latitude":   loc.Latitude,
			"longitude":  loc.Longitude,
			"timestamp":  loc.Timestamp.Unix(),
		})
	}

	return ctx.JSON(response)
}

// GetLastLocation returns the last known location of a vehicle
func (c *VehicleController) GetLastLocation(ctx *fiber.Ctx) error {
	vehicleID := ctx.Params("id")
	var location models.VehicleLocation

	result := config.DB.Where("vehicle_id = ?", vehicleID).
		Order("timestamp DESC").
		First(&location)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"vehicle_id": vehicleID,
		"latitude":   location.Latitude,
		"longitude":  location.Longitude,
		"timestamp":  location.Timestamp.Unix(),
	})
}
