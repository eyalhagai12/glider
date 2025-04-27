package gitrepo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneRepository(savePath string, gitRepoUrl string, branchName string) (*git.Repository, error) {
	repo, err := git.PlainClone(savePath, false, &git.CloneOptions{
		URL:           gitRepoUrl,
		ReferenceName: plumbing.ReferenceName(branchName),
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}
