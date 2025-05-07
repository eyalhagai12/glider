package main

import (
	"context"
	"fmt"
	"glider/api"
	"glider/workerpool"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	os.Setenv("DOCKER_BUILDKIT", "1")

	db, err := sqlx.Connect("postgres", "user=glider password=glider123 dbname=glider sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	workerPool := workerpool.NewWorkerPool(10, 10)
	workerPool.Run(context.Background())

	deployHandlers := api.NewDeployHandlers(db, workerPool, cli)
	nodesHandlers := api.NewNodeHandlers(db)

	r.POST("/deploy", api.HandlerFromFunc(deployHandlers.Deploy, http.StatusAccepted))
	r.POST("/nodes", api.HandlerFromFunc(nodesHandlers.RegisterNewNode, http.StatusCreated))
	r.Run()
	workerPool.Close()
}
