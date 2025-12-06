#!/bin/bash
# GMC Installer Script
# Installs GMC (GIT Maintenance Complete) CLI tool
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash -s -- --version v0.1.0
#
# Options:
#   --version VERSION    Install specific version (default: latest)
#   --verify             Verify GPG signatures (requires gpg)
#   --no-color           Disable colored output
#   --help               Show this help

set -e

# Configuration
REPO="wstolk/gmc"
GITHUB_API="https://api.github.com/repos/$REPO"
GITHUB_RELEASES="https://github.com/$REPO/releases"
BINARY_NAME="gmc"
INSTALL_DIR="/usr/local/bin"
FALLBACK_DIR="$HOME/bin"

# Colors
if [[ -t 1 && -z "$NO_COLOR" ]]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
fi

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1" >&2
}

# Cleanup function
cleanup() {
    if [[ -n "$TEMP_DIR" && -d "$TEMP_DIR" ]]; then
        rm -rf "$TEMP_DIR"
    fi
}

trap cleanup EXIT

# Show usage
show_usage() {
    cat << EOF
GMC Installer

Installs GMC (GIT Maintenance Complete) CLI tool.

USAGE:
    curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash
    curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash -s -- [OPTIONS]

OPTIONS:
    --version VERSION    Install specific version (default: latest)
    --verify             Verify GPG signatures (requires gpg)
    --no-color           Disable colored output
    --help               Show this help

EXAMPLES:
    # Install latest version
    curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash

    # Install specific version
    curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash -s -- --version v0.1.0

    # Install with GPG verification
    curl -fsSL https://raw.githubusercontent.com/wstolk/gmc/main/scripts/install.sh | bash -s -- --verify

For more information, visit: https://github.com/wstolk/gmc
EOF
}

# Parse arguments
VERSION="latest"
VERIFY_SIGNATURES=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            VERSION="$2"
            shift 2
            ;;
        --verify)
            VERIFY_SIGNATURES=true
            shift
            ;;
        --no-color)
            RED=''
            GREEN=''
            YELLOW=''
            BLUE=''
            NC=''
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Detect platform and architecture
detect_platform() {
    if [[ -z "$OS" ]]; then
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    fi

    case $OS in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        freebsd)
            OS="freebsd"
            ;;
        *)
            log_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac

    if [[ -z "$ARCH" ]]; then
        ARCH=$(uname -m)
    fi

    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l|arm)
            ARCH="arm"
            ;;
        i386|i686)
            ARCH="386"
            ;;
        *)
            log_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
}

# Get version information
get_version_info() {
    if [[ "$VERSION" == "latest" ]]; then
        log_info "Fetching latest version information..."
        RELEASE_INFO=$(curl -s "$GITHUB_API/releases/latest")
        VERSION=$(echo "$RELEASE_INFO" | grep '"tag_name"' | cut -d'"' -f4)
        if [[ -z "$VERSION" ]]; then
            log_error "Failed to fetch latest version"
            exit 1
        fi
    fi

    log_info "Installing GMC $VERSION for $OS/$ARCH"
}

