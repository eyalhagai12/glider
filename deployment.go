package backend

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	DeploymentStatusPending   = "pending"
	DeploymentStatusDeploying = "deploying"
	DeploymentStatusFailed    = "failed"
	DeploymentStatusReady     = "ready"
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
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	ProjectID   uuid.UUID `json:"projectId"`
	Status      string    `json:"status"`
	ImageID     uuid.UUID `json:"imageId"`
	Replicas    int       `json:"replicas"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`

	Tags           []Tag          `json:"tags,omitempty"`
	DeployMetadata map[string]any `json:"deployMetadata,omitempty"`
}

type DeploymentService interface {
	Create(ctx context.Context, deployment *Deployment) (*Deployment, error)
}

func TagsFromList(tags []string) []Tag {
	result := make([]Tag, len(tags))
	for i, tag := range tags {
		result[i] = Tag{
			Name:     tag,
			IsSystem: false,
		}
	}

	return result
}
