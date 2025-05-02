package api

import (
	"encoding/base64"
	"encoding/json"
	"glider/containers"
	"glider/workerpool"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type NodeDeploymentHandlers struct {
	workerPool *workerpool.WorkerPool
	dockerCli  *client.Client
}

type NodeDeployRequest struct {
	DeploymentName string
	Replicas       int
	Image          string
	Registry       string
}

func NewNodeDeploymentHandlers(workerPool *workerpool.WorkerPool, dockerCli *client.Client) NodeDeploymentHandlers {
	return NodeDeploymentHandlers{
		workerPool: workerPool,
		dockerCli:  dockerCli,
	}
}

func (h NodeDeploymentHandlers) Deploy(c *gin.Context, request NodeDeployRequest) ([]containers.Container, error) {
	// pull docker image from registry
	auth := struct{}{}
	authData, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}
	encodedAuth := base64.URLEncoding.EncodeToString(authData)

	imageName := request.Registry + "/" + request.Image
	imagePullOptions := image.PullOptions{
		RegistryAuth: encodedAuth,
	}

	out, err := h.dockerCli.ImagePull(c, imageName, imagePullOptions)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// create containers
	for i := 0; i < request.Replicas; i++ {
		resp, err := h.dockerCli.ContainerCreate(c, &container.Config{
			Image: request.Image,
		}, nil, nil, nil, request.DeploymentName+"-"+string(i+1))
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
