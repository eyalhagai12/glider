package main

import (
	"fmt"
	"glider/api"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()

	db, err := sqlx.Connect("postgres", "user=glider password=glider123 dbname=glider sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	deployHandlers := api.NewDeployHandlers(db)

	r.POST("/deploy", api.HandlerFromFunc(deployHandlers.Deploy, http.StatusAccepted))
	r.Run()
}
