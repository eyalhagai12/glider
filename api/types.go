package api

import "github.com/google/uuid"

type NodeDeployRequest struct {
	Image          string    `json:"image"`
	DeploymentName string    `json:"deploymentName"`
	DeploymentUUID uuid.UUID `json:"deploymentUUID"`
	NodeUUID       uuid.UUID `json:"nodeUUID"`
	Replicas       int       `json:"replicas"`
}

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
	Replicas       int    `json:"replicas"`
	GithubRepo     string `json:"githubRepo"`
	GithubBranch   string `json:"githubBranch"`
}
