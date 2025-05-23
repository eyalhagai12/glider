package main

import (
	"context"
	"glider/api"
	"glider/nodeapi"
	"glider/nodes"
	"glider/workerpool"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	logger.Info("Starting glider node...")
	logger.Info("Loading environment variables...")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	logger.Info("Environment variables loaded successfully.")

	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	dockerCli.Ping(ctx)
	defer dockerCli.Close()

	logger.Info("Registering node...")
	_, err = nodes.RegisterNode("http://172.18.0.7:8080")
	if err != nil {
		log.Fatalf("Error registering node: %v", err)
	}
	logger.Info("Node registered successfully.")
	logger.Info("Starting glider node...")
	logger.Info("Starting worker pool...")

	wp := workerpool.NewWorkerPool(10, 10)
	wp.Run(ctx)
	defer wp.Close()

	logger.Info("Worker pool started successfully.")

	nodeDeploymentHandler := nodeapi.NodeAgentHandlers(wp, dockerCli, logger)

	logger.Info("Starting HTTP server...")
	r := gin.Default()
	r.POST("/deploy", api.HandlerFromFunc(nodeDeploymentHandler.Deploy, http.StatusCreated))
	r.GET("/metrics", api.HandlerFromFunc(nodeDeploymentHandler.ReportMetrics, http.StatusOK))
	r.POST("/connect", api.HandlerFromFunc(nodeDeploymentHandler.ConnectToVPN, http.StatusOK))

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
