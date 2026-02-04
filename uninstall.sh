#!/usr/bin/env bash
set -euo pipefail

# ============================================
# JarRun Uninstaller
# ============================================

# Color codes for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly MAGENTA='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly BOLD='\033[1m'
readonly DIM='\033[2m'
readonly RESET='\033[0m'

# Icons
readonly TRASH="🗑️"
readonly CHECK="✅"
readonly WARNING="⚠️"
readonly INFO="ℹ️"
readonly CROSS="❌"
readonly QUESTION="❓"

# Paths
readonly BIN_PATH="/usr/local/bin/jarrun"
readonly LOG_DIR="$HOME/.jarrun/logs"
readonly CONFIG_DIR="$HOME/.jarrun/config"


# functions
print_header() {
    echo -e "\n${RED}${BOLD}╔══════════════════════════════════════╗${RESET}"
    echo -e "${RED}${BOLD}║    JarRun CLI Uninstaller             ║${RESET}"
    echo -e "${RED}${BOLD}╚══════════════════════════════════════╝${RESET}\n"
}

print_warning() {
    echo -e "${YELLOW}${WARNING} $1${RESET}"
}

print_success() {
    echo -e "${GREEN}${CHECK} $1${RESET}"
}

print_error() {
    echo -e "${RED}${CROSS} $1${RESET}" >&2
}

print_info() {
    echo -e "${BLUE}${INFO} $1${RESET}"
}

print_step() {
    echo -e "${MAGENTA}▶ $1${RESET}"
}

print_removed() {
    echo -e "${GREEN}  ✓${RESET} Removed: ${DIM}$1${RESET}"
}

print_skipped() {
    echo -e "${YELLOW}  ○${RESET} Not found: ${DIM}$1${RESET}"
}

ask_for_confirmation() {
    echo -e "\n${YELLOW}${WARNING} This will uninstall JarRun and remove:${RESET}"
    echo -e "  ${CYAN}•${RESET} Binary: ${BOLD}$BIN_PATH${RESET}"
    echo -e "  ${CYAN}•${RESET} Logs directory: ${BOLD}$LOG_DIR${RESET}"
    echo -e "  ${CYAN}•${RESET} Configuration: ${BOLD}$CONFIG_DIR${RESET}"
    
    echo -e "\n${YELLOW}${QUESTION} Are you sure you want to continue?${RESET}"
    read -p "Type 'yes' to confirm: " -r CONFIRMATION
    
    if [[ ! "$CONFIRMATION" =~ ^[Yy][Ee][Ss]$ ]]; then
        echo -e "\n${BLUE}${INFO} Uninstallation cancelled.${RESET}"
        exit 0
    fi
}


# main uninstallation script
print_header

echo -e "${TRASH} ${RED}${BOLD}Uninstalling JarRun...${RESET}"

ask_for_confirmation

echo -e "\n${CYAN}${BOLD}Removing Components:${RESET}"

# remove binary
print_step "Binary file..."
if [ -f "$BIN_PATH" ]; then
    if sudo rm -f "$BIN_PATH" 2>/dev/null; then
        print_removed "$BIN_PATH"
    else
        print_error "Failed to remove binary. You may need to run with sudo."
        echo -e "${YELLOW}Try: sudo $(basename "$0")${RESET}"
        exit 1
    fi
else
    print_skipped "$BIN_PATH"
fi

# remove logs directory
print_step "Logs directory..."
if [ -d "$LOG_DIR" ]; then
    if rm -rf "$LOG_DIR"; then
        print_removed "$LOG_DIR"
    else
        print_warning "Could not fully remove logs directory"
    fi
else
    print_skipped "$LOG_DIR"
fi

# remove config directory
print_step "Configuration directory..."
if [ -d "$CONFIG_DIR" ]; then
    echo -e "${YELLOW}  ⚠️  Configuration directory contains user settings${RESET}"
    echo -e "${BLUE}  Contents:${RESET}"
    if command -v tree &>/dev/null && [ -d "$CONFIG_DIR" ]; then
        tree "$CONFIG_DIR" -L 2 2>/dev/null || ls -la "$CONFIG_DIR/"
    else
        ls -la "$CONFIG_DIR/" 2>/dev/null || echo "    (cannot list contents)"
    fi
    
    read -p "$(echo -e ${YELLOW}"Remove configuration? [y/N]: "${RESET})" -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if rm -rf "$CONFIG_DIR"; then
            print_removed "$CONFIG_DIR"
        else
            print_warning "Could not fully remove config directory"
        fi
    else
        echo -e "${BLUE}  Configuration preserved at: $CONFIG_DIR${RESET}"
    fi
else
    print_skipped "$CONFIG_DIR"
fi

echo -e "\n${CYAN}${BOLD}Verification:${RESET}"

# Check if binary still exists
if command -v jarrun &>/dev/null || [ -f "$BIN_PATH" ]; then
    print_warning "Binary might still be accessible"
    echo -e "${YELLOW}  Run 'which jarrun' to check installation${RESET}"
else
    print_success "Binary successfully removed from system PATH"
fi

echo -e "\n${RED}${BOLD}╔══════════════════════════════════════╗${RESET}"
echo -e "${RED}${BOLD}║      Uninstallation Complete          ║${RESET}"
echo -e "${RED}${BOLD}╚══════════════════════════════════════╝${RESET}"

echo -e "\n${GREEN}${BOLD}JarRun has been uninstalled successfully!${RESET}"


echo -e "\n${DIM}Thank you for using JarRun! 👋${RESET}\n"

