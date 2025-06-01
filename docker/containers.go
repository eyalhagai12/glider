package docker

import (
	"context"
	"database/sql"
	backend "glider"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerContainerService struct {
	cli *client.Client
	db  *sql.DB
}

func NewDockerContainerService(cli *client.Client, db *sql.DB) *DockerContainerService {
	return &DockerContainerService{
		cli: cli,
		db:  db,
	}
}

func (s *DockerContainerService) Create(ctx context.Context, cont *backend.Container, image *backend.Image) (*backend.Container, error) {
	createResp, err := s.cli.ContainerCreate(ctx, &container.Config{
		Image: image.ImagePath(),
	}, nil, nil, nil, cont.Name)
	if err != nil {
		return nil, err
	}
	cont.PlatformID = createResp.ID

	err = s.cli.ContainerStart(ctx, cont.PlatformID, container.StartOptions{})
	if err != nil {
		return nil, err
	}

	resp, err := s.cli.ContainerInspect(ctx, cont.PlatformID)
	if err != nil {
		return nil, err
	}

	// cont.Port = resp.NetworkSettings.Ports["8080/tcp"][0].HostPort
	cont.Port = "8080"
	cont.Host = resp.NetworkSettings.IPAddress

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO containers (id, platform_id, name, deployment_id, image_id, host, port)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, cont.ID, cont.PlatformID, cont.Name, cont.DeploymentID, cont.ImageID, cont.Host, cont.Port)
	if err != nil {
		return nil, err
	}

	return cont, nil
}
