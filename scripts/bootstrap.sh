#!/usr/bin/env bash
set -euo pipefail

# bootstrap.sh — Full automated setup for Pocket Star

PSTAR_DATA="${PSTAR_DATA:-./pstar-data}"

info()  { echo -e "\033[1;34m[*]\033[0m $*"; }
ok()    { echo -e "\033[1;32m[+]\033[0m $*"; }
err()   { echo -e "\033[1;31m[-]\033[0m $*"; }

check_prereqs() {
    local missing=0
    for cmd in curl wget ip modprobe; do
        if ! command -v "$cmd" &>/dev/null; then
            err "$cmd not found (try: apt install -y $cmd)"
            missing=1
        fi
    done
    return "$missing"
}

install_pstar() {
    if command -v pstar &>/dev/null; then
        ok "pstar already installed"
        return
    fi
    info "Installing pstar..."
    bash <(curl -sSL https://github.com/anomalyco/pocket-star/raw/main/scripts/quick-start.sh)
    ok "pstar installed"
}

setup_hub() {
    info "Setting up as HUB..."
    pstar hub init --docker --pass --data-dir "$PSTAR_DATA"
    cd "$PSTAR_DATA" && docker compose up -d
    ok "Hub is running!"
    info "Your public IP: $(curl -s ifconfig.me)"
}

setup_node() {
    local hub_addr="${1:-}"
    if [ -z "$hub_addr" ]; then
        err "Usage: bootstrap.sh node <HUB_IP>"
        exit 1
    fi
    info "Connecting to hub at $hub_addr..."
    pstar node join --hub "$hub_addr" --docker --data-dir "$PSTAR_DATA"
    cd "$PSTAR_DATA" && docker compose up -d
    ok "Node connected!"
}

case "${1:-}" in
    hub)
        check_prereqs
        install_pstar
        setup_hub
        ;;
    node)
        check_prereqs
        install_pstar
        setup_node "${2:-}"
        ;;
    *)
        echo "Usage: $0 {hub|node <HUB_IP>}"
        echo ""
        echo "  hub          — Set up this machine as the central VPS hub"
        echo "  node <IP>    — Connect this machine as a node to the hub at IP"
        exit 1
        ;;
esac
