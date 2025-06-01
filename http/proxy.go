package http

import (
	"errors"
	"fmt"
	backend "glider"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/eyalhagai12/hagio/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterProxyHandler() {
	routeGroup := s.apiRoutes.Group("/proxy")
	routeGroup.Any("/:id/*path", handler.FromFunc(s.handleRequest, http.StatusOK))
}

func (s *Server) handleRequest(c *gin.Context, _ any) (any, error) {
	deploymentUuid := uuid.MustParse(c.Param("id"))
	path := c.Param("path")

	if path == "" {
		path = "/"
	}

	logger := s.logger.With("deployment_uuid", deploymentUuid, "path", path)
	logger.Debug("proxying request")

	deployment, err := s.deploymentService.GetByID(c.Request.Context(), deploymentUuid)
	if err != nil {
		return nil, err
	}

	if deployment.DeletedAt != nil {
		return nil, errors.Join(backend.ErrNotFound, errors.New("deployment not found"))
	}

	if deployment.Status != backend.DeploymentStatusReady {
		return nil, errors.Join(backend.ErrInvalidState, errors.New("deployment not ready"))
	}

	containers, err := s.containerService.GetContainersByDeploymentID(c.Request.Context(), deploymentUuid)
	if len(containers) == 0 {
		return nil, errors.Join(backend.ErrNotFound, errors.New("no containers found"))
	}

	container := containers[rand.Intn(len(containers))]
	requestUrl := fmt.Sprintf("http://%s:%s/%s", container.Host, container.Port, path)

	logger.Debug("proxying request to", slog.String("request_url", requestUrl))

	url, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(&http.Request{
		Method: c.Request.Method,
		URL:    url,
		Header: c.Request.Header,
		Body:   c.Request.Body,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
