package docker

import (
	"context"
	"database/sql"
	"fmt"
	backend "glider"
	"glider/gitrepo"

	"github.com/docker/docker/client"
)

type LocalDeployer struct {
	cli               *client.Client
	imageService      backend.ImageService
	containerService  backend.ContainerService
	sourceCodeService backend.SourceCodeService
}

func NewLocalDeployer(db *sql.DB) (*LocalDeployer, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	imageService := NewDockerImageService(cli, db)
	containerService := NewDockerContainerService(cli, db)
	sourceCodeService := gitrepo.NewGitsSourceCodeService()

	return &LocalDeployer{
		cli:               cli,
		imageService:      imageService,
		containerService:  containerService,
		sourceCodeService: sourceCodeService,
	}, nil
}

func (d *LocalDeployer) Deploy(ctx context.Context, deployment *backend.Deployment, image *backend.Image, deploymentMetadata backend.DeploymentMetadata) (*backend.Deployment, error) {
	path, err := d.sourceCodeService.Fetch(ctx, deploymentMetadata)
	if err != nil {
		return nil, err
	}

	image, err = d.imageService.BuildImage(ctx, image, path)
	if err != nil {
		return nil, err
	}

	if _, err := d.imageService.PullImage(ctx, image); err != nil {
		return nil, err
	}

	for replica := range deployment.Replicas {
		_, err := d.containerService.Create(ctx, &backend.Container{
			DeploymentID: deployment.ID,
			ImageID:      image.ID,
			Name:         fmt.Sprintf("%s-%d", deployment.Name, replica),
		}, image)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
