package backend

import (
	"context"

	"github.com/google/uuid"
)

type Network struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ProjectID  uuid.UUID `json:"projectId"`
	PlatformID uint64    `json:"platformId"`
}

func NewNetwork(name string, projectID uuid.UUID, platformID uint64) *Network {
	return &Network{
		ID:         uuid.New(),
		Name:       name,
		ProjectID:  projectID,
		PlatformID: platformID,
	}
}

type NetworkService interface {
	Create(ctx context.Context, network *Network) (*Network, error)
}
