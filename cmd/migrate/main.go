package main

import (
	"log"

	"go-auth-backend-api/migrations"
	"go-auth-backend-api/pkg/database"
)

func main() {

	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := migrations.Migrate(); err != nil {
		log.Fatal(err)
	}

	log.Println("Migration completed successfully")
}
