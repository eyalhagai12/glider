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

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	dockerCli.Ping(ctx)
	defer dockerCli.Close()

	err = nodes.RegisterNode("http://localhost:8080")
	if err != nil {
		log.Fatalf("Error registering node: %v", err)
	}

	wp := workerpool.NewWorkerPool(10, 10)
	wp.Run(ctx)
	defer wp.Close()

	nodeDeploymentHandler := nodeapi.NewNodeDeploymentHandlers(wp, dockerCli)

	r := gin.Default()
	r.POST("/deploy", api.HandlerFromFunc(nodeDeploymentHandler.Deploy, http.StatusAccepted))
	r.GET("/metrics", api.HandlerFromFunc(nodeDeploymentHandler.ReportMetrics, http.StatusOK))

	if err := r.Run(); err != nil {
		panic(err)
	}
}
