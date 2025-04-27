package main

import (
	"glider/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	deployHandlers := api.NewDeployHandlers()

	r.POST("/deploy", api.HandlerFromFunc(deployHandlers.Deploy, http.StatusAccepted))
	r.Run()
}
