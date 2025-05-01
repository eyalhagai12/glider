package images

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

const registryURL = "localhost:5000"

func StoreImage(ctx context.Context, cli *client.Client, name string, namespace string, tag string) error {
	imageName := registryURL + "/" + namespace + "/" + name + ":" + tag

	auth := struct{}{}
	authData, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	encodedAuth := base64.URLEncoding.EncodeToString(authData)

	pushResponse, err := cli.ImagePush(ctx, imageName, image.PushOptions{RegistryAuth: encodedAuth})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, pushResponse)
	defer pushResponse.Close()

	return nil
}
