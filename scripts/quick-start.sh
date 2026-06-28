#!/usr/bin/env bash
set -euo pipefail

PSTAR_VERSION="v0.1.0"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

URL="https://github.com/reneget/Pocket-Star/releases/download/${PSTAR_VERSION}/pstar-${OS}-${ARCH}.tar.gz"

echo "Downloading pstar ${PSTAR_VERSION} for ${OS}/${ARCH}..."
if command -v curl &>/dev/null; then
    curl -sSL "$URL" | tar xz -C /usr/local/bin/ pstar
elif command -v wget &>/dev/null; then
    wget -qO- "$URL" | tar xz -C /usr/local/bin/ pstar
else
    echo "Need curl or wget"
    exit 1
fi

chmod +x /usr/local/bin/pstar
echo "✔ pstar installed. Run 'pstar doctor' to verify."
