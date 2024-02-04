package main

import (
	"fmt"
	"log"
	"p2p/initializers"
	"p2p/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	initializers.DB.AutoMigrate(&models.User{}, &models.OTPCode{})
	fmt.Println("👍 Migration complete")
}
