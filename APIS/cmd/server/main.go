package main

import (
	"log"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/configs"
)

func main() {
	// Load database configuration
	_, err := configs.LoadDbConfig(".")
	if err != nil {
		log.Fatalf("failed to load database config: %v", err)
	}

	// Load web server configuration
	_, err = configs.LoadWebConfig(".")
	if err != nil {
		log.Fatalf("failed to load web config: %v", err)
	}

	// Get DSN for logging (optional)
	dsn, err := configs.GetDSN()
	if err != nil {
		log.Printf("warning: could not get DSN: %v", err)
	} else {
		log.Printf("Database DSN configured: %s", dsn)
	}

	// Initialize database connection
	db, err := configs.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection established")

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Configuration loaded successfully
	log.Println("Configuration loaded successfully")
}