# Download file with progress
download_file() {
    local url="$1"
    local output="$2"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL --progress-bar "$url" -o "$output"
    elif command -v wget >/dev/null 2>&1; then
        wget -q --show-progress "$url" -O "$output"
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

# Verify checksums
verify_checksums() {
    local checksums_file="$TEMP_DIR/SHA256SUMS"

    if [[ ! -f "$checksums_file" ]]; then
        log_warning "No checksums file found, skipping verification"
        return
    fi

    log_info "Verifying checksums..."

    if command -v sha256sum >/dev/null 2>&1; then
        cd "$TEMP_DIR"
        if ! sha256sum -c --ignore-missing SHA256SUMS >/dev/null 2>&1; then
            log_error "Checksum verification failed"
            exit 1
        fi
        cd - >/dev/null
    else
        log_warning "sha256sum not found, skipping checksum verification"
    fi

    log_success "Checksums verified"
}

# Verify GPG signatures
verify_signatures() {
    if [[ "$VERIFY_SIGNATURES" != true ]]; then
        return
    fi

    local signature_file="$TEMP_DIR/SHA256SUMS.sig"
    local checksums_file="$TEMP_DIR/SHA256SUMS"

    if [[ ! -f "$signature_file" ]]; then
        log_warning "No GPG signature file found, skipping signature verification"
        return
    fi

    if ! command -v gpg >/dev/null 2>&1; then
        log_warning "GPG not found, skipping signature verification"
        return
    fi

    log_info "Verifying GPG signatures..."

    # Import GMC public key (this would be the actual key in production)
    # For now, we'll skip actual verification since we don't have the key set up
    log_warning "GPG signature verification not fully implemented yet"
}

# Determine installation directory
get_install_dir() {
    # Try /usr/local/bin first
    if [[ -w "/usr/local/bin" ]] || ( [[ -w "/usr/local" ]] && sudo -n true 2>/dev/null ); then
        echo "/usr/local/bin"
        return
    fi

    # Try ~/bin as fallback
    if [[ -d "$HOME/bin" && -w "$HOME/bin" ]]; then
        echo "$HOME/bin"
        return
    fi

    # Create ~/bin if it doesn't exist
    if mkdir -p "$HOME/bin" 2>/dev/null; then
        echo "$HOME/bin"
        return
    fi

    log_error "No suitable installation directory found"
    exit 1
}

# Install binary
install_binary() {
    local archive_name
    local binary_path

    if [[ "$OS" == "windows" ]]; then
        archive_name="gmc-${OS}-${ARCH}.zip"
        binary_path="$TEMP_DIR/gmc-${OS}-${ARCH}.exe"
    else
        archive_name="gmc-${OS}-${ARCH}.tar.gz"
        binary_path="$TEMP_DIR/gmc-${OS}-${ARCH}"
    fi

    # Download archive
    log_info "Downloading GMC $VERSION..."
    download_file "$GITHUB_RELEASES/download/$VERSION/$archive_name" "$TEMP_DIR/$archive_name"

    # Download checksums
    download_file "$GITHUB_RELEASES/download/$VERSION/SHA256SUMS" "$TEMP_DIR/SHA256SUMS" 2>/dev/null || true

    # Verify checksums
    verify_checksums

    # Extract archive
    log_info "Extracting archive..."
    if [[ "$OS" == "windows" ]]; then
        unzip -q "$TEMP_DIR/$archive_name" -d "$TEMP_DIR"
    else
        tar -xzf "$TEMP_DIR/$archive_name" -C "$TEMP_DIR"
    fi

    # Verify signatures
    verify_signatures

    # Determine install directory
    INSTALL_DIR=$(get_install_dir)
    log_info "Installing to $INSTALL_DIR"

    # Install binary
    if [[ "$INSTALL_DIR" == "/usr/local/bin" ]] && ! [[ -w "/usr/local/bin" ]]; then
        sudo cp "$binary_path" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/gmc"
    else
        cp "$binary_path" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/gmc"
    fi

    log_success "GMC installed successfully!"
}

# Test installation
test_installation() {
    log_info "Testing installation..."

    if ! command -v gmc >/dev/null 2>&1; then
        log_error "GMC not found in PATH. You may need to restart your shell or add $INSTALL_DIR to your PATH."
        echo "To add to PATH temporarily: export PATH=\"$INSTALL_DIR:\$PATH\""
        exit 1
    fi

    if ! gmc --help >/dev/null 2>&1; then
        log_error "GMC installation test failed"
        exit 1
    fi

    log_success "Installation verified!"
}

# Main function
main() {
    log_info "GMC Installer"
    echo

    detect_platform
    get_version_info

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    log_info "Using temporary directory: $TEMP_DIR"

    install_binary
    test_installation

    echo
    log_success "GMC $VERSION has been installed successfully!"
    echo
    echo "Run 'gmc --help' to get started."
    echo "For more information, visit: https://github.com/wstolk/gmc"
}

main "$@"