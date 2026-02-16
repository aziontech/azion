#!/bin/bash
# Azion CLI installer script
# Usage: curl -fsSL https://cli.azion.app/install.sh | bash
#
# Environment variables:
#   AZION_VERSION     - Pin a specific version (default: latest)
#   AZION_INSTALL_DIR - Custom install directory (default: $HOME/.azion/bin)

set -euo pipefail

# --- Error handling -----------------------------------------------------------

trap 'on_error $LINENO' ERR

on_error() {
    error "installation failed at line $1"
    exit 1
}

# --- Color output (terminal-aware) --------------------------------------------

setup_colors() {
    if [ -t 1 ] && command -v tput >/dev/null 2>&1 && tput colors >/dev/null 2>&1; then
        RED=$(tput setaf 1)
        GREEN=$(tput setaf 2)
        YELLOW=$(tput setaf 3)
        BOLD=$(tput bold)
        RESET=$(tput sgr0)
    else
        RED=""
        GREEN=""
        YELLOW=""
        BOLD=""
        RESET=""
    fi
}

info() {
    printf '%s[info]%s %s\n' "${GREEN}" "${RESET}" "$1"
}

warn() {
    printf '%s[warn]%s %s\n' "${YELLOW}" "${RESET}" "$1"
}

error() {
    printf '%s[error]%s %s\n' "${RED}" "${RESET}" "$1" >&2
}

# --- OS and architecture detection --------------------------------------------

detect_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')

    case "$os" in
        linux*)  PLATFORM="linux" ;;
        darwin*) PLATFORM="darwin" ;;
        freebsd*) PLATFORM="freebsd" ;;
        mingw*|msys*|cygwin*)
            error "Windows is not supported. Please use WSL or download from https://github.com/aziontech/azion/releases"
            exit 1
            ;;
        *)
            error "unsupported operating system: $os"
            exit 1
            ;;
    esac
}

detect_arch() {
    local arch
    arch=$(uname -m)

    case "$arch" in
        x86_64|amd64)   ARCH="amd64" ;;
        aarch64|arm64)   ARCH="arm64" ;;
        armv7l|armv7)    ARCH="armv7" ;;
        i686|i386|i586)  ARCH="386" ;;
        ppc64|ppc64le)   ARCH="ppc64" ;;
        *)
            error "unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# --- HTTP client --------------------------------------------------------------

detect_http_client() {
    if command -v curl >/dev/null 2>&1; then
        HTTP_CLIENT="curl"
    elif command -v wget >/dev/null 2>&1; then
        HTTP_CLIENT="wget"
    else
        error "either curl or wget is required"
        exit 1
    fi
}

http_get() {
    local url="$1"
    if [ "$HTTP_CLIENT" = "curl" ]; then
        curl -fsSL "$url"
    else
        wget -qO- "$url"
    fi
}

http_download() {
    local url="$1"
    local output="$2"
    if [ "$HTTP_CLIENT" = "curl" ]; then
        curl -fsSL -o "$output" "$url"
    else
        wget -q -O "$output" "$url"
    fi
}

# --- Version resolution -------------------------------------------------------

resolve_version() {
    if [ -n "${AZION_VERSION:-}" ]; then
        VERSION="$AZION_VERSION"
        info "using pinned version: $VERSION"
    else
        info "fetching latest version from GitHub..."
        local api_response
        api_response=$(http_get "https://api.github.com/repos/aziontech/azion/releases/latest")
        VERSION=$(printf '%s' "$api_response" | grep '"tag_name"' | sed -E 's/.*"tag_name":[ ]*"([^"]+)".*/\1/')
        if [ -z "$VERSION" ]; then
            error "failed to determine latest version"
            exit 1
        fi
        info "latest version: $VERSION"
    fi
}

# --- Package manager detection ------------------------------------------------

