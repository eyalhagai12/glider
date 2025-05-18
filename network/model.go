package network

import "github.com/google/uuid"

type Network struct {
	ID            uuid.UUID `db:"id"`
	InterfaceName string    `db:"interface_name"`
	IpAddress     string    `db:"ip_address"`
	Port          string    `db:"port"`
	ProjectID     uuid.UUID `db:"project_id"`
}

func NewNetwork(interfaceName string, ipAddress string, port string, projectID uuid.UUID) *Network {
	return &Network{
		ID:            uuid.New(),
		InterfaceName: interfaceName,
		IpAddress:     ipAddress,
		Port:          port,
		ProjectID:     projectID,
	}
}
