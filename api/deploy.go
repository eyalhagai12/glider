package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"glider/containers"
	"glider/deployments"
	"glider/gitrepo"
	"glider/images"
	"glider/nodes"
	"glider/workerpool"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DeployHandlers struct {
	db         *sqlx.DB
	workerPool *workerpool.WorkerPool
	dockerCli  *client.Client
	logger     *slog.Logger
}

func NewDeployHandlers(db *sqlx.DB, workerPool *workerpool.WorkerPool, logger *slog.Logger, dockerCli *client.Client) DeployHandlers {
	return DeployHandlers{
		db:         db,
		workerPool: workerPool,
		dockerCli:  dockerCli,
		logger:     logger,
	}
}

func (d DeployHandlers) Deploy(c *gin.Context, request DeployRequest) (*deployments.Deployment, error) {
	deployment := deployments.NewDeployment(request.DeploymentName, request.GithubRepo, request.GithubBranch, request.Replicas)
	err := deployments.StoreDeployment(d.db, deployment)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to store deployment: %v", err))
	}

	d.workerPool.Submit(d.deploy(deployment, request))

	return deployment, nil
}

func (d DeployHandlers) uploadImage(ctx context.Context, cli *client.Client, path string, repo string, branch string, deploymentName string, dockerFilePath string, tag string, namespace string) error {
	defer os.RemoveAll(path) // might need to make sure that the error is used and handled

	_, err := gitrepo.CloneRepository(path, repo, branch)
	if err != nil {
		return err
	}

	image, err := images.BuildImage(ctx, cli, deploymentName, namespace, path, dockerFilePath, tag)
	if err != nil {
		return err
	}
	defer image.Body.Close()

	imageName := images.FormatImageName(deploymentName, namespace, tag)
	if err := images.StoreImage(ctx, cli, imageName, images.RegistryAuth{}); err != nil {
		return err
	}

	return nil
}

func (d DeployHandlers) sendRequestToAgent(deployment *deployments.Deployment, node *nodes.Node, namespace string, tag string) ([]containers.Container, error) {
	agentClient := &http.Client{}
	req, err := http.NewRequest("POST", node.DeploymentURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	nodeDeployRequest := NodeDeployRequest{
		Image:          images.FormatImageName(deployment.Name, namespace, tag),
		DeploymentName: deployment.Name,
		DeploymentUUID: deployment.ID,
		Replicas:       deployment.TargetReplicaCount,
		NodeUUID:       node.ID,
	}

	body, err := json.Marshal(nodeDeployRequest)
	if err != nil {
		d.logger.Error("Failed to marshal request body\n", "error", err)
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	resp, err := agentClient.Do(req)
	if err != nil {
		d.logger.Error("Failed to send request to agent\n", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		d.logger.Error("Failed to deploy on node\n", "status", resp.Status)
		return nil, fmt.Errorf("failed to deploy on node: %s", resp.Status)
	}

	var containerList []containers.Container
	if err := json.NewDecoder(resp.Body).Decode(&containerList); err != nil {
		d.logger.Error("Failed to decode response\n", "error", err)
		return nil, err
	}
	if len(containerList) == 0 {
		d.logger.Error("No containers deployed\n")
		return nil, fmt.Errorf("no containers deployed")
	}

	return containerList, nil
}

func (d DeployHandlers) deploy(deployment *deployments.Deployment, request DeployRequest) workerpool.Task {
	return func(ctx context.Context, id int) error {
		defer func(deployment *deployments.Deployment) {
			d.logger.Info("Updating deployment status\n", "deployment ID", deployment.ID)
			err := deployments.UpdateDeployment(d.db, deployment)
			if err != nil {
				d.logger.Error("Failed to update deployment\n", "error", err)
			} else {
				d.logger.Info("Deployment updated successfully\n", "deployment name", deployment.Name, "status", deployment.Status)
			}
		}(deployment)

		logger := slog.Default()

		path := "./tmp/clones/" + request.DeploymentName
		err := d.uploadImage(ctx, d.dockerCli, path, request.GithubRepo, request.GithubBranch, request.DeploymentName, request.DockerfilePath, request.Tag, request.Namespace)
		if err != nil {
			logger.Error("Failed to upload image\n", "error", err)
			return err
		}

		deployment.Status = deployments.DeploymentStatusImageUploaded
		deployment.UpdatedAt = time.Now()

		logger.Info("Image uploaded successfully\n", "image", images.FormatImageName(deployment.Name, request.Namespace, request.Tag))
		logger.Info("Fetching available nodes\n")

		nodes, err := nodes.GetAvailableNodes(d.db)
		if err != nil {
			logger.Error("Failed to fetch available nodes\n", "error", err)
			return err
		}

		logger.Info("Found available nodes\n", "count", len(nodes))

		if len(nodes) == 0 {
			logger.Error("No available nodes\n")
			return fmt.Errorf("no available nodes")
		}

		logger.Info("Selecting a node\n")

		randomIndex := rand.Intn(len(nodes))
		node := nodes[randomIndex]

		deployment.Status = deployments.DeploymentStatusDeploying
		deployment.ReplicaCount = 0
		deployment.TargetReplicaCount = request.Replicas
		deployment.GithubRepo = request.GithubRepo
		deployment.GithubBranch = request.GithubBranch
		deployment.Name = request.DeploymentName
		deployment.DeletedAt = pq.NullTime{}
		deployment.CreatedAt = time.Now()
		deployment.UpdatedAt = time.Now()

		logger.Info("Sending deploy request to node\n", "node deploy address", node.DeploymentURL, "node uuid", node.ID)
		logger.Info("Deploying", "deployment name", deployment.Name, "image", images.FormatImageName(deployment.Name, request.Namespace, request.Tag), "replicas", request.Replicas, "deployment uuid", deployment.ID)

		containers, err := d.sendRequestToAgent(deployment, &node, request.Namespace, request.Tag)
		if err != nil {
			logger.Error("Failed to send request to agent\n", "error", err)
			return err
		}

		deployment.Status = deployments.DeploymentStatusReady
		deployment.ReplicaCount = len(containers)
		deployment.UpdatedAt = time.Now()

		logger.Info("Deployment is ready\n", "deployment name", deployment.Name, "active replicas", len(containers))
		logger.Info("Containers deployed\n", "count", containers)
		logger.Info("Deployment completed successfully\n", "deployment name", deployment.Name)

		return nil
	}
}