detect_package_manager() {
    PKG_MANAGER=""
    PKG_FORMAT=""

    if [ "$PLATFORM" = "darwin" ]; then
        if command -v brew >/dev/null 2>&1; then
            PKG_MANAGER="brew"
            PKG_FORMAT="brew"
        fi
        return
    fi

    if [ "$PLATFORM" != "linux" ]; then
        return
    fi

    if command -v apt-get >/dev/null 2>&1 && command -v dpkg >/dev/null 2>&1; then
        PKG_MANAGER="dpkg"
        PKG_FORMAT="deb"
    elif command -v dnf >/dev/null 2>&1; then
        PKG_MANAGER="dnf"
        PKG_FORMAT="rpm"
    elif command -v yum >/dev/null 2>&1; then
        PKG_MANAGER="yum"
        PKG_FORMAT="rpm"
    elif command -v apk >/dev/null 2>&1; then
        PKG_MANAGER="apk"
        PKG_FORMAT="apk"
    fi
}

# --- Download and verify ------------------------------------------------------

build_download_url() {
    local base="https://github.com/aziontech/azion/releases/download/${VERSION}"
    local ext="$1"
    ASSET_NAME="azion_${VERSION}_${PLATFORM}_${ARCH}.${ext}"
    DOWNLOAD_URL="${base}/${ASSET_NAME}"
    CHECKSUM_URL="${base}/azion_v${VERSION}_checksum"
}

verify_checksum() {
    local file="$1"
    local checksum_file="$2"
    local expected

    expected=$(grep "${ASSET_NAME}" "$checksum_file" | awk '{print $1}')
    if [ -z "$expected" ]; then
        error "checksum not found for ${ASSET_NAME}"
        exit 1
    fi

    local actual
    if command -v sha256sum >/dev/null 2>&1; then
        actual=$(sha256sum "$file" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        actual=$(shasum -a 256 "$file" | awk '{print $1}')
    else
        warn "no sha256 tool found, skipping checksum verification"
        return 0
    fi

    if [ "$actual" != "$expected" ]; then
        error "checksum mismatch"
        error "  expected: $expected"
        error "  actual:   $actual"
        exit 1
    fi

    info "checksum verified"
}

download_and_verify() {
    local ext="$1"
    build_download_url "$ext"

    info "downloading ${ASSET_NAME}..."
    http_download "$DOWNLOAD_URL" "${TMP_DIR}/${ASSET_NAME}"

    info "downloading checksum file..."
    http_download "$CHECKSUM_URL" "${TMP_DIR}/checksum"

    verify_checksum "${TMP_DIR}/${ASSET_NAME}" "${TMP_DIR}/checksum"
}

# --- Installation -------------------------------------------------------------

install_with_package_manager() {
    local pkg_file="${TMP_DIR}/${ASSET_NAME}"

    info "installing with ${PKG_MANAGER}..."

    case "$PKG_MANAGER" in
        dpkg)
            sudo dpkg -i "$pkg_file"
            ;;
        dnf)
            sudo dnf install -y "$pkg_file"
            ;;
        yum)
            sudo yum localinstall -y "$pkg_file"
            ;;
        apk)
            sudo apk add --allow-untrusted "$pkg_file"
            ;;
    esac
}

install_with_brew() {
    info "installing with Homebrew..."
    brew install azion
}

install_binary() {
    local install_dir="${AZION_INSTALL_DIR:-$HOME/.azion/bin}"
    local zip_file="${TMP_DIR}/${ASSET_NAME}"

    if ! command -v unzip >/dev/null 2>&1; then
        error "unzip is required for binary installation"
        exit 1
    fi

    info "extracting to ${install_dir}..."
    mkdir -p "$install_dir"
    unzip -o -q "$zip_file" -d "${TMP_DIR}/extract"

    # Find the azion binary in the extracted contents
    local binary
    binary=$(find "${TMP_DIR}/extract" -name "azion" -type f | head -n 1)
    if [ -z "$binary" ]; then
        error "azion binary not found in archive"
        exit 1
    fi

    cp "$binary" "${install_dir}/azion"
    chmod +x "${install_dir}/azion"

    INSTALL_LOCATION="${install_dir}/azion"
    NEEDS_PATH_UPDATE=true

    configure_path "$install_dir"
}

