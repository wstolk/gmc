package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CheckoutMainBranch checks out the main branch (tries main first, then master)
func (r *Repository) CheckoutMainBranch() error {
	w, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Try main branch first
	mainRef := plumbing.ReferenceName("refs/heads/main")
	err = w.Checkout(&git.CheckoutOptions{
		Branch: mainRef,
	})

	if err != nil {
		// If main doesn't exist, try master
		masterRef := plumbing.ReferenceName("refs/heads/master")
		err = w.Checkout(&git.CheckoutOptions{
			Branch: masterRef,
		})
		if err != nil {
			return fmt.Errorf("failed to checkout main or master branch: %w", err)
		}
	}

	return nil
}

// FetchAndPrune fetches all remote branches and prunes stale remote references
func (r *Repository) FetchAndPrune(remoteName string) error {
	err := r.repo.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		Prune:      true,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to fetch from remote %s: %w", remoteName, err)
	}

	return nil
}

// GetStaleBranches returns local branches that no longer exist on the remote
func (r *Repository) GetStaleBranches(remoteName string) ([]string, error) {
	var staleBranches []string

	// Get all branches
	branches, err := r.repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	// Get remote references
	remotes, err := r.repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("failed to get remotes: %w", err)
	}

	var remoteRefs map[string]bool
	for _, remote := range remotes {
		if remote.Config().Name == remoteName {
			refs, err := remote.List(&git.ListOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to list remote references: %w", err)
			}

			remoteRefs = make(map[string]bool)
			for _, ref := range refs {
				if ref.Name().IsBranch() {
					// Store branch name without refs/heads/ prefix
					branchName := ref.Name().Short()
					remoteRefs[branchName] = true
				}
			}
			break
		}
	}

	if remoteRefs == nil {
		return nil, fmt.Errorf("remote %s not found", remoteName)
	}

	// Check each local branch
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branchName := ref.Name().Short()

			// Skip current branch (HEAD)
			head, err := r.repo.Head()
			if err == nil && head.Name() == ref.Name() {
				return nil
			}

			// Check if remote branch exists
			if _, exists := remoteRefs[branchName]; !exists {
				staleBranches = append(staleBranches, branchName)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	return staleBranches, nil
}

// DeleteBranches deletes the specified local branches
func (r *Repository) DeleteBranches(branches []string) error {
	for _, branch := range branches {
		err := r.repo.DeleteBranch(branch)
		if err != nil {
			return fmt.Errorf("failed to delete branch %s: %w", branch, err)
		}
	}
	return nil
}
