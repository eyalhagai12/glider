package main

import (
	"glider/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	deployHandlers := api.NewDeployHandlers()

	r.POST("/deploy", api.HandlerRFromFunc(deployHandlers.Deploy))
	r.Run()
}
