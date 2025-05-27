package main

import (
	"glider/http"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	server := http.NewServer("0.0.0.0", "8080")

	server.RegisterUserRoutes()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
