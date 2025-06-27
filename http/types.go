package http

import (
	backend "glider"

	"github.com/google/uuid"
)

type createUser struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	HashedPassword string `json:"hashedPassword" binding:"required"`
}

type loginUser struct {
	Username       string `json:"username" binding:"required"`
	HashedPassword string `json:"hashedPassword" binding:"required"`
}

type deployRequest struct {
	Name           string                     `json:"name" binding:"required"`
	Description    string                     `json:"description"`
	Version        string                     `json:"version" binding:"required"`
	Replicas       int                        `json:"replicas" binding:"required"`
	ProjectID      uuid.UUID                  `json:"projectId" binding:"required"`
	Environment    string                     `json:"environment" binding:"required"`
	Tags           []string                   `json:"tags"`
	DeployMetadata backend.DeploymentMetadata `json:"deployMetadata" binding:"required"`
}

type fetchDeployment struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type DeploymentAction struct {
	Name   string `json:"name" binding:"required"`
	Method string `json:"method" binding:"required"`
	Path   string `json:"path" binding:"required"`
	Body   string `json:"body,omitempty"`
}

type DeploymentResponse struct {
	Deployment *backend.Deployment `json:"deployment"`
	Actions    []DeploymentAction  `json:"actions"`
	URL        string              `json:"url"`
}

type createProject struct {
	Name string `json:"name" binding:"required"`
}
