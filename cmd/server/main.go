package main

import (
	"log"
	"net/http"

	"github.com/anandabhimanyu/employee-backend-api/internal/config"
)

func main() {
	cfg := config.Load()

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	log.Println("Starting employee-backend-api...")

	server := &http.Server{
		Addr: ":" + cfg.AppPort,
	}

	log.Println("Server running on port", cfg.AppPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
