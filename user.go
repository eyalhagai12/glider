package backend

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
	RoleGuest = "guest"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	HashedPassword string    `json:"hashedPassword"`

	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.Join(ErrInvalidInput, errors.New("username is required"))
	}
	if u.Email == "" {
		return errors.Join(ErrInvalidInput, errors.New("email is required"))
	}
	if u.HashedPassword == "" {
		return errors.Join(ErrInvalidInput, errors.New("hashed password is required"))
	}

	return nil
}

type UserService interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
