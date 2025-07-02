package routes

import (
	"tj_techtest/app/http/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Initialize controllers
	vehicleController := &controllers.VehicleController{}
	geofenceController := &controllers.GeofenceController{}

	// Vehicle routes
	vehicles := app.Group("/vehicles")
	vehicles.Get("/", vehicleController.GetVehicles)
	vehicles.Post("/", vehicleController.CreateVehicle)
	vehicles.Get("/:id", vehicleController.GetVehicle)
	vehicles.Get("/:id/history", vehicleController.GetVehicleLocations)
	vehicles.Get("/:id/location", vehicleController.GetLastLocation)

	// Geofence routes
	geofences := app.Group("/geofences")
	geofences.Get("/", geofenceController.GetGeofences)
	geofences.Post("/", geofenceController.CreateGeofence)
	geofences.Get("/:id", geofenceController.GetGeofence)
	geofences.Put("/:id", geofenceController.UpdateGeofence)
	geofences.Delete("/:id", geofenceController.DeleteGeofence)
}
