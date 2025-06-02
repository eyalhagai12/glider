package backend

import (
	"context"

	"github.com/google/uuid"
)

type Project struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProjectService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	Create(ctx context.Context, project *Project) (*Project, error)
}
