# GMC: GIT Maintenance Complete

A powerful CLI tool for comprehensive Git repository maintenance. Keep your local Git repositories clean and up-to-date with automated branch management and cleanup.

## Features

- **🔄 Smart Branch Management**: Automatically switches to main/master branch
- **📥 Remote Synchronization**: Fetches and prunes all remote branches
- **🧹 Stale Branch Cleanup**: Safely removes local branches deleted remotely
- **🛡️ Safety First**: Dry-run mode and force confirmation for destructive operations
- **🌍 Cross-Platform**: Works on Linux, macOS, and Windows
- **🎨 Beautiful CLI**: Colored output with progress indicators
- **⚡ Fast & Efficient**: Optimized for performance with minimal dependencies

## Installation

### Quick Install (Recommended)

#### Linux/macOS
```bash
curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash
```

#### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.ps1 | iex
```

### Install Specific Version

```bash
# Install specific version
curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash -s -- --version v0.1.0

# Install with GPG verification (requires gpg)
curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash -s -- --verify
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/wstolk/gmc/releases)
2. Extract the archive
3. Move the binary to a directory in your PATH

#### Linux/macOS
```bash
# Download and install manually
wget https://github.com/wstolk/gmc/releases/download/v0.1.0/gmc-linux-amd64.tar.gz
tar -xzf gmc-linux-amd64.tar.gz
sudo mv gmc /usr/local/bin/
```

#### Windows
```powershell
# Download and install manually
Invoke-WebRequest -Uri "https://github.com/wstolk/gmc/releases/download/v0.1.0/gmc-windows-amd64.zip" -OutFile "gmc.zip"
Expand-Archive -Path "gmc.zip" -DestinationPath "."
Move-Item -Path "gmc.exe" -Destination "$env:USERPROFILE\bin\gmc.exe"
```

### Package Managers

#### macOS (Homebrew)
```bash
brew tap wstolk/gmc
brew install gmc
```

#### Build from Source

**Requirements**: Go 1.21+

```bash
git clone https://github.com/wstolk/gmc
cd gmc
go build -o gmc .
sudo mv gmc /usr/local/bin/
```

## Usage

### Basic Usage

```bash
# Preview maintenance operations (safe)
gmc --dry-run

# Run maintenance with verbose output
gmc --verbose

# Force cleanup (required for actual deletions)
gmc --force

# Specify custom remote
gmc --remote upstream
```

### Advanced Examples

```bash
# Dry run with verbose output
gmc --dry-run --verbose

# Maintenance on a repository with custom remote
gmc --remote origin --force

# Get help
gmc --help
```

### What GMC Does

1. **Checkout Main Branch**: Switches to `main` or `master` branch
2. **Fetch & Prune**: Updates all remote branches and removes stale remote references
3. **Identify Stale Branches**: Finds local branches that no longer exist on remote
4. **Safe Cleanup**: Removes stale branches (with confirmation in dry-run mode)

### Output Examples

#### Successful Run
```
ℹ Starting Git maintenance in: /path/to/repo
ℹ Checking out main branch...
✓ Checked out main branch
ℹ Fetching from remote 'origin' with pruning...
✓ Fetched and pruned remote branches
ℹ Identifying stale local branches...
⚠ Found 2 stale local branch(es):
  - feature/old-feature
  - bugfix/fixed-bug
ℹ Dry run: would delete 2 branch(es)
```

#### Error Handling
```
✗ Failed to checkout main branch: reference not found
```

## Security

- **Checksum Verification**: All downloads are verified using SHA256 checksums
- **GPG Signatures**: Optional GPG signature verification for enhanced security
- **No Telemetry**: GMC doesn't collect or transmit any data
- **Local Operation**: All operations are performed locally on your system

### Verifying Downloads

The installer automatically verifies downloads, but you can manually verify:

```bash
# Download checksums
curl -fsSL https://github.com/wstolk/gmc/releases/download/v0.1.0/SHA256SUMS -o SHA256SUMS

# Verify
sha256sum -c SHA256SUMS
```

## Configuration

GMC works out-of-the-box with standard Git configurations. It automatically:

- Detects main branch (`main` or `master`)
- Uses `origin` as default remote
- Respects your Git authentication settings

## Troubleshooting

### "Not a Git repository"
```
✗ Not a Git repository: /current/directory
```
**Solution**: Run GMC from within a Git repository directory.

### "Remote not found"
```
⚠ Skipping fetch/prune: remote origin not found
```
**Solution**: This is normal for local-only repositories. GMC will continue with local branch cleanup.

### "Permission denied"
```
✗ Failed to install: permission denied
```
**Solution**: The installer needs sudo access for `/usr/local/bin`. Use `sudo` or install to `~/bin`.

### Command not found after installation
```bash
gmc: command not found
```
**Solution**: Add the installation directory to your PATH:
```bash
export PATH="$HOME/bin:$PATH"
# Or add to your shell profile
```

## Contributing

Contributions are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/wstolk/gmc/issues)
- **Discussions**: [GitHub Discussions](https://github.com/wstolk/gmc/discussions)
- **Documentation**: [GitHub Wiki](https://github.com/wstolk/gmc/wiki)

---

**Made with ❤️ for the Git community**
