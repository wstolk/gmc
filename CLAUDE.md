# CLAUDE.md - AI Assistant Guide for GMC

This document provides comprehensive information about the GMC (GIT Maintenance Complete) codebase for AI assistants working on this project.

## Project Overview

**GMC** is a powerful CLI tool for comprehensive Git repository maintenance written in Go. It automates branch management and cleanup operations to keep local Git repositories clean and up-to-date.

### Core Functionality
1. **Smart Branch Management**: Automatically switches to main/master branch
2. **Remote Synchronization**: Fetches and prunes all remote branches
3. **Stale Branch Cleanup**: Safely removes local branches deleted remotely
4. **Safety First**: Dry-run mode and force confirmation for destructive operations

### Repository
- **GitHub**: https://github.com/wstolk/gmc
- **License**: MIT
- **Language**: Go 1.25.4+

## Architecture and Codebase Structure

### Directory Layout

```
gmc/
├── cmd/                    # CLI command definitions (Cobra)
│   └── root.go            # Root command and flag handling
├── internal/              # Internal packages (not importable by external code)
│   ├── git/               # Git operations and repository management
│   │   ├── operations.go  # Git operations (checkout, fetch, branch management)
│   │   ├── repository.go  # Repository wrapper and validation
│   │   └── repository_test.go  # Tests for git package
│   └── ui/                # User interface and display utilities
│       └── display.go     # Colored output, progress bars
├── scripts/               # Build, install, and packaging scripts
│   ├── build.sh          # Multi-platform build script
│   ├── install.sh        # Installation script (Linux/macOS)
│   └── package.sh        # Packaging script
├── .github/workflows/     # CI/CD pipelines
│   └── release.yml       # Release workflow (builds, tests, creates releases)
├── main.go               # Application entry point
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
├── README.md             # User-facing documentation
├── AGENTS.md             # Quick reference for AI agents
└── CLAUDE.md             # This file - comprehensive AI assistant guide
```

### Package Organization

#### `cmd/` - Command Layer
- **Purpose**: CLI interface using Cobra framework
- **Key File**: `root.go`
- **Responsibilities**:
  - Command-line argument parsing (flags: --dry-run, --verbose, --force, --remote)
  - High-level orchestration of git operations
  - Error handling and user feedback
  - Workflow execution (checkout → fetch → identify → delete)

#### `internal/git/` - Git Operations Layer
- **Purpose**: Encapsulates all git operations using go-git library
- **Key Files**:
  - `repository.go`: Repository struct and validation (`OpenRepository`, `IsValidRepository`)
  - `operations.go`: Core git operations (`CheckoutMainBranch`, `FetchAndPrune`, `GetStaleBranches`, `DeleteBranches`)
- **Design Pattern**: Repository pattern wrapping go-git functionality
- **Safety**: All operations are non-destructive except `DeleteBranches`

#### `internal/ui/` - User Interface Layer
- **Purpose**: Consistent, colored terminal output
- **Key File**: `display.go`
- **Utilities**:
  - `PrintSuccess`: Green checkmark for successful operations
  - `PrintInfo`: Blue info icon for informational messages
  - `PrintWarning`: Yellow warning icon for warnings
  - `PrintError`: Red X for errors
  - `CreateProgressBar`: Progress visualization (currently not used in main workflow)

## Key Dependencies

### Core Dependencies
- **github.com/spf13/cobra**: CLI framework for commands and flags
- **github.com/go-git/go-git/v5**: Pure Go git implementation
- **github.com/fatih/color**: Terminal color output
- **github.com/schollz/progressbar/v3**: Progress bar display

### Development Tools
- Go 1.25.4+ (required)
- `golangci-lint` (for linting)
- `gofmt` (for code formatting)
- `go vet` (for static analysis)

## Development Workflows

### Building and Running

```bash
# Build the binary
go build -o gmc .

# Run without building
go run . [flags]

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o gmc-linux-amd64 .

# Run with flags (examples)
go run . --dry-run --verbose
go run . --force --remote origin
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/git

# Run specific test function
go test -v -run TestCheckoutMainBranch ./internal/git

# Run tests with coverage
go test -cover ./...
```

### Code Quality

```bash
# Format code (must do before committing)
gofmt -w .

# Run linter
golangci-lint run

# Static analysis
go vet ./...

# Check for common issues
go vet ./... && golangci-lint run
```

## Testing Conventions

### Test Structure
- Tests use standard Go testing framework
- Tests are co-located with source files (e.g., `repository_test.go` alongside `repository.go`)
- Tests that require git use `findGitRoot()` helper to locate project root
- Tests gracefully skip when prerequisites are missing (e.g., `t.Skip()` when not in git repo)

