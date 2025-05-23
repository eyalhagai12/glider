package api

import (
	"glider/nodes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type NodeHandlers struct {
	db *sqlx.DB
}

type RegisterNodeRequest struct {
	DeploymentURL string `json:"deployment_url"`
	MetricsURL    string `json:"metrics_url"`
	HealthURL     string `json:"health_url"`
	ConnectionURL string `json:"connection_url"`
}

func NewNodeHandlers(db *sqlx.DB) NodeHandlers {
	return NodeHandlers{
		db: db,
	}
}

func (n NodeHandlers) RegisterNewNode(c *gin.Context, request RegisterNodeRequest) (*nodes.Node, error) {
	id := uuid.New()

	node := nodes.Node{
		ID:            id,
		DeploymentURL: request.DeploymentURL,
		MetricsURL:    request.MetricsURL,
		HealthURL:     request.HealthURL,
		ConnectionURL: request.ConnectionURL,
		Status:        nodes.StatusActive,
	}

	_, err := n.db.Exec(`
		INSERT INTO nodes (id, deployment_url, metrics_url, health_url, connection_url, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, node.ID, node.DeploymentURL, node.MetricsURL, node.HealthURL, node.ConnectionURL, node.Status)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}

	return &node, nil
}
