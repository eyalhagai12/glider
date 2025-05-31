package main

import (
	"glider/http"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")

	server := http.NewServer("0.0.0.0", "8080")

	server.RegisterUserRoutes()
	server.RegisterDeploymentsHandler()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
