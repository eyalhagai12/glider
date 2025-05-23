package api

import (
	"fmt"
	"glider/network"
	"glider/nodes"
	"glider/projects"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ProjectHandlers struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewProjectHandlers(db *sqlx.DB, logger *slog.Logger) ProjectHandlers {
	return ProjectHandlers{
		db:     db,
		logger: logger,
	}
}

func (p ProjectHandlers) CreateProject(c *gin.Context, request NewProjectRequest) (*projects.Project, error) {
	p.logger.Info("Creating new project", "name", request.Name, "description", request.Description)

	tx, err := p.db.BeginTxx(c, nil)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to begin transaction: %v", err))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				p.logger.Error("failed to rollback transaction", "error", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				p.logger.Error("failed to commit transaction", "error", commitErr)
			}
		}
	}()

	project := projects.NewProject(request.Name, request.Description)
	err = projects.StoreProject(tx, project)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to store project: %v", err))
	}

	netIP := "10.0.0.1/24" // need to be unique for each project
	netPort := "51820"     // need to be unique for each project
	interfaceName := fmt.Sprintf("wg-%s", project.ID.String()[:8])

	net := network.NewNetwork(interfaceName, netIP, netPort, project.ID)
	err = network.InitializeVPN(p.logger, net)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to initialize VPN: %v", err))
	}

	err = network.StoreNetwork(tx, net)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to store network: %v", err))
	}

	err = broadcastNetworkToAgentNodes(p.db, net)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to broadcast network to agent nodes: %v", err))
	}

	p.logger.Info("Network initialized successfully", "networkID", net.ID)
	p.logger.Info("Successfully broadcasted network to agent nodes", "networkID", net.ID)
	p.logger.Info("Project created successfully", "projectID", project.ID)

	return project, nil
}

func broadcastNetworkToAgentNodes(db *sqlx.DB, network *network.Network) error {
	nodes, err := nodes.GetAvailableNodes(db)
	if err != nil {
		return fmt.Errorf("failed to get available nodes: %v", err)
	}

	for _, node := range nodes {
		err := node.AddNetwork(network, network.PublicKey, []string{})
		if err != nil {
			return fmt.Errorf("failed to send network info to node %s: %v", node.ID, err)
		}
	}

	return nil
}
