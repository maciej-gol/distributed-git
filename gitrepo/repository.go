package gitrepo

import (
	"fmt"
	g2g "gopkg.in/libgit2/git2go.v24"
)

type GitRepository interface{}

type gitRepository struct {
	repository *g2g.Repository
}

func OpenGitRepository(path string) (GitRepository, error) {
	repo, err := g2g.OpenRepository(path)
	if err != nil {
		return nil, fmt.Errorf("error opening repository: %s", err.Error())
	}
	return &gitRepository{repository: repo}, nil
}

func InitGitRepository(path string) (GitRepository, error) {
	repo, err := g2g.InitRepository(path, true)
	if err != nil {
		return nil, fmt.Errorf("error creating repository: %s", err.Error())
	}
	return &gitRepository{repository: repo}, nil
}
