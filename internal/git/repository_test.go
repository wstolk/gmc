package git

import (
	"testing"
)

func TestIsValidRepository(t *testing.T) {
	// Test with project root directory (should be valid)
	repoPath := "../../" // Go up two levels from internal/git to project root
	if !IsValidRepository(repoPath) {
		t.Error("Project root should be a valid Git repository")
	}

	// Test with temp directory (should be invalid)
	tempDir := t.TempDir()
	if IsValidRepository(tempDir) {
		t.Error("Temp directory should not be a valid Git repository")
	}
}

func TestOpenRepository(t *testing.T) {
	repoPath := "../../" // Go up two levels from internal/git to project root

	repo, err := OpenRepository(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repository: %v", err)
	}

	if repo.path != repoPath {
		t.Errorf("Expected path %s, got %s", repoPath, repo.path)
	}
}

func TestCheckoutMainBranch(t *testing.T) {
	repoPath := "../../" // Go up two levels from internal/git to project root

	repo, err := OpenRepository(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repository: %v", err)
	}

	// Should not error since we're already on main
	err = repo.CheckoutMainBranch()
	if err != nil {
		t.Errorf("Failed to checkout main branch: %v", err)
	}
}

func TestGetStaleBranches_NoRemote(t *testing.T) {
	repoPath := "../../" // Go up two levels from internal/git to project root

	repo, err := OpenRepository(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repository: %v", err)
	}

	// Should return empty list when no remote exists
	stale, err := repo.GetStaleBranches("origin")
	if err != nil {
		t.Errorf("Should not error when no remote exists: %v", err)
	}

	if len(stale) != 0 {
		t.Errorf("Expected no stale branches, got %d", len(stale))
	}
}
