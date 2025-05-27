package http

import (
	backend "glider"
	"glider/pg"
	"log"

	"github.com/gin-gonic/gin"
)

const apiRouteGroup = "api"

type Server struct {
	engine    *gin.Engine
	apiRoutes *gin.RouterGroup

	Port string
	Host string

	userService backend.UserService
}

func NewServer(host string, port string) *Server {
	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.Use(gin.ErrorLogger())

	server := &Server{
		engine:    engine,
		apiRoutes: engine.Group(apiRouteGroup),
		Port:      port,
		Host:      host,
	}

	db, err := pg.NewDatabase()
	if err != nil {
		log.Fatal("Failed to connect to the database: " + err.Error())
	}

	userService := pg.NewUserService(db)
	server.userService = userService

	return server
}

func (s *Server) Start() error {
	if err := s.engine.Run(s.Host + ":" + s.Port); err != nil {
		return err
	}

	return nil
}
