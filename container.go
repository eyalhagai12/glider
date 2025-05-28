package backend

import (
	"context"

	"github.com/google/uuid"
)

type Container struct {
	ID           uuid.UUID
	PlatformID   string
	Name         string
	DeploymentID uuid.UUID
	ImageID      uuid.UUID

	Host string
	Port string
}

type ContainerService interface {
	Create(ctx context.Context, container *Container, image *Image) (*Container, error)
}
