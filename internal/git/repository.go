package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

// Repository represents a Git repository with operations
type Repository struct {
	repo *git.Repository
	path string
}

// OpenRepository opens a Git repository at the given path
func OpenRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository at %s: %w", path, err)
	}

	return &Repository{
		repo: repo,
		path: path,
	}, nil
}

// IsValidRepository checks if the current directory is a valid Git repository
func IsValidRepository(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}
