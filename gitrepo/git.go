package gitrepo

import (
	"context"
	"errors"
	"fmt"
	backend "glider"
	"os"

	"github.com/go-git/go-git/v5"
)

const (
	RepoUrlKey  = "repoUrl"
	RepoNameKey = "repoName"
	TempDirPath = "/tmp/repos"
)

type GitSourceCodeService struct{}

func NewGitsSourceCodeService() *GitSourceCodeService {
	return &GitSourceCodeService{}
}

func (s *GitSourceCodeService) Fetch(ctx context.Context, deploymentMetadata backend.DeploymentMetadata) (string, error) {
	if err := validateDeploymentMetadata(deploymentMetadata); err != nil {
		return "", err
	}

	repoUrl := deploymentMetadata[RepoUrlKey].(string)
	repoName := deploymentMetadata[RepoNameKey].(string)

	if _, err := os.Stat(TempDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(TempDirPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	path := fmt.Sprintf("%s/%s", TempDirPath, repoName)

	_, err := git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL: repoUrl,
	})
	if err != nil {
		return "", err
	}

	return path, nil
}

func validateDeploymentMetadata(deploymentMetadata backend.DeploymentMetadata) error {
	if _, ok := deploymentMetadata[RepoUrlKey]; !ok {
		return errors.Join(backend.ErrInvalidInput, errors.New(fmt.Sprintf("'%s' key is required when deploying with git", RepoUrlKey)))
	}
	if _, ok := deploymentMetadata[RepoNameKey]; !ok {
		return errors.Join(backend.ErrInvalidInput, errors.New(fmt.Sprintf("'%s' key is required when deploying with git", RepoNameKey)))
	}

	return nil
}
