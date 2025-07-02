package controllers

import (
	"strconv"
	"tj_techtest/app/models"
	"tj_techtest/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GeofenceController struct{}

type CreateGeofenceRequest struct {
	Name      string  `json:"name" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Radius    float64 `json:"radius" validate:"required,min=1"`
}

var validate = validator.New()

// GetGeofences returns all geofences
func (c *GeofenceController) GetGeofences(ctx *fiber.Ctx) error {
	var geofences []models.Geofence
	result := config.DB.Find(&geofences)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting geofences",
			"error":   result.Error.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Geofences retrieved successfully",
		"data":    geofences,
	})
}

// CreateGeofence creates a new geofence
func (c *GeofenceController) CreateGeofence(ctx *fiber.Ctx) error {
	var req CreateGeofenceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	geofence := models.Geofence{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Radius:    req.Radius,
	}

	result := config.DB.Create(&geofence)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating geofence",
			"error":   result.Error.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Geofence created successfully",
		"data":    geofence,
	})
}

// GetGeofence returns a specific geofence
func (c *GeofenceController) GetGeofence(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	geofenceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid geofence ID",
		})
	}

	var geofence models.Geofence
	result := config.DB.First(&geofence, geofenceID)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Geofence not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Geofence retrieved successfully",
		"data":    geofence,
	})
}

// UpdateGeofence updates a geofence
func (c *GeofenceController) UpdateGeofence(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	geofenceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid geofence ID",
		})
	}

	var req CreateGeofenceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	var geofence models.Geofence
	result := config.DB.First(&geofence, geofenceID)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Geofence not found",
		})
	}

	geofence.Name = req.Name
	geofence.Latitude = req.Latitude
	geofence.Longitude = req.Longitude
	geofence.Radius = req.Radius

	result = config.DB.Save(&geofence)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating geofence",
			"error":   result.Error.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Geofence updated successfully",
		"data":    geofence,
	})
}

// DeleteGeofence deletes a geofence
func (c *GeofenceController) DeleteGeofence(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	geofenceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid geofence ID",
		})
	}

	result := config.DB.Delete(&models.Geofence{}, geofenceID)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting geofence",
			"error":   result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Geofence not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Geofence deleted successfully",
	})
}
