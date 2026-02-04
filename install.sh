#!/usr/bin/env bash
set -euo pipefail

# ============================================
# JarRun Installer
# ============================================

# Color codes for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly MAGENTA='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly BOLD='\033[1m'
readonly RESET='\033[0m'

# Icons
readonly CHECK="✅"
readonly ROCKET="🚀"
readonly WARNING="⚠️"
readonly ERROR="❌"
readonly INFO="ℹ️"

# Configuration
readonly BINARY_NAME="jarrun"
readonly INSTALL_DIR="/usr/local/bin"
readonly REPO_URL="https://github.com/Caleb-ne1/JarRun"


# functions
print_header() {
    echo -e "\n${CYAN}${BOLD}╔══════════════════════════════════════╗${RESET}"
    echo -e "${CYAN}${BOLD}║    JarRun Installer               ║${RESET}"
    echo -e "${CYAN}${BOLD}╚══════════════════════════════════════╝${RESET}\n"
}

print_success() {
    echo -e "${GREEN}${CHECK} $1${RESET}"
}

print_error() {
    echo -e "${RED}${ERROR} $1${RESET}" >&2
}

print_info() {
    echo -e "${BLUE}${INFO} $1${RESET}"
}

print_step() {
    echo -e "${MAGENTA}▶ $1${RESET}"
}

# cleanup on interrupt or error
cleanup() {
    if [[ -f "./$BINARY_NAME" ]]; then
        rm -f "./$BINARY_NAME"
        print_info "Cleaned up temporary files"
    fi
}


# main installation script
trap cleanup EXIT ERR

print_header

echo -e "${ROCKET} ${YELLOW}${BOLD}Installing JarRun...${RESET}\n"

# Detect OS and Architecture
print_step "Detecting system configuration..."

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

print_info "Operating System: $OS"
print_info "Architecture: $ARCH"

# map architecture names
case "$ARCH" in
    "x86_64")
        ARCH="amd64"
        print_info "Mapped to: $ARCH"
        ;;
    "arm64"|"aarch64")
        ARCH="arm64"
        print_info "Mapped to: $ARCH"
        ;;
    *)
        print_error "Unsupported architecture: $ARCH"
        echo -e "${YELLOW}Supported architectures: x86_64, arm64, aarch64${RESET}"
        exit 1
        ;;
esac

# construct download url
DOWNLOAD_URL="$REPO_URL/raw/main/dist/$BINARY_NAME-$OS-$ARCH"

print_step "Downloading $BINARY_NAME for $OS-$ARCH..."
echo -e "${BLUE}└─ Source: $DOWNLOAD_URL${RESET}"

# download binary
if ! curl -fsSL "$DOWNLOAD_URL" -o "$BINARY_NAME"; then
    print_error "Failed to download $BINARY_NAME"
    print_error "Please check your internet connection and try again"
    exit 1
fi

print_success "Download completed"

# make binary executable
print_step "Setting executable permissions..."
chmod +x "$BINARY_NAME"

# install to system directory
print_step "Installing to $INSTALL_DIR/..."

if ! sudo mv "$BINARY_NAME" "$INSTALL_DIR/" 2>/dev/null; then
    print_error "Installation failed. You may need to run with sudo:"
    echo -e "${YELLOW}  sudo $0${RESET}"
    exit 1
fi

print_success "Installation completed successfully!"


# Post-installation

echo -e "\n${GREEN}${BOLD}╔══════════════════════════════════════╗${RESET}"
echo -e "${GREEN}${BOLD}║          Installation Complete        ║${RESET}"
echo -e "${GREEN}${BOLD}╚══════════════════════════════════════╝${RESET}"

echo -e "\n${BOLD}Installation Details:${RESET}"
echo -e "  ${CYAN}•${RESET} Binary location: ${YELLOW}$INSTALL_DIR/$BINARY_NAME${RESET}"
echo -e "  ${CYAN}•${RESET} Repository: ${BLUE}$REPO_URL${RESET}"

echo -e "\n${GREEN}${BOLD}Happy coding! 🎉${RESET}\n"