### Test Categories
1. **Unit Tests**: Test individual functions in isolation
2. **Integration Tests**: Test with real git repositories (temporary or project repo)
3. **Skippable Tests**: Tests skip when environment doesn't support them (no git, no repo, etc.)

### Test Patterns
```go
// Pattern 1: Using project repository
func TestSomething(t *testing.T) {
    repoPath, err := findGitRoot()
    if err != nil {
        t.Skip("Skipping test: not in a git repository")
    }
    // Test with real repo...
}

// Pattern 2: Using temporary repository
func TestWithTemp(t *testing.T) {
    tempDir := t.TempDir()
    // Initialize git repo in tempDir...
    // Run tests...
}
```

### Test Robustness
- Tests must handle CI environments (GitHub Actions runs tests on ubuntu-latest)
- Tests should not assume specific branch names exist
- Tests should clean up after themselves (t.TempDir() auto-cleanup)
- Tests should work regardless of current working directory

## Code Style Guidelines

### General Go Style
- **Formatting**: Use `gofmt` (standard Go formatting, tabs not spaces)
- **Naming**:
  - PascalCase for exported identifiers (`CheckoutMainBranch`, `Repository`)
  - camelCase for unexported identifiers (`findGitRoot`, `remoteName`)
  - Package names: lowercase, single word (`git`, `ui`, `cmd`)
- **Imports**: Group in order:
  1. Standard library (e.g., `fmt`, `os`)
  2. Third-party packages (e.g., `github.com/spf13/cobra`)
  3. Local packages (e.g., `wstolk/gmc/internal/git`)

### Error Handling
- **Always return errors** using `fmt.Errorf` with `%w` verb for wrapping
- **Never panic** except for truly exceptional unrecoverable situations
- **Provide context** in error messages: `fmt.Errorf("failed to checkout branch: %w", err)`
- **Check errors** immediately after operations

### Function Design
- **Single Responsibility**: Each function should do one thing well
- **Small Functions**: Prefer multiple small functions over large complex ones
- **Clear Names**: Function names should describe what they do (verb-based)
- **Documentation**: Export functions should have doc comments

### Project-Specific Conventions

#### CLI Flag Handling
- Use `cobra.Command` with `RunE` for error handling
- Define flags in `init()` function
- Use long-form names (`--dry-run`, not `-d`)
- Provide single-letter shortcuts for common flags (`-v` for `--verbose`)

#### Git Operations Safety
- Always check if operation will be destructive
- Respect `--dry-run` flag for previewing operations
- Require `--force` flag for actual deletions
- Provide clear user feedback before destructive operations

#### User Feedback Pattern
```go
// Step pattern used throughout cmd/root.go
ui.PrintInfo("Starting operation...")
if verbose {
    fmt.Println("  Additional details...")
}
if err := doOperation(); err != nil {
    ui.PrintError("Failed to do operation: %v", err)
    return err
}
ui.PrintSuccess("Operation completed")
```

## Release Process

### Versioning
- Semantic versioning: `v1.2.3` (major.minor.patch)
- Git tags trigger releases: `git tag v1.2.3 && git push origin v1.2.3`

### Release Workflow (`.github/workflows/release.yml`)
The project uses **GoReleaser** to automate the release process.

1. **Triggered by**: Pushing a tag matching `v*.*.*`
2. **Steps**:
   - Checkout code with full history
   - Set up Go 1.21
   - Run all tests (`go test ./...`)
   - Run GoReleaser to:
     - Build binaries for multiple platforms:
       - Linux: amd64, arm64, arm
       - macOS (Darwin): amd64, arm64
       - Windows: amd64, arm64
     - Create archives (.tar.gz for Unix, .zip for Windows)
     - Generate SHA256 checksums (SHA256SUMS file)
     - Create GitHub release with formatted notes
     - Upload all artifacts

### GoReleaser Configuration (`.goreleaser.yml`)
The `.goreleaser.yml` file configures:
- **Build settings**: Cross-platform compilation with version injection
- **Archive formats**: .tar.gz for Unix, .zip for Windows
- **Checksum generation**: SHA256SUMS for verification
- **Changelog**: Auto-generated, grouped by feature/fix/enhancement
- **Release notes**: Custom header and footer with installation instructions

### Build Artifacts
- **Linux/macOS**: `.tar.gz` archives containing binary and documentation
- **Windows**: `.zip` archives containing binary and documentation
- **Naming**: `gmc-{os}-{arch}.{ext}` (e.g., `gmc-linux-amd64.tar.gz`)
- **Checksums**: `SHA256SUMS` file for verification
- **Includes**: LICENSE, README.md, CLAUDE.md, AGENTS.md

### Manual Release
```bash
# Tag a release
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3

# GitHub Actions will automatically build and publish using GoReleaser
```

