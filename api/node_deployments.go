package api

import (
	"glider/containers"
	"glider/workerpool"

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
}

func NewNodeDeploymentHandlers(workerPool *workerpool.WorkerPool, dockerCli *client.Client) NodeDeploymentHandlers {
	return NodeDeploymentHandlers{
		workerPool: workerPool,
		dockerCli:  dockerCli,
	}
}

func (h NodeDeploymentHandlers) Deploy(c *gin.Context, request NodeDeployRequest) (*containers.Container, error) {
	return nil, nil
}
