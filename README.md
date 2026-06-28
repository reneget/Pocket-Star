# Pocket Star

> Turn your old computer into a personal server without a public IP. Your home, reachable from anywhere.

Pocket Star solves the problem when you have an old PC you want to use as a server (file storage, local AI inference, dev environment), but your ISP doesn't give you a public IP. Renting a VPS with the same specs as your old PC would cost a fortune.

Instead, you rent the **cheapest VPS** you can find (~$3/mo) — it acts as a hub with a public IP. Your server and your laptop connect to it via AmneziaWG, forming a virtual local network. Just like that, your server is accessible from anywhere.

```
VPS (Hub)         Old PC (Server)     Laptop (Client)
10.0.0.1 ──────── 10.0.0.2            10.0.0.3
    │                   │                   │
    └── AmneziaWG ──────┴───────────────────┘
    Docker Compose      Docker Compose       Native WG GUI
                                             or pstar CLI
```

## Quick Start

### 1. On VPS (hub)

```bash
curl -sSL https://github.com/reneget/Pocket-Star/raw/main/scripts/quick-start.sh | bash

pstar hub init --docker --pass
# → Creates: hub.conf, docker-compose.yml, encrypted client configs
# → Set a master password (protects node configs)

docker compose -f pstar-data/docker-compose.yml up -d
```

### 2. On old PC (server)

```bash
pstar node join --hub <YOUR_VPS_IP>:51820 --docker --name server
# → Enter the master password you set on the hub

docker compose -f pstar-data/docker-compose.yml up -d
```

**No Docker?** Copy `pstar-data/clients/server.conf` from the hub and import it into WireGuard / AmneziaWG GUI.

### 3. On laptop (client)

```bash
pstar node join --hub <YOUR_VPS_IP>:51820 --name work-laptop
# Enter master password → connected.
```

**Or via GUI:** copy `work-laptop.conf` → import into WireGuard/Amnezia → hit Connect.

### 4. Verify

```bash
ping 10.0.0.2               # Ping the server
ssh user@10.0.0.2           # SSH into it
scp file user@10.0.0.2:~/  # Copy files directly
```

## Commands

| Command | What it does |
|---------|--------------|
| `pstar hub init --docker --pass` | Initialize hub on VPS |
| `pstar hub status` | List connected peers |
| `pstar node join --hub X` | Connect this machine as a node |
| `pstar node status` | VPN connection status |
| `pstar module list` | List available modules |
| `pstar decrypt <file>` | Decrypt a config with master password |
| `pstar doctor` | Run system diagnostics |
| `pstar version` | Print version |

## Master Password

Node configs are encrypted with **AES-256-GCM** using a master password. This protects private keys during transfer over SCP or any other channel.

```bash
# On hub: configs are saved as *.conf.enc
# On node: run `pstar decrypt server.conf.enc` and enter the password
# Or just: `pstar node join` — it prompts for the password automatically
```

## Module System

Everything is built as pluggable modules. You can write your own by implementing one interface:

```go
type Module interface {
    Name() string
    Priority() int
    Dependencies() []string
    Init(*Context) error
    Start() error
    Stop() error
    Status() Status
}
```

### Planned modules

- **core** — VPN + routing (MVP)
- **sftp** — file sharing over SFTP
- **port-fwd** — forward ports through the hub
- **dashboard** — web UI

## Architecture

```
pocket-star/
├── cmd/pstar/           # Entry point
├── pkg/
│   ├── cli/             # CLI commands (hub, node, module, doctor)
│   ├── config/          # AES-256-GCM encryption
│   ├── vpn/             # AmneziaWG key + config generation
│   ├── module/          # Module interface + registry
│   └── docker/          # docker-compose generation
├── modules/
│   └── core/            # Core module (VPN)
├── images/
│   ├── hub/Dockerfile
│   └── node/Dockerfile
└── scripts/
    ├── quick-start.sh   # curl-to-bash installer
    └── bootstrap.sh     # Fully automated setup
```

## Development

```bash
# Build
go build -o pstar ./cmd/pstar

# Test
go test ./...

# Cross-compile
GOOS=linux GOARCH=arm64 go build -o pstar-arm64 ./cmd/pstar
```

## License

MIT
