package docker

import (
	"context"
	"database/sql"
	"fmt"
	backend "glider"
	"log/slog"
	"os"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

type LocalDeployer struct {
	cli    *client.Client
	logger *slog.Logger

	imageService      backend.ImageService
	containerService  backend.ContainerService
	sourceCodeService backend.SourceCodeService
}

func NewLocalDeployer(db *sql.DB, logger *slog.Logger, sourceCodeService backend.SourceCodeService) (*LocalDeployer, error) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		return nil, err
	}

	imageService := NewDockerImageService(cli, db, logger)
	containerService := NewDockerContainerService(cli, db)

	return &LocalDeployer{
		cli:               cli,
		logger:            logger,
		imageService:      imageService,
		containerService:  containerService,
		sourceCodeService: sourceCodeService,
	}, nil
}

func (d *LocalDeployer) Deploy(ctx context.Context, deployment *backend.Deployment, image *backend.Image, deploymentMetadata backend.DeploymentMetadata) (*backend.Deployment, error) {
	logger := d.logger.With("deployment_id", deployment.ID)

	logger.Debug("fetching source code")
	path, err := d.sourceCodeService.Fetch(ctx, deploymentMetadata)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(path)

	logger.Debug("building image")
	image, err = d.imageService.BuildImage(ctx, image, path)
	if err != nil {
		return nil, err
	}

	logger.Debug("pulling image")
	if _, err := d.imageService.PullImage(ctx, image); err != nil {
		return nil, err
	}

	logger.Debug("creating containers")
	for replica := range deployment.Replicas {
		logger.Debug("creating container", slog.Int("replica", replica))
		_, err := d.containerService.Create(ctx, &backend.Container{
			ID:           uuid.New(),
			DeploymentID: deployment.ID,
			ImageID:      image.ID,
			Name:         fmt.Sprintf("%s-%d", deployment.Name, replica),
		}, image)
		if err != nil {
			return nil, err
		}
	}

	logger.Debug("deployed successfully")
	return deployment, nil
}
