package http

import "github.com/google/uuid"

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
	Name           string
	Description    string
	Version        string
	Replicas       int
	ProjectID      uuid.UUID
	Environment    string
	ImageID        uuid.UUID
	Tags           []string
	DeployMetadata map[string]any
}
