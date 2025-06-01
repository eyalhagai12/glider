package http

import (
	"errors"
	backend "glider"
	"glider/docker"
	"glider/gitrepo"
	"log/slog"
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
	logger := s.logger.With("deployment_name", request.Name)

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
	defer func(deployment *backend.Deployment) {
		if err != nil {
			logger.Debug("rolling back deployment", slog.String("error", err.Error()), slog.Any("deployment", deployment))
			deployment.Status = backend.DeploymentStatusFailed
			_, err := s.deploymentService.Update(c.Request.Context(), deployment)
			if err != nil {
				logger.Error("failed to update deployment", slog.String("error", err.Error()))
			}
		}
	}(deployment)

	sourceCodeService := gitrepo.NewGitSourceCodeService()
	deployer, err := docker.NewLocalDeployer(s.db, logger, sourceCodeService)
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

	deployment.Status = backend.DeploymentStatusReady

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
