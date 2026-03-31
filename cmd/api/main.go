package main

import (
	"go-auth-backend-api/app"
	"log"
	"net/http"
	"os"
)

func main() {
	h, err := app.Handler()
	if err != nil {
		log.Fatalf("❌ Application initialization failed: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("🚀 Server running on :%s", port)
	if err := http.ListenAndServe(":"+port, h); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
