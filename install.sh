#!/bin/bash

set -euo pipefail

function info() {
    echo "[INFO] $1" >&2
}

function error() {
    echo "[ERROR] $1" >&2
}

function get_os() {
    if [[ "$OSTYPE" == "linux-gnu" ]]; then
        echo "Linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "Darwin"
    fi
}

function get_arch() {
    uname -m
}

function get_latest_release_url() {
    local os="$1"
    local arch="$2"
    local url="https://api.github.com/repos/vitorfhc/queryxss/releases/latest"
    local os_arch="$1_$2"
    local ret=$(curl -s "$url" | grep "browser_download_url.*$os_arch.*"| cut -d '"' -f 4)
    if [[ -z "$ret" ]]; then
        echo "No release found for $os_arch" >&2
        exit 1
    fi
    echo "$ret"
}

function download() {
    local url="$1"
    local file="$2"
    curl -Lso "$file" "$url"
}

function has_curl() {
    command -v curl > /dev/null
}

if ! has_curl; then
    error "curl is required to run this script"
    exit 1
fi

info "Finding latest release..."
LTS_URL=`get_latest_release_url $(get_os) $(get_arch)`

info "Downloading..."
TMP=`mktemp -d`
TMP_TAR="$TMP/queryxss.tar.gz"
download "$LTS_URL" "$TMP_TAR"

info "Extracting..."
tar -xzf "$TMP_TAR" -C "$TMP"
info "Moving to /usr/local/bin"
mv "$TMP/queryxss" /usr/local/bin 2> /dev/null || {
    error "Failed to move to /usr/local/bin"
    error "Please run this script as root or with sudo"
    exit 1
}
