package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting employee-backend-api...")

	server := &http.Server{
		Addr: ":8080",
	}

	log.Println("Server running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
