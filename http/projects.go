package http

import (
	"crypto/rand"
	"encoding/binary"
	backend "glider"
	"net"
	"net/http"

	"github.com/eyalhagai12/hagio/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterProjectsHandler() {
	routeGroup := s.apiRoutes.Group("/projects")
	routeGroup.POST("/", handler.FromFunc(s.CreateProject, http.StatusCreated))
}

func (s *Server) CreateProject(c *gin.Context, request createProject) (*backend.Project, error) {
	logger := s.logger.With("project_name", request.Name)

	project := &backend.Project{
		ID:   uuid.New(),
		Name: request.Name,
	}

	createdProject, err := s.projectService.Create(c.Request.Context(), project)
	if err != nil {
		return nil, err
	}

	privateKey, _, err := s.networkService.GenerateKeys()
	if err != nil {
		logger.Error("failed to generate keys", "error", err)
		return nil, err
	}

	port, err := getFreePort()
	if err != nil {
		logger.Error("failed to get free port", "error", err)
		return nil, err
	}
	network := backend.NewNetwork(project.Name, project.ID, 0, generateRandomIPv4(), port, privateKey)
	projectNetwork, err := s.networkService.Create(c.Request.Context(), network)
	if err != nil {
		return nil, err
	}
	logger.Debug("created network", "network_id", projectNetwork.ID)

	return createdProject, nil
}

func generateRandomIPv4() backend.IPAddress {
	bytes := make([]byte, 4)
	rand.Read(bytes[:3])
	bytes[3] = 1
	return backend.IPAddress(binary.BigEndian.Uint32(bytes))
}

func getFreePort() (uint16, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return uint16(addr.Port), nil
}
