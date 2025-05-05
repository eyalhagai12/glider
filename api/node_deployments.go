package api

import (
	"encoding/base64"
	"encoding/json"
	"glider/containers"
	"glider/workerpool"

	"fmt"

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
	DeploymentUUID string `json:"deploymentUUID"`
	DeploymentName string `json:"deploymentName"`
	Replicas       int    `json:"replicas"`
	Image          string `json:"image"`
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

	imageName := request.Image
	imagePullOptions := image.PullOptions{
		RegistryAuth: encodedAuth,
	}

	out, err := h.dockerCli.ImagePull(c, imageName, imagePullOptions)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// create containers
	containerList := make([]containers.Container, request.Replicas)
	for i := range request.Replicas {
		containerName := fmt.Sprintf("%s-%d", request.DeploymentName, i+1)

		h.dockerCli.ContainerRemove(c, containerName, container.RemoveOptions{})

		resp, err := h.dockerCli.ContainerCreate(c, &container.Config{
			Image: request.Image,
		}, nil, nil, nil, containerName)
		if err != nil {
			return nil, err
		}

		err = h.dockerCli.ContainerStart(c, resp.ID, container.StartOptions{})
		if err != nil {
			return nil, err
		}

		containerInspection, err := h.dockerCli.ContainerInspect(c, resp.ID)
		if err != nil {
			return nil, err
		}

		var hostPort string
		for _, bindings := range containerInspection.NetworkSettings.Ports {
			for _, binding := range bindings {
				hostPort = binding.HostPort
			}
		}

		ipAddress := containerInspection.NetworkSettings.IPAddress

		containerList[i] = containers.Container{
			ID:             resp.ID,
			Name:           request.DeploymentName,
			DeploymentUUID: request.DeploymentUUID,
			NodeID:         "adasdfadfasdfadsfadfadsf",
			IP:             ipAddress,
			Port:           hostPort,
		}
	}
	return containerList, nil
}