### Testing Release Locally
```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser/v2@latest

# Check configuration
goreleaser check

# Build snapshot (test without releasing)
goreleaser build --snapshot --clean

# Full release dry-run
goreleaser release --snapshot --clean
```

## Common Development Tasks

### Adding a New Feature

1. **Create feature branch**:
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Implement changes**:
   - Add code in appropriate package (`internal/git/`, `internal/ui/`, or `cmd/`)
   - Follow existing patterns and conventions
   - Add error handling with context

3. **Add tests**:
   ```bash
   # Add test in corresponding _test.go file
   # Run tests to verify
   go test -v ./...
   ```

4. **Format and lint**:
   ```bash
   gofmt -w .
   golangci-lint run
   ```

5. **Commit and push**:
   ```bash
   git add .
   git commit -m "Add feature: description"
   git push origin feature/my-feature
   ```

### Fixing a Bug

1. **Write a failing test** that demonstrates the bug
2. **Fix the bug** in the source code
3. **Verify the test passes**: `go test -v -run TestName ./...`
4. **Check for regressions**: `go test ./...`
5. **Commit**: Use clear commit message describing the fix

### Adding a New CLI Flag

1. **Define flag in `cmd/root.go`**:
   ```go
   var myFlag string

   func init() {
       rootCmd.Flags().StringVar(&myFlag, "my-flag", "default", "Description")
   }
   ```

2. **Use flag in `runMaintenance` function**
3. **Update README.md** with new flag documentation
4. **Test with various flag combinations**

### Adding a New Git Operation

1. **Add method to `Repository` struct** in `internal/git/operations.go`:
   ```go
   func (r *Repository) NewOperation() error {
       // Implementation
       return nil
   }
   ```

2. **Add error handling** with wrapped errors
3. **Write tests** in `repository_test.go`
4. **Call from `cmd/root.go`** in appropriate workflow step
5. **Add UI feedback** using `ui.Print*` functions

## Important Files and Their Purposes

### Source Files

| File | Purpose | Key Functions/Types |
|------|---------|-------------------|
| `main.go` | Entry point | `main()` |
| `cmd/root.go` | CLI command definition | `runMaintenance()`, flag definitions |
| `internal/git/repository.go` | Repository abstraction | `Repository`, `OpenRepository()`, `IsValidRepository()` |
| `internal/git/operations.go` | Git operations | `CheckoutMainBranch()`, `FetchAndPrune()`, `GetStaleBranches()`, `DeleteBranches()` |
| `internal/ui/display.go` | Terminal UI | `PrintSuccess()`, `PrintError()`, `PrintInfo()`, `PrintWarning()` |

### Configuration Files

| File | Purpose |
|------|---------|
| `go.mod` | Go module definition, dependency versions |
| `go.sum` | Cryptographic checksums for dependencies |
| `.goreleaser.yml` | GoReleaser configuration for automated releases |
| `.github/workflows/release.yml` | CI/CD release pipeline using GoReleaser |

### Documentation Files

| File | Audience |
|------|----------|
| `README.md` | End users (installation, usage, features) |
| `AGENTS.md` | AI agents (quick reference, build/test commands) |
| `CLAUDE.md` | AI assistants (comprehensive guide - this file) |

### Scripts

| Script | Purpose |
|--------|---------|
| `scripts/build.sh` | Multi-platform binary builds |
| `scripts/install.sh` | Installation script for Linux/macOS |
| `scripts/package.sh` | Packaging script |

## Git Workflow Execution Order

Understanding the exact order of operations is critical when modifying the tool:

1. **Validation**: Check if current directory is a git repository
2. **Open Repository**: Initialize go-git repository object
3. **Checkout Main Branch**:
   - Get current HEAD
   - Skip if already on main/master
   - Try `main` first, fallback to `master`
4. **Fetch and Prune**:
   - Verify remote exists
   - Fetch all branches from remote
   - Prune stale remote-tracking branches
   - Skip gracefully if no remote (local-only repos)
5. **Identify Stale Branches**:
   - List all local branches
   - List all remote branches
   - Compare: find locals not on remote
   - Exclude current branch (HEAD)
6. **Delete Branches** (conditional):
   - Skip if dry-run mode
   - Skip if no stale branches found
   - Require --force flag
   - Delete each stale branch
   - Report results

## Safety Mechanisms

### Flags and Their Purposes

- `--dry-run`: Preview operations without making changes (default: false)
- `--force`: Required for actual branch deletions (default: false)
- `--verbose` (`-v`): Show detailed operation information (default: false)
- `--remote`: Specify remote name (default: "origin")

### Safety Checks

