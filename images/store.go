package images

import (
	"context"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

const registryURL = "http://localhost:5000"

func StoreImage(ctx context.Context, name string, namespace string, tag string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	imageName := registryURL + "/" + namespace + "/" + name + ":" + tag
	pushResponse, err := cli.ImagePush(ctx, imageName, image.PushOptions{})
	if err != nil {
		return err
	}
	defer pushResponse.Close()

	return nil
}
