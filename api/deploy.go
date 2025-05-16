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
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DeployHandlers struct {
	db         *sqlx.DB
	workerPool *workerpool.WorkerPool
	dockerCli  *client.Client
}

func NewDeployHandlers(db *sqlx.DB, workerPool *workerpool.WorkerPool, dockerCli *client.Client) DeployHandlers {
	return DeployHandlers{
		db:         db,
		workerPool: workerPool,
		dockerCli:  dockerCli,
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
	logger := log.Default()
	logger.SetOutput(os.Stdout)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	logger.SetPrefix("glider: ")

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
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	resp, err := agentClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		logger.Printf("Failed to deploy on node: %s\n", resp.Status)
		return nil, fmt.Errorf("failed to deploy on node: %s", resp.Status)
	}

	var containerList []containers.Container
	if err := json.NewDecoder(resp.Body).Decode(&containerList); err != nil {
		logger.Printf("Failed to decode response: %v\n", err)
		return nil, err
	}
	if len(containerList) == 0 {
		logger.Printf("No containers deployed\n")
		return nil, fmt.Errorf("no containers deployed")
	}

	return containerList, nil
}

func (d DeployHandlers) deploy(deployment *deployments.Deployment, request DeployRequest) workerpool.Task {
	return func(ctx context.Context, id int) error {
		defer deployments.StoreDeployment(d.db, deployment)

		logger := log.Default()
		logger.SetOutput(os.Stdout)
		logger.SetFlags(log.LstdFlags | log.Lshortfile)
		logger.SetPrefix("glider: ")
		logger.Printf("Worker %d started for deployment %s\n", id, request.DeploymentName)

		path := "./tmp/clones/" + request.DeploymentName

		logger.Printf("Cloning repository %s at branch %s\n", request.GithubRepo, request.GithubBranch)
		logger.Printf("Building image for deployment %s\n", request.DeploymentName)
		logger.Printf("Using Dockerfile at %s\n", request.DockerfilePath)
		logger.Printf("Using tag %s\n", request.Tag)
		logger.Printf("Using namespace %s\n", request.Namespace)
		logger.Printf("Using replicas %d\n", request.Replicas)

		err := d.uploadImage(ctx, d.dockerCli, path, request.GithubRepo, request.GithubBranch, request.DeploymentName, request.DockerfilePath, request.Tag, request.Namespace)
		if err != nil {
			return err
		}

		deployment.Status = deployments.DeploymentStatusImageUploaded
		deployment.UpdatedAt = time.Now()

		logger.Printf("Image %s uploaded successfully\n", request.DeploymentName)
		logger.Printf("Fetching available nodes\n")

		nodes, err := nodes.GetAvailableNodes(d.db)
		if err != nil {
			return err
		}

		logger.Printf("Found %d available nodes\n", len(nodes))

		if len(nodes) == 0 {
			return fmt.Errorf("no available nodes")
		}

		logger.Printf("Selecting a node\n")

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
		deployment.ID = uuid.New()

		logger.Printf("Sending deploy request to node at: %s\n", node.DeploymentURL)
		logger.Printf("Node UUID: %s\n", node.ID)
		logger.Printf("Deployment UUID: %s\n", deployment.ID)
		logger.Printf("Deployment name: %s\n", deployment.Name)

		containers, err := d.sendRequestToAgent(deployment, &node, request.Namespace, request.Tag)
		if err != nil {
			return err
		}

		deployment.Status = deployments.DeploymentStatusReady
		deployment.ReplicaCount = len(containers)
		deployment.UpdatedAt = time.Now()

		logger.Printf("Deployment %s is ready with %d replicas\n", deployment.Name, len(containers))
		logger.Printf("Containers deployed: %v\n", containers)
		logger.Printf("Deployment %s completed successfully\n", deployment.Name)
		logger.Printf("Worker %d finished for deployment %s\n", id, request.DeploymentName)

		return nil
	}
}
