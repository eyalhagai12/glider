package backend

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	DeploymentStatusPending   = "pending"
	DeploymentStatusDeploying = "deploying"
	DeploymentStatusFailed    = "failed"
	DeploymentStatusReady     = "ready"
)

const DeploymentURLFormat = "/deployments/%s"
const ProxyURLFormat = "/proxy/%s"

type DeploymentMetadata map[string]any

func (dm DeploymentMetadata) Validate() error {
	if _, ok := dm["type"]; !ok {
		return errors.Join(ErrInvalidInput, errors.New("type is required"))
	}

	return nil
}

type Tag struct {
	ID           int       `json:"id"`
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

	Tags           []Tag              `json:"tags,omitempty"`
	DeployMetadata DeploymentMetadata `json:"deployMetadata,omitempty"`
}

func (d *Deployment) URL() string {
	return fmt.Sprintf(DeploymentURLFormat, d.ID)
}

func (d *Deployment) ProxyURL() string {
	return fmt.Sprintf(ProxyURLFormat, d.ID)
}

type DeploymentService interface {
	Create(ctx context.Context, deployment *Deployment) (*Deployment, error)
	Update(ctx context.Context, deployment *Deployment) (*Deployment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Deployment, error)
}

type SourceCodeService interface {
	// this function should load the source code into a directory and returng the path to it
	Fetch(ctx context.Context, deploymentMetadata DeploymentMetadata) (string, error)
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
