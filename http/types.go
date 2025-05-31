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
