package nodes

import "github.com/google/uuid"

type Node struct {
	ID            uuid.UUID `db:"id" json:"id"`
	DeploymentURL string    `db:"deployment_url" json:"deployment_url"`
	MetricsURL    string    `db:"metrics_url" json:"metrics_url"`
	HealthURL     string    `db:"health_url" json:"health_url"`
	ConnectionURL string    `db:"connection_url" json:"connection_url"`
	Status        string    `db:"status" json:"status"`
}
