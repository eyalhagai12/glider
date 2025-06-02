package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	backend "glider"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"

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
	path = strings.TrimLeft(path, "/")
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

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, requestUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header = c.Request.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData map[string]any
	if err := json.Unmarshal(data, &responseData); err != nil {
		return nil, err
	}

	logger.Debug("response data", slog.Any("response_data", responseData))

	return responseData, nil
}
