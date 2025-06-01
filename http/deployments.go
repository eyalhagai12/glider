package http

import (
	"context"
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
	routeGroup.GET("/:id", handler.FromFunc(s.getDeployment, http.StatusOK))
	routeGroup.POST("/", handler.FromFunc(s.createDeployment, http.StatusAccepted))
}

func (s *Server) getDeployment(c *gin.Context, _ any) (DeploymentResponse, error) {
	deploymentUuid := uuid.MustParse(c.Param("id"))
	deployment, err := s.deploymentService.GetByID(c.Request.Context(), deploymentUuid)
	if err != nil {
		return DeploymentResponse{}, err
	}

	if deployment.DeletedAt != nil {
		return DeploymentResponse{}, backend.ErrNotFound
	}

	actions := []DeploymentAction{}
	deploymentUrl := deployment.URL()

	if deployment.Status == backend.DeploymentStatusFailed {
		actions = append(actions, DeploymentAction{
			Name:   "retry",
			Method: "POST",
			Path:   deploymentUrl + "/retry",
		})
	}

	if deployment.Status == backend.DeploymentStatusPending {
		actions = append(actions, DeploymentAction{
			Name:   "cancel",
			Method: "POST",
			Path:   deploymentUrl + "/cancel",
		}) // to actually implement this i need to store the context cancel function and call when needed, this is a todo
	}

	if deployment.Status != backend.DeploymentStatusPending {
		actions = append(actions, DeploymentAction{
			Name:   "delete",
			Method: "DELETE",
			Path:   deploymentUrl,
		})

		actions = append(actions, DeploymentAction{
			Name:   "redeploy",
			Method: "POST",
			Path:   deploymentUrl + "/redeploy",
		})
	}

	return DeploymentResponse{
		Deployment: deployment,
		Actions:    actions,
		URL:        deployment.ProxyURL(),
	}, nil
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

	task := s.workerpool.SubmitErr(s.deployTask(context.Background(), deployment))

	go func() {
		if err := task.Wait(); err != nil {
			s.logger.Error("failed to deploy", slog.String("error", err.Error()), slog.Any("deployment", deployment))
		}
	}()

	return deployment, nil
}

func (s *Server) deployTask(ctx context.Context, deployment *backend.Deployment) func() error {
	return func() error {
		logger := s.logger.With("deployment_id", deployment.ID)

		deployment, err := s.deploymentService.Create(ctx, deployment)
		if err != nil {
			return errors.Join(err, errors.New("failed to create deployment"))
		}
		defer func(deployment *backend.Deployment) {
			if err != nil {
				logger.Debug("rolling back deployment", slog.String("error", err.Error()), slog.Any("deployment", deployment))
				deployment.Status = backend.DeploymentStatusFailed
				_, err := s.deploymentService.Update(ctx, deployment)
				if err != nil {
					logger.Error("failed to update deployment", slog.String("error", err.Error()))
				}
			}
		}(deployment)

		sourceCodeService := gitrepo.NewGitSourceCodeService()
		deployer, err := docker.NewLocalDeployer(s.db, logger, sourceCodeService)
		if err != nil {
			return errors.Join(err, errors.New("failed to create deployer"))
		}

		image := &backend.Image{
			ID:          uuid.New(),
			Name:        deployment.Name,
			Version:     deployment.Version,
			RegistryURL: "localhost:5000",
		}

		deployment.ImageID = image.ID

		deployment, err = deployer.Deploy(ctx, deployment, image, deployment.DeployMetadata)
		if err != nil {
			return errors.Join(err, errors.New("failed to deploy"))
		}

		deployment.Status = backend.DeploymentStatusReady

		deployment, err = s.deploymentService.Update(ctx, deployment)
		if err != nil {
			return errors.Join(err, errors.New("failed to update deployment"))
		}

		return nil
	}
}
