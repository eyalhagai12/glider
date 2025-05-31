package pg

import (
	"context"
	"database/sql"
	backend "glider"
)

type DeploymentService struct {
	db *sql.DB
}

func NewDeploymentService(db *sql.DB) *DeploymentService {
	return &DeploymentService{
		db: db,
	}
}

func (s *DeploymentService) Create(ctx context.Context, dep *backend.Deployment) (*backend.Deployment, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO deployments (id, name, version, environment, project_id, status, replicas, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())",
		dep.ID, dep.Name, dep.Version, dep.Environment, dep.ProjectID, dep.Status, dep.Replicas)
	if err != nil {
		return nil, err
	}

	for _, tag := range dep.Tags {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO tags (deployment_id, name, is_system) VALUES ($1, $2, $3)",
			dep.ID, tag.Name, tag.IsSystem)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return dep, nil
}
