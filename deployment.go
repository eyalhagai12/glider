package backend

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID           uuid.UUID `json:"id"`
	DeploymentID uuid.UUID `json:"deployment_id"`
	Name         string    `json:"name"`
	IsSystem     bool      `json:"isSystem"`
}

type Deployment struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	ProjectID   uuid.UUID `json:"project_id"`
	Status      string    `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Tags           []Tag          `json:"tags,omitempty"`
	DeployMetadata map[string]any `json:"deploy_metadata,omitempty"`
}

type DeploymentService interface {
	Create(ctx context.Context, deployment *Deployment) (*Deployment, error)
}
