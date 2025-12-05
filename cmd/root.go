package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"wstolk/gmc/internal/git"
	"wstolk/gmc/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "gmc",
	Short: "GIT Maintenance Complete - Clean up your Git repositories",
	Long: `GMC performs comprehensive Git repository maintenance:

1. Checkout main branch
2. Pull all remote branches
3. Cleanup stale local branches that no longer exist remotely

This tool helps keep your local Git repositories clean and up-to-date.`,
	RunE: runMaintenance,
}

var (
	dryRun  bool
	verbose bool
	remote  string
	force   bool
)

func init() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVar(&remote, "remote", "origin", "Remote name to use (default: origin)")
	rootCmd.Flags().BoolVar(&force, "force", false, "Force deletion of branches with uncommitted changes")
}

func runMaintenance(cmd *cobra.Command, args []string) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if we're in a Git repository
	if !git.IsValidRepository(cwd) {
		ui.PrintError("Not a Git repository: %s", cwd)
		return fmt.Errorf("not a Git repository")
	}

	ui.PrintInfo("Starting Git maintenance in: %s", cwd)

	// Open repository
	repo, err := git.OpenRepository(cwd)
	if err != nil {
		ui.PrintError("Failed to open repository: %v", err)
		return err
	}

	// Step 1: Checkout main branch
	ui.PrintInfo("Checking out main branch...")
	if verbose {
		fmt.Println("  Looking for main or master branch...")
	}
	if err := repo.CheckoutMainBranch(); err != nil {
		ui.PrintError("Failed to checkout main branch: %v", err)
		return err
	}
	ui.PrintSuccess("Checked out main branch")

	// Step 2: Fetch and prune (skip if no remote)
	ui.PrintInfo("Fetching from remote '%s' with pruning...", remote)
	if verbose {
		fmt.Println("  This will update local remote-tracking branches...")
	}
	if err := repo.FetchAndPrune(remote); err != nil {
		ui.PrintWarning("Skipping fetch/prune: %v", err)
		if verbose {
			fmt.Println("  No remote repository found, proceeding with local cleanup only...")
		}
	} else {
		ui.PrintSuccess("Fetched and pruned remote branches")
	}

	// Step 3: Identify stale branches
	ui.PrintInfo("Identifying stale local branches...")
	if verbose {
		fmt.Println("  Comparing local branches with remote branches...")
	}
	staleBranches, err := repo.GetStaleBranches(remote)
	if err != nil {
		ui.PrintError("Failed to identify stale branches: %v", err)
		return err
	}

	if len(staleBranches) == 0 {
		ui.PrintSuccess("No stale branches found")
		return nil
	}

	// Show what would be deleted
	ui.PrintWarning("Found %d stale local branch(es):", len(staleBranches))
	for _, branch := range staleBranches {
		fmt.Printf("  - %s\n", branch)
	}

	// Delete branches (unless dry run)
	if !dryRun {
		if !force && len(staleBranches) > 0 {
			ui.PrintWarning("Use --force to actually delete branches, or --dry-run to preview")
			return fmt.Errorf("refusing to delete branches without --force flag")
		}

		ui.PrintInfo("Deleting stale branches...")
		if verbose {
			fmt.Printf("  Deleting %d branch(es)...\n", len(staleBranches))
		}
		if err := repo.DeleteBranches(staleBranches); err != nil {
			ui.PrintError("Failed to delete branches: %v", err)
			return err
		}
		ui.PrintSuccess("Deleted %d stale branch(es)", len(staleBranches))
	} else {
		ui.PrintInfo("Dry run: would delete %d branch(es)", len(staleBranches))
	}

	ui.PrintSuccess("Git maintenance completed successfully!")
	return nil
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
