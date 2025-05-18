package api

import (
	"fmt"
	"glider/network"
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

	netIP := "10.0.0.1/24"
	netPort := "51820"

	net := network.NewNetwork(project.Name, netIP, netPort, project.ID)
	err = network.InitializeVPN(p.logger, net)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to initialize VPN: %v", err))
	}
	
	err = network.StoreNetwork(tx, net)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to store network: %v", err))
	}

	p.logger.Info("Network initialized successfully", "networkID", net.ID)
	p.logger.Info("Project created successfully", "projectID", project.ID)

	return project, nil
}
