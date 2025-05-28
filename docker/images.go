package docker

import (
	"context"
	"database/sql"
	backend "glider"

	"github.com/docker/docker/client"
)

type DockerImageService struct {
	cli *client.Client
	db  *sql.DB
}

func NewDockerImageService(cli *client.Client, db *sql.DB) *DockerImageService {
	return &DockerImageService{
		cli: cli,
		db:  db,
	}
}

func (s *DockerImageService) BuildImage(ctx context.Context, image *backend.Image) (*backend.Image, error) {
	
	return nil, nil
}
