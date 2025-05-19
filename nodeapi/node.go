package nodeapi

import (
	"glider/containers"
	"glider/images"
	"glider/network"
	"glider/resources"
	"glider/workerpool"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type NodeHandlers struct {
	workerPool *workerpool.WorkerPool
	dockerCli  *client.Client
	logger     *slog.Logger
}

func NodeAgentHandlers(workerPool *workerpool.WorkerPool, dockerCli *client.Client, logger *slog.Logger) NodeHandlers {
	return NodeHandlers{
		workerPool: workerPool,
		dockerCli:  dockerCli,
		logger:     logger,
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

func (h NodeHandlers) ConnectToVPN(c *gin.Context, request NodeConnectRequest) error {
	h.logger.Info("Connecting to VPN", "interface", request.Interface, "ip", request.IP, "endpoint", request.Endpoint)

	if err := network.ConnectToVPN(h.logger, request.Interface, request.IP, request.PublicKey, request.Endpoint, request.AllowedIPs); err != nil {
		h.logger.Error("Failed to connect to VPN", "error", err)
		return err
	}

	h.logger.Info("Connected to VPN successfully", "interface", request.Interface)
	
	return nil
}
