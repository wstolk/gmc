package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// findGitRoot finds the git repository root by walking up the directory tree
func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}

func TestIsValidRepository(t *testing.T) {
	// Test with project root directory (should be valid)
	repoPath, err := findGitRoot()
	if err != nil {
		t.Skip("Skipping test: not in a git repository")
	}

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
	repoPath, err := findGitRoot()
	if err != nil {
		t.Skip("Skipping test: not in a git repository")
	}

	repo, err := OpenRepository(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repository: %v", err)
	}

	if repo.path != repoPath {
		t.Errorf("Expected path %s, got %s", repoPath, repo.path)
	}
}

func TestCheckoutMainBranch(t *testing.T) {
	repoPath, err := findGitRoot()
	if err != nil {
		t.Skip("Skipping test: not in a git repository")
	}

	repo, err := OpenRepository(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repository: %v", err)
	}

	// Should not error since we're already on main or master
	err = repo.CheckoutMainBranch()
	if err != nil {
		t.Errorf("Failed to checkout main branch: %v", err)
	}
}

func TestGetStaleBranches_NoRemote(t *testing.T) {
	repoPath, err := findGitRoot()
	if err != nil {
		t.Skip("Skipping test: not in a git repository")
	}

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

// Test with a temporary git repository
func TestWithTempRepository(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Initialize a git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Skip("Skipping test: git not available")
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	cmd.Run()

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	cmd.Run()

	// Create and commit a file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Now test our functions
	if !IsValidRepository(tempDir) {
		t.Error("Temp directory should be a valid Git repository")
	}

	repo, err := OpenRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to open temp repository: %v", err)
	}

	// Test checkout main branch
	err = repo.CheckoutMainBranch()
	if err != nil {
		t.Errorf("Failed to checkout main branch in temp repo: %v", err)
	}

	// Test get stale branches (should be empty)
	stale, err := repo.GetStaleBranches("origin")
	if err != nil {
		t.Errorf("Should not error when no remote exists: %v", err)
	}

	if len(stale) != 0 {
		t.Errorf("Expected no stale branches in temp repo, got %d", len(stale))
	}
}
