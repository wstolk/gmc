#!/bin/bash
# GMC Build Script
# Builds GMC binaries for multiple platforms

set -e

# Configuration
VERSION=${VERSION:-"v0.1.0"}
OUTPUT_DIR="dist"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# Build matrix
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

main() {
    log_info "Building GMC $VERSION for multiple platforms"

    # Clean output directory
    rm -rf "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"

    # Build for each platform
    for PLATFORM in "${PLATFORMS[@]}"; do
        GOOS=${PLATFORM%/*}
        GOARCH=${PLATFORM#*/}

        BINARY_NAME="gmc"
        if [ "$GOOS" = "windows" ]; then
            BINARY_NAME="gmc.exe"
        fi

        OUTPUT_NAME="gmc-${GOOS}-${GOARCH}"
        if [ "$GOOS" = "windows" ]; then
            OUTPUT_NAME="gmc-${GOOS}-${GOARCH}.exe"
        fi

        log_info "Building for $GOOS/$GOARCH..."

        # Build with version info
        GOOS=$GOOS GOARCH=$GOARCH go build \
            -ldflags "-X main.version=$VERSION -s -w" \
            -o "$OUTPUT_DIR/$OUTPUT_NAME" \
            .

        log_success "Built $OUTPUT_NAME"
    done

    log_success "All binaries built successfully in $OUTPUT_DIR/"
    ls -la "$OUTPUT_DIR/"
}

main "$@"