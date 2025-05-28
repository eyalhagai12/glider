package pg

import (
	"context"
	"database/sql"
	backend "glider"

	"github.com/google/uuid"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) Create(ctx context.Context, user *backend.User) (*backend.User, error) {
	currentUser := backend.GetUserFromContext(ctx)
	if currentUser != nil {
		return nil, backend.ErrUnauthorized
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (id, username, email, hashed_password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, username, email, role, created_at, updated_at
	`
	row := s.db.QueryRowContext(ctx, query, user.ID, user.Username, user.Email, user.HashedPassword, user.Role)
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*backend.User, error) {
	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := s.db.QueryRowContext(ctx, query, id)

	user := &backend.User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, backend.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*backend.User, error) {
	query := `
		SELECT id, username, email, role, hashed_password, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	row := s.db.QueryRowContext(ctx, query, username)

	user := &backend.User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, backend.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}
