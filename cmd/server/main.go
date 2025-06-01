package main

import (
	"glider/http"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server := http.NewServer("0.0.0.0", "8080")

	server.RegisterUserRoutes()
	server.RegisterDeploymentsHandler()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
