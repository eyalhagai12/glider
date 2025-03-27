package api

import (
	"github.com/gin-gonic/gin"
)

type DeployRequest struct {
	DeploymentName string `json:"deploymentName"`
	Replicas       int    `json:"replicas"`
	GithubRepo     string `json:"githubRepo"`
	GithubBranch   string `json:"githubBranch"`
	GithubToken    string `json:"githubToken"`
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
	return DeployResponse{
		DeploymentName: request.DeploymentName,
		Status:         "deploying",
	}, nil
}
