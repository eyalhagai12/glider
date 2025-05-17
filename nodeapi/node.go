package nodeapi

import (
	"glider/containers"
	"glider/images"
	"glider/resources"
	"glider/workerpool"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type NodeHandlers struct {
	workerPool *workerpool.WorkerPool
	dockerCli  *client.Client
}

func NewNodeDeploymentHandlers(workerPool *workerpool.WorkerPool, dockerCli *client.Client) NodeHandlers {
	return NodeHandlers{
		workerPool: workerPool,
		dockerCli:  dockerCli,
	}
}

func (h NodeHandlers) Deploy(c *gin.Context, request NodeDeployRequest) ([]containers.Container, error) {
	logger := log.Default()
	logger.SetOutput(os.Stdout)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	logger.SetPrefix("Agent: ")

	logger.Println("Starting deployment...")
	logger.Printf("Pulling image %s...\n", request.Image)
	err := images.PullImage(c, h.dockerCli, request.Image, images.RegistryAuth{})
	if err != nil {
		return nil, err
	}

	logger.Printf("Deploying %s...\n", request.DeploymentName)

	containerList, err := containers.DeployConainers(c, h.dockerCli, request.DeploymentName, request.DeploymentUUID, request.Image, request.Replicas, request.NodeUUID)
	if err != nil {
		logger.Printf("Error deploying containers: %v\n", err)
		return nil, err
	}

	logger.Printf("Deployment %s completed successfully.\n", request.DeploymentName)
	logger.Printf("Containers deployed: %d\n", len(containerList))
	return containerList, nil
}

func (h NodeHandlers) ReportMetrics(c *gin.Context, _ any) (resources.ResourceResponse, error) {
	disks, err := resources.FetchDisks()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return resources.ResourceResponse{}, err
	}
	cpuData, err := resources.FetchCPUs()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return resources.ResourceResponse{}, err
	}
	memoryData, err := resources.FetchMemory()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return resources.ResourceResponse{}, err
	}

	response := resources.ResourceResponse{
		NodeResources: resources.NodeResources{
			CPUs:   cpuData,
			Memory: memoryData,
			Disks:  disks,
		},
		Containers: []resources.ContainerUtilization{},
	}

	return response, nil
}
