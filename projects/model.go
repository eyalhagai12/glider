package projects

import "github.com/google/uuid"

type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	DeletedAt   string    `json:"deleted_at"`
}

func NewProject(name string, description string) *Project {
	return &Project{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
	}
}
