package api

import (
	"context"
	"fmt"
	"glider/deployments"
	"glider/gitrepo"
	"glider/images"
	"glider/workerpool"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type DeployRequest struct {
	DeploymentName string `json:"deploymentName"`
	Replicas       int    `json:"replicas"`
	GithubRepo     string `json:"githubRepo"`
	GithubBranch   string `json:"githubBranch"`
	GithubToken    string `json:"githubToken"`
	Tag            string `json:"tag"`
	Namespace      string `json:"namespace"`
	DockerfilePath string `json:"dockerfilePath"`
}

type DeployResponse struct {
	DeploymentName string `json:"deploymentName"`
	Status         string `json:"status"`
	Replicas       int    `json:"replicas"`
	GithubRepo     string `json:"githubRepo"`
	GithubBranch   string `json:"githubBranch"`
}

type DeployParams struct {
	Deployment *deployments.Deployment
	Request    DeployRequest
}

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

	d.workerPool.Submit(d.deploy(c, deployment, request))

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

	if err := images.StoreImage(ctx, cli, deploymentName, namespace, tag); err != nil {
		return err
	}

	return nil
}

func (d DeployHandlers) deploy(ctx context.Context, deployment *deployments.Deployment, request DeployRequest) workerpool.Task {
	return func(ctx context.Context, id int) error {
		path := "./tmp/clones/" + request.DeploymentName

		err := d.uploadImage(ctx, d.dockerCli, path, request.GithubRepo, request.GithubBranch, request.DeploymentName, request.DockerfilePath, request.Tag, request.Namespace)
		if err != nil {
			return err
		}

		// TODO
		// get list of available agents
		// choose agent node to useb based on database based on best fit method (or random at start) and based on cloud vs on-prem preferences later
		// send request to agent to deploy the new deployment
		// set deployment status based on the result of the previus step
		// send an update to websockets if needed to update the UI live
		// store deployment in database

		return nil
	}
}
