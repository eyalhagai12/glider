package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	backend "glider"

	"github.com/google/uuid"
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
		"INSERT INTO deployments (id, name, version, environment, project_id, status, replicas, deploy_metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())",
		dep.ID, dep.Name, dep.Version, dep.Environment, dep.ProjectID, dep.Status, dep.Replicas, dep.DeployMetadata)
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

func (s *DeploymentService) Update(ctx context.Context, dep *backend.Deployment) (*backend.Deployment, error) {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE deployments SET name = $1, version = $2, environment = $3, project_id = $4, status = $5, replicas = $6, image_id = $7, updated_at = NOW() WHERE id = $8",
		dep.Name, dep.Version, dep.Environment, dep.ProjectID, dep.Status, dep.Replicas, dep.ImageID, dep.ID,
	)
	if err != nil {
		return nil, err
	}

	for _, tag := range dep.Tags {
		_, err = s.db.ExecContext(
			ctx,
			"INSERT INTO tags (deployment_id, name, is_system) VALUES ($1, $2, $3)",
			dep.ID, tag.Name, tag.IsSystem)
		if err != nil {
			return nil, err
		}
	}
	return dep, nil
}

func (s *DeploymentService) GetByID(ctx context.Context, id uuid.UUID) (*backend.Deployment, error) {
	row, err := s.db.QueryContext(ctx, "SELECT * FROM deployments WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	row.Next()

	deployMetadata := json.RawMessage{}

	deployment := &backend.Deployment{}
	if err := row.Scan(
		&deployment.ID, &deployment.Name, &deployment.Version,
		&deployment.Environment, &deployment.ProjectID,
		&deployment.Status, &deployment.ImageID,
		&deployment.Replicas, &deployment.CreatedAt,
		&deployment.UpdatedAt, &deployment.DeletedAt,
		&deployMetadata); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(deployMetadata, &deployment.DeployMetadata); err != nil {
		return nil, err
	}

	return deployment, nil
}
