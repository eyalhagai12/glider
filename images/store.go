package images

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

const registryURL = "localhost:5000"

func StoreImage(ctx context.Context, cli *client.Client, name string, namespace string, tag string) error {
	authConfig := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "glider",
		Password: "glider123",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imageName := registryURL + "/" + namespace + "/" + name + ":" + tag
	pushResponse, err := cli.ImagePush(ctx, imageName, image.PushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return err
	}
	defer pushResponse.Close()

	return nil
}
