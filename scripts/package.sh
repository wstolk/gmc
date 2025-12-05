#!/bin/bash
# GMC Packaging Script
# Creates archives and checksums for distribution

set -e

# Configuration
VERSION=${VERSION:-"v0.1.0"}
INPUT_DIR="dist"
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

create_archives() {
    log_info "Creating archives..."

    cd "$INPUT_DIR"

    # Create tar.gz archives for Unix-like systems
    for binary in gmc-linux-* gmc-darwin-*; do
        if [[ -f "$binary" ]]; then
            archive_name="${binary}.tar.gz"
            log_info "Creating $archive_name..."
            tar -czf "$archive_name" "$binary"
            log_success "Created $archive_name"
        fi
    done

    # Create zip archives for Windows
    for binary in gmc-windows-*.exe; do
        if [[ -f "$binary" ]]; then
            archive_name="${binary%.exe}.zip"
            log_info "Creating $archive_name..."
            zip -q "$archive_name" "$binary"
            log_success "Created $archive_name"
        fi
    done

    cd - >/dev/null
}

generate_checksums() {
    log_info "Generating SHA256 checksums..."

    cd "$INPUT_DIR"

    # Generate checksums for all files
    sha256sum *.tar.gz *.zip *.exe 2>/dev/null | tee SHA256SUMS

    log_success "Checksums generated in SHA256SUMS"
    cat SHA256SUMS

    cd - >/dev/null
}

main() {
    log_info "Packaging GMC $VERSION"

    if [[ ! -d "$INPUT_DIR" ]]; then
        log_error "Input directory $INPUT_DIR does not exist. Run build.sh first."
        exit 1
    fi

    create_archives
    generate_checksums

    log_success "Packaging completed successfully!"
    log_info "Files ready for release:"
    ls -la "$INPUT_DIR/"*.tar.gz "$INPUT_DIR/"*.zip "$INPUT_DIR/"SHA256SUMS 2>/dev/null || true
}

main "$@"