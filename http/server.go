package http

import (
	"database/sql"
	backend "glider"
	"glider/docker"
	"glider/pg"
	"log"
	"log/slog"
	"os"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

const apiRouteGroup = "api"

type Server struct {
	engine    *gin.Engine
	apiRoutes *gin.RouterGroup

	logger *slog.Logger
	db     *sql.DB

	Port string
	Host string

	userService       backend.UserService
	imageService      backend.ImageService
	deploymentService backend.DeploymentService
	containerService  backend.ContainerService
}

func NewServer(host string, port string) *Server {
	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.Use(gin.ErrorLogger())

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)

	server := &Server{
		engine:    engine,
		apiRoutes: engine.Group(apiRouteGroup),
		Port:      port,
		Host:      host,
		logger:    logger,
	}

	db, err := pg.NewDatabase()
	if err != nil {
		log.Fatal("Failed to connect to the database: " + err.Error())
	}
	server.db = db

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal("failed to connect to client")
	}

	userService := pg.NewUserService(db)
	containerService := docker.NewDockerContainerService(cli, db)
	imageService := docker.NewDockerImageService(cli, db)
	deploymentService := pg.NewDeploymentService(db)

	server.userService = userService
	server.containerService = containerService
	server.imageService = imageService
	server.deploymentService = deploymentService

	return server
}

func (s *Server) Start() error {
	if err := s.engine.Run(s.Host + ":" + s.Port); err != nil {
		return err
	}

	return nil
}
