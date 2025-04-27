package gitrepo

import (
	"github.com/go-git/go-git/v5"
)

func CloneRepository(savePath string, gitRepoUrl string, branchName string) (*git.Repository, error) {
	repo, err := git.PlainClone(savePath, false, &git.CloneOptions{
		URL: gitRepoUrl,
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}
