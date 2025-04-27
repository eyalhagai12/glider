package api

import (
	"context"
	"fmt"
	"glider/deployments"
	"glider/gitrepo"
	"glider/images"
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

type DeployHandlers struct {
	db *sqlx.DB
}

func NewDeployHandlers(db *sqlx.DB) DeployHandlers {
	return DeployHandlers{
		db: db,
	}
}

func (d DeployHandlers) Deploy(c *gin.Context, request DeployRequest) (*deployments.Deployment, error) {
	path := "./tmp/clones/" + request.DeploymentName

	deployment := deployments.NewDeployment(request.DeploymentName, request.GithubRepo, request.GithubBranch, request.Replicas)
	err := deployments.StoreDeployment(d.db, deployment)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to store deployment: %v", err))
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create docker client: %v", err))
	}
	defer cli.Close()

	err = d.uploadImage(c, cli, path, request.GithubRepo, request.GithubBranch, request.DeploymentName, request.DockerfilePath, request.Tag, request.Namespace)
	if err != nil {
		return nil, c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to upload image: %v", err))
	}

	return deployment, nil
}

func (d DeployHandlers) uploadImage(ctx context.Context, cli *client.Client, path string, repo string, branch string, deploymentName string, dockerFilePath string, tag string, namespace string) error {
	_, err := gitrepo.CloneRepository(path, repo, branch)
	if err != nil {
		return err
	}

	image, err := images.BuildImage(ctx, cli, deploymentName, path, dockerFilePath, tag)
	if err != nil {
		return err
	}
	defer image.Body.Close()

	if err := images.StoreImage(ctx, cli, deploymentName, namespace, tag); err != nil {
		return err
	}

	err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}
