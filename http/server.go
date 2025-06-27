package http

import (
	"database/sql"
	backend "glider"
	"glider/docker"
	"glider/pg"
	"glider/wireguard"
	"log"
	"log/slog"
	"os"

	"github.com/alitto/pond/v2"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

const apiRouteGroup = "api"

type Server struct {
	engine    *gin.Engine
	apiRoutes *gin.RouterGroup

	logger     *slog.Logger
	db         *sql.DB
	workerpool pond.Pool

	Port string
	Host string

	userService       backend.UserService
	imageService      backend.ImageService
	deploymentService backend.DeploymentService
	containerService  backend.ContainerService
	sourceCodeService backend.SourceCodeService
	projectService    backend.ProjectService
	networkService    backend.NetworkService
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

	workerpool := pond.NewPool(10, pond.WithQueueSize(1000), pond.WithNonBlocking(true))

	server := &Server{
		engine:     engine,
		apiRoutes:  engine.Group(apiRouteGroup),
		Port:       port,
		Host:       host,
		logger:     logger,
		workerpool: workerpool,
	}

	db, err := pg.NewDatabase(logger)
	if err != nil {
		log.Fatal("Failed to connect to the database: " + err.Error())
	}
	server.db = db

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		log.Fatal("failed to connect to client")
	}

	userService := pg.NewUserService(db)
	containerService := docker.NewDockerContainerService(cli, db, logger)
	imageService := docker.NewDockerImageService(cli, db, logger)
	deploymentService := pg.NewDeploymentService(db)
	projectService := pg.NewProjectService(db)
	networkService := wireguard.NewWireGuardService(db, logger)

	server.userService = userService
	server.containerService = containerService
	server.imageService = imageService
	server.deploymentService = deploymentService
	server.projectService = projectService
	server.networkService = networkService

	return server
}

func (s *Server) Start() error {
	if err := s.engine.Run(s.Host + ":" + s.Port); err != nil {
		return err
	}

	return nil
}
