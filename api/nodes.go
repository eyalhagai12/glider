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
	}

	_, err := n.db.Exec(`
		INSERT INTO nodes (id, deployment_url, metrics_url)
		VALUES ($1, $2, $3)
	`, node.ID, node.DeploymentURL, node.MetricsURL)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}

	return &node, nil
}
