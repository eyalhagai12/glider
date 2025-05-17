package main

import (
	"context"
	"glider/api"
	"glider/nodeapi"
	"glider/nodes"
	"glider/workerpool"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	logger := log.Default()
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	logger.SetPrefix("glider: ")
	logger.Println("Starting glider node...")
	logger.Println("Loading environment variables...")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	logger.Println("Environment variables loaded successfully.")

	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	dockerCli.Ping(ctx)
	defer dockerCli.Close()

	logger.Println("Registering node...")
	_, err = nodes.RegisterNode("http://localhost:8080")
	if err != nil {
		log.Fatalf("Error registering node: %v", err)
	}
	logger.Println("Node registered successfully.")
	logger.Println("Starting glider node...")
	logger.Println("Starting worker pool...")

	wp := workerpool.NewWorkerPool(10, 10)
	wp.Run(ctx)
	defer wp.Close()

	logger.Println("Worker pool started successfully.")

	nodeDeploymentHandler := nodeapi.NewNodeDeploymentHandlers(wp, dockerCli)

	logger.Println("Starting HTTP server...")
	r := gin.Default()
	r.POST("/deploy", api.HandlerFromFunc(nodeDeploymentHandler.Deploy, http.StatusCreated))
	r.GET("/metrics", api.HandlerFromFunc(nodeDeploymentHandler.ReportMetrics, http.StatusOK))

	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
