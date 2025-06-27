package pg

import (
	"context"
	"database/sql"
	backend "glider"

	"github.com/google/uuid"
)

type ProjectService struct {
	db *sql.DB
}

func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{
		db: db,
	}
}

func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*backend.Project, error) {
	row, err := s.db.QueryContext(ctx, "SELECT * FROM projects WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if !row.Next() {
		return nil, backend.ErrNotFound
	}

	project := &backend.Project{}
	if err := row.Scan(&project.ID, &project.Name); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Create(ctx context.Context, project *backend.Project) (*backend.Project, error) {
	query := `
		INSERT INTO projects (id, name)
		VALUES ($1, $2)
		RETURNING id, name
	`
	row := s.db.QueryRowContext(ctx, query, project.ID, project.Name)

	if err := row.Scan(&project.ID, &project.Name); err != nil {
		return nil, err
	}

	return project, nil
}
