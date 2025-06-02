package backend

import (
	"context"

	"github.com/google/uuid"
)

type Network struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ProjectID uuid.UUID `json:"projectId"`
	Mask      string    `json:"mask"`
}

func NewNetwork(name string, projectID uuid.UUID, mask string) *Network {
	return &Network{
		ID:        uuid.New(),
		Name:      name,
		ProjectID: projectID,
		Mask:      mask,
	}
}

type NetworkService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Network, error)
	Create(ctx context.Context, network *Network) (*Network, error)
}
