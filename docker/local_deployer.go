package docker

import (
	"context"
	"fmt"
	backend "glider"

	"github.com/docker/docker/client"
)

type LocalDeployer struct {
	cli              *client.Client
	imageService     backend.ImageService
	containerService backend.ContainerService
}

func NewLocalDeployer() (*LocalDeployer, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &LocalDeployer{
		cli: cli,
	}, nil
}

func (d *LocalDeployer) Deploy(ctx context.Context, deployment *backend.Deployment, image *backend.Image) (*backend.Deployment, error) {
	image, err := d.imageService.BuildImage(ctx, image)
	if err != nil {
		return nil, err
	}

	if err := d.imageService.PullImage(ctx, image); err != nil {
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
