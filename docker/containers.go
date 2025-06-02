package docker

import (
	"context"
	"database/sql"
	backend "glider"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type DockerContainerService struct {
	cli    *client.Client
	db     *sql.DB
	logger *slog.Logger
}

func NewDockerContainerService(cli *client.Client, db *sql.DB, logger *slog.Logger) *DockerContainerService {
	return &DockerContainerService{
		cli:    cli,
		db:     db,
		logger: logger,
	}
}

func (s *DockerContainerService) Create(ctx context.Context, cont *backend.Container, image *backend.Image) (*backend.Container, error) {
	createResp, err := s.cli.ContainerCreate(ctx, &container.Config{
		Image: image.ImagePath(),
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "",
				},
			},
		},
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"backend_devcontainer_default": {},
		},
	}, nil, cont.Name)
	if err != nil {
		return nil, err
	}
	cont.PlatformID = createResp.ID

	err = s.cli.ContainerStart(ctx, cont.PlatformID, container.StartOptions{})
	if err != nil {
		return nil, err
	}

	containerNetworkSettings, err := pollContainerNetworkSettings(ctx, s.cli, cont.PlatformID, 10)
	if containerNetworkSettings == nil {
		return nil, err
	}

	if len(containerNetworkSettings.Ports["8080/tcp"]) > 0 {
		cont.Host = cont.Name
		cont.Port = "8080"
		s.logger.Debug("container started", "host", cont.Host, "port", cont.Port)
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO containers (id, platform_id, name, deployment_id, image_id, host, port)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, cont.ID, cont.PlatformID, cont.Name, cont.DeploymentID, cont.ImageID, cont.Host, cont.Port)
	if err != nil {
		return nil, err
	}

	return cont, nil
}

func (s *DockerContainerService) GetContainersByDeploymentID(ctx context.Context, deploymentID uuid.UUID) ([]*backend.Container, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, platform_id, name, deployment_id, image_id, host, port
		FROM containers
		WHERE deployment_id = $1
	`, deploymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []*backend.Container
	for rows.Next() {
		var cont backend.Container
		if err := rows.Scan(&cont.ID, &cont.PlatformID, &cont.Name, &cont.DeploymentID, &cont.ImageID, &cont.Host, &cont.Port); err != nil {
			return nil, err
		}
		containers = append(containers, &cont)
	}

	return containers, nil
}

func pollContainerNetworkSettings(ctx context.Context, cli *client.Client, containerID string, retries int) (*container.NetworkSettings, error) {
	var settings *container.NetworkSettings
	for range retries {
		resp, err := cli.ContainerInspect(ctx, containerID)
		if err != nil {
			return nil, err
		}

		settings = resp.NetworkSettings
		ports := settings.Ports["8080/tcp"]
		if len(ports) > 0 {
			return settings, nil
		}

		time.Sleep(time.Second)
	}

	return nil, backend.ErrNotFound
}