1. **Repository Validation**: Ensures we're in a git repository before any operations
2. **Main Branch Protection**: Never deletes main/master branches
3. **Current Branch Protection**: Never deletes the currently checked-out branch
4. **Remote Existence Check**: Gracefully handles missing remotes
5. **Force Flag Requirement**: Prevents accidental deletions
6. **Dry Run Mode**: Allows safe previewing of operations

## Edge Cases and Error Handling

### Handled Edge Cases

1. **No Remote**: Tool continues with local-only operations
2. **No Main/Master**: Returns clear error message
3. **Already on Main**: Skips checkout step
4. **No Stale Branches**: Exits gracefully with success message
5. **Repository Without .git**: Clear error before attempting operations

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| "Not a Git repository" | Running outside git repo | `cd` to git repo directory |
| "Failed to checkout main branch" | No main/master branch exists | Create main or master branch |
| "Remote not found" | Specified remote doesn't exist | Check remote name with `git remote -v` |
| "Refusing to delete branches without --force" | Missing --force flag | Add `--force` flag for actual deletion |

## Key Concepts

### Repository Pattern
The code uses a Repository pattern to encapsulate git operations:
- `Repository` struct wraps `*git.Repository` from go-git
- All git operations are methods on `Repository`
- This provides a clean abstraction layer for testing and maintenance

### Stale Branch Detection Algorithm
```
For each local branch:
  1. Skip if it's the current branch (HEAD)
  2. Check if remote has matching branch name
  3. If no remote match → mark as stale
Return list of stale branches
```

### Error Wrapping Strategy
Errors are consistently wrapped to provide context:
```go
if err != nil {
    return fmt.Errorf("context about what failed: %w", err)
}
```
This creates an error chain that can be unwrapped with `errors.Unwrap()`.

## Debugging Tips

### Enable Verbose Mode
Always use `--verbose` when debugging:
```bash
go run . --verbose --dry-run
```

### Check Git State
Before running GMC, inspect git state:
```bash
git status
git branch -a
git remote -v
```

### Test in Temporary Repository
Create a test repository for experimentation:
```bash
mkdir test-repo && cd test-repo
git init
git config user.name "Test"
git config user.email "test@example.com"
echo "test" > test.txt
git add . && git commit -m "Initial commit"
# Now run GMC in this directory
```

### Common Debug Commands
```bash
# See what go-git is doing
go run . --verbose 2>&1 | tee debug.log

# Test with specific remote
go run . --remote upstream --dry-run

# Test deletion safety
go run . --force --dry-run
```

## Contributing Guidelines

When contributing to GMC:

1. **Follow Go idioms**: Write idiomatic Go code
2. **Add tests**: All new features need tests
3. **Update documentation**: Keep README.md and this file current
4. **Format code**: Always run `gofmt -w .` before committing
5. **Lint**: Ensure `golangci-lint run` passes
6. **Test thoroughly**: Run full test suite with `go test ./...`
7. **Keep it simple**: Prefer simple, readable code over clever solutions
8. **Maintain backwards compatibility**: Don't break existing CLI flags or behavior

## Future Considerations

### Potential Enhancements
- Configuration file support (`.gmcrc`, `.gmc.yaml`)
- Interactive mode for branch selection
- Git hooks integration
- Multi-repository batch operations
- Branch age-based cleanup
- Protected branch patterns
- Custom cleanup rules

### Extension Points
- Additional git operations can be added to `internal/git/operations.go`
- New UI helpers can be added to `internal/ui/display.go`
- Subcommands can be added using Cobra's command structure
- Platform-specific installers can be added to `scripts/`

## Troubleshooting Development Issues

### Go Module Issues
```bash
# Reset go modules
go clean -modcache
go mod download
go mod tidy
```

### Build Issues
```bash
# Clean build cache
go clean -cache
go build -a .
```

### Test Failures in CI
- Ensure tests don't depend on local git configuration
- Use `t.Skip()` for tests requiring specific environments
- Check GitHub Actions logs for platform-specific issues

## Quick Reference

### Essential Commands
```bash
# Development
go run . --dry-run --verbose    # Test locally
go test ./...                   # Run tests
gofmt -w .                      # Format code
golangci-lint run              # Lint

# Building
go build -o gmc .              # Build binary
./gmc --help                   # Test binary

# Release
git tag v1.2.3                 # Create release tag
git push origin v1.2.3         # Trigger release workflow
```

### File Locations
- Entry point: `main.go:10`
- Main logic: `cmd/root.go:39` (`runMaintenance` function)
- Git operations: `internal/git/operations.go`
- Repository wrapper: `internal/git/repository.go`
- UI helpers: `internal/ui/display.go`
- Tests: `internal/git/repository_test.go`

---

**Last Updated**: 2025-12-16
**For**: Claude AI Assistant and other AI development tools
**Maintainer**: wstolk (https://github.com/wstolk)
