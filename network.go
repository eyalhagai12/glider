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
	Address    string    `json:"address"`
	ListenPort uint16    `json:"listenPort"`
	PrivateKey string    `json:"privateKey"`
}

func NewNetwork(name string, projectID uuid.UUID, platformID uint64, address string, listenPort uint16, privateKey string) *Network {
	return &Network{
		ID:         uuid.New(),
		Name:       name,
		ProjectID:  projectID,
		PlatformID: platformID,
		Address:    address,
		ListenPort: listenPort,
		PrivateKey: privateKey,
	}
}

func (n *Network) SetKey(key string) {
	n.PrivateKey = key
}

type NetworkService interface {
	Create(ctx context.Context, network *Network) (*Network, error)
	GenerateKeys() (string, string, error)
}
