package main

import (
	"log"
	"tj_techtest/config"
	"tj_techtest/database/seeders"
)

func main() {
	config.ConnectDB()
	log.Println("Seeding vehicles and locations...")
	seeders.SeedVehicles()
	log.Println("Seeding completed.")
}
