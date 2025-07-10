package backend

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/google/uuid"
)

type IPAddress uint32

func (i IPAddress) String() string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, uint32(i))
	return fmt.Sprintf("%s/24", ip.String())
}

type Network struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ProjectID  uuid.UUID `json:"projectId"`
	PlatformID uint64    `json:"platformId"`
	Address    IPAddress `json:"address"`
	ListenPort uint16    `json:"listenPort"`
	PrivateKey string    `json:"privateKey"`
}

func NewNetwork(name string, projectID uuid.UUID, platformID uint64, address IPAddress, listenPort uint16, privateKey string) *Network {
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
