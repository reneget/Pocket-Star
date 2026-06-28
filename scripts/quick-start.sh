#!/usr/bin/env bash
set -euo pipefail

MPS_VERSION="v0.1.0"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

URL="https://github.com/anomalyco/my-pretty-star/releases/download/${MPS_VERSION}/mps-${OS}-${ARCH}.tar.gz"

echo "Downloading mps ${MPS_VERSION} for ${OS}/${ARCH}..."
if command -v curl &>/dev/null; then
    curl -sSL "$URL" | tar xz -C /usr/local/bin/ mps
elif command -v wget &>/dev/null; then
    wget -qO- "$URL" | tar xz -C /usr/local/bin/ mps
else
    echo "Need curl or wget"
    exit 1
fi

chmod +x /usr/local/bin/mps
echo "✔ mps installed. Run 'mps doctor' to verify."