# --- PATH configuration (binary fallback only) --------------------------------

configure_path() {
    local install_dir="$1"

    # Check if already in PATH
    if echo "$PATH" | tr ':' '\n' | grep -qx "$install_dir"; then
        NEEDS_PATH_UPDATE=false
        return
    fi

    local export_line="export PATH=\"${install_dir}:\$PATH\""
    UPDATED_SHELL_CONFIGS=()
    local shell_configs=()

    # Detect relevant shell config files
    if [ -f "$HOME/.bashrc" ]; then
        shell_configs+=("$HOME/.bashrc")
    fi
    if [ -f "$HOME/.zshrc" ]; then
        shell_configs+=("$HOME/.zshrc")
    fi
    if [ -f "$HOME/.bash_profile" ]; then
        shell_configs+=("$HOME/.bash_profile")
    elif [ -f "$HOME/.profile" ]; then
        shell_configs+=("$HOME/.profile")
    fi

    # If no config files exist, create .profile
    if [ ${#shell_configs[@]} -eq 0 ]; then
        shell_configs+=("$HOME/.profile")
    fi

    for config in "${shell_configs[@]}"; do
        if [ -f "$config" ] && grep -qF "$install_dir" "$config" 2>/dev/null; then
            continue
        fi
        printf '\n# Added by Azion CLI installer\n%s\n' "$export_line" >> "$config"
        UPDATED_SHELL_CONFIGS+=("$config")
    done

    if [ ${#UPDATED_SHELL_CONFIGS[@]} -gt 0 ]; then
        info "updated PATH in: ${UPDATED_SHELL_CONFIGS[*]}"
        NEEDS_PATH_UPDATE=true
    else
        NEEDS_PATH_UPDATE=false
    fi
}

# --- Post-install verification ------------------------------------------------

verify_installation() {
    # For binary installs, temporarily add to PATH for verification
    local install_dir="${AZION_INSTALL_DIR:-$HOME/.azion/bin}"
    export PATH="${install_dir}:$PATH"

    if command -v azion >/dev/null 2>&1; then
        local installed_version
        installed_version=$(azion --version 2>/dev/null || echo "unknown")
        printf '\n'
        info "${BOLD}Azion CLI installed successfully!${RESET}"
        info "version:  ${installed_version}"
        info "location: $(command -v azion)"
    else
        warn "installation completed but 'azion' was not found in PATH"
    fi

    if [ "${NEEDS_PATH_UPDATE:-false}" = true ]; then
        printf '\n'
        warn "restart your shell or run:"
        printf '  %ssource %s%s\n' "${BOLD}" "${UPDATED_SHELL_CONFIGS[0]:-"your shell config"}" "${RESET}"
    fi

    printf '\n'
    info "documentation: https://www.azion.com/en/documentation/products/azion-cli/overview/"
}

# --- Main ---------------------------------------------------------------------

main() {
    setup_colors

    printf '%s\n' "${BOLD}Azion CLI Installer${RESET}"
    printf '\n'

    detect_os
    detect_arch
    info "detected platform: ${PLATFORM}/${ARCH}"

    NEEDS_PATH_UPDATE=false
    INSTALL_LOCATION=""

    detect_package_manager

    if [ "$PKG_MANAGER" = "brew" ]; then
        info "detected package manager: Homebrew"
        install_with_brew
        INSTALL_LOCATION=$(brew --prefix)/bin/azion
    else
        detect_http_client
        resolve_version

        # Create temp directory with cleanup
        TMP_DIR=$(mktemp -d)
        trap 'rm -rf "$TMP_DIR"' EXIT

        if [ -n "$PKG_MANAGER" ]; then
            info "detected package manager: ${PKG_MANAGER}"
            download_and_verify "$PKG_FORMAT"
            install_with_package_manager
            INSTALL_LOCATION=$(command -v azion 2>/dev/null || echo "/usr/bin/azion")
        else
            info "using binary installation (zip)"
            download_and_verify "zip"
            install_binary
        fi
    fi

    verify_installation
}

main