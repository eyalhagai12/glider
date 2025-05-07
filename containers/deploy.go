package containers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func DeployConainers(c context.Context, cli *client.Client, deploymentName string, deploymentUUID string, image string, replicas int, nodeID string) ([]Container, error) {
	containerList := make([]Container, replicas)
	for i := range replicas {
		containerName := fmt.Sprintf("%s-%d", deploymentName, i+1)

		cli.ContainerRemove(c, containerName, container.RemoveOptions{})

		resp, err := cli.ContainerCreate(c, &container.Config{
			Image: image,
		}, nil, nil, nil, containerName)
		if err != nil {
			return nil, err
		}

		err = cli.ContainerStart(c, resp.ID, container.StartOptions{})
		if err != nil {
			return nil, err
		}

		containerInspection, err := cli.ContainerInspect(c, resp.ID)
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

		containerList[i] = Container{
			ID:             resp.ID,
			Name:           deploymentName,
			DeploymentUUID: deploymentUUID,
			NodeID:         nodeID,
			IP:             ipAddress,
			Port:           hostPort,
		}
	}

	return containerList, nil
}
