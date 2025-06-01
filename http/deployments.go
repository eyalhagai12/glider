package http

import (
	"errors"
	backend "glider"
	"glider/docker"
	"glider/gitrepo"
	"net/http"

	"github.com/eyalhagai12/hagio/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterDeploymentsHandler() {
	routeGroup := s.apiRoutes.Group("/deployments")
	routeGroup.POST("/", handler.FromFunc(s.createDeployment, http.StatusAccepted))
}

func (s *Server) createDeployment(c *gin.Context, request deployRequest) (*backend.Deployment, error) {
	tags := backend.TagsFromList(request.Tags)
	deployment := &backend.Deployment{
		ID:             uuid.New(),
		Name:           request.Name,
		ProjectID:      request.ProjectID,
		Environment:    request.Environment,
		Tags:           tags,
		DeployMetadata: request.DeployMetadata,
		Status:         backend.DeploymentStatusPending,
		Version:        request.Version,
		Replicas:       request.Replicas,
	}

	if err := request.DeployMetadata.Validate(); err != nil {
		return nil, errors.Join(err, errors.New("invalid deploy metadata"))
	}

	deployment, err := s.deploymentService.Create(c.Request.Context(), deployment)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to create deployment"))
	}

	sourceCodeService := gitrepo.NewGitSourceCodeService()
	deployer, err := docker.NewLocalDeployer(s.db, s.logger, sourceCodeService)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to create deployer"))
	}

	image := &backend.Image{
		ID:          uuid.New(),
		Name:        request.Name,
		Version:     request.Version,
		RegistryURL: "localhost:5000",
	}

	deployment.ImageID = image.ID

	deployment, err = deployer.Deploy(c.Request.Context(), deployment, image, request.DeployMetadata)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to deploy"))
	}

	deployment, err = s.deploymentService.Update(c.Request.Context(), deployment)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to update deployment"))
	}

	deployment, err = s.deploymentService.GetByID(c.Request.Context(), deployment.ID)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to get deployment"))
	}

	return deployment, nil
}
