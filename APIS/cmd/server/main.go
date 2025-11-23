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

	// Configuration loaded successfully
	log.Println("Configuration loaded successfully")
}
