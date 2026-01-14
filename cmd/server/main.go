package main

import (
	"log"
	"net/http"

	"github.com/anandabhimanyu/employee-backend-api/internal/config"
	"github.com/anandabhimanyu/employee-backend-api/internal/db"
)

func main() {
	cfg := config.Load()

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	log.Println("Starting employee-backend-api...")

	postgres, err := db.NewPostgres(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer postgres.DB.Close()

	log.Println("Connected to PostgreSQL database")

	server := &http.Server{
		Addr: ":" + cfg.AppPort,
	}

	log.Println("Server running on port", cfg.AppPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
