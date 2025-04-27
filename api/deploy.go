package api

import (
	"fmt"
	"glider/gitrepo"
	"glider/images"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeployRequest struct {
	DeploymentName string `json:"deploymentName"`
	Replicas       int    `json:"replicas"`
	GithubRepo     string `json:"githubRepo"`
	GithubBranch   string `json:"githubBranch"`
	GithubToken    string `json:"githubToken"`
	Tag            string `json:"tag"`
	Namespace      string `json:"namespace"`
	DockerfilePath string `json:"dockerfilePath"`
}

type DeployResponse struct {
	DeploymentName string `json:"deploymentName"`
	Status         string `json:"status"`
}

type DeployHandlers struct{}

func NewDeployHandlers() DeployHandlers {
	return DeployHandlers{}
}

func (d DeployHandlers) Deploy(c *gin.Context, request DeployRequest) (DeployResponse, error) {
	path := "./tmp/clones/" + request.DeploymentName
	_, err := gitrepo.CloneRepository(path, request.GithubRepo, request.GithubBranch)
	if err != nil {
		return DeployResponse{}, c.AbortWithError(http.StatusInternalServerError, err)
	}

	image, err := images.BuildImage(request.DeploymentName, path, request.DockerfilePath, request.Tag)
	if err != nil {
		return DeployResponse{}, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to build image: %v", err))
	}
	defer image.Body.Close()

	if err := images.StoreImage(c, request.DeploymentName, request.Namespace, request.Tag); err != nil {
		return DeployResponse{}, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to store image: %v", err))
	}

	return DeployResponse{
		DeploymentName: request.DeploymentName,
		Status:         "deploying",
	}, nil
}
