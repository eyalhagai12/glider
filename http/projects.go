package http

import (
	backend "glider"
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

	projectNetwork, err := s.networkService.Create(c.Request.Context(), &backend.Network{
		ID:        uuid.New(),
		Name:      project.Name,
		ProjectID: project.ID,
	})
	if err != nil {
		return nil, err
	}
	logger.Debug("created network", "network_id", projectNetwork.ID)

	return createdProject, nil
}
