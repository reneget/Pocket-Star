# ⭐ Pocket Star

<div align="center">

<img src="pkg/icon/logo.svg" width="120" alt="Pocket Star Logo"/>

[![Readme in Russian](https://img.shields.io/badge/README-%D0%A0%D1%83%D1%81%D1%81%D0%BA%D0%B8%D0%B9-blue.svg)](README.ru.md)
![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)
![Go](https://img.shields.io/badge/go-1.25+-00ADD8.svg)
![AmneziaWG](https://img.shields.io/badge/vpn-AmneziaWG-green.svg)
![Docker](https://img.shields.io/badge/docker-ready-blue.svg)
![License](https://img.shields.io/badge/license-MIT-yellow.svg)

**Turn your old computer into a personal server without a public IP. Your home, reachable from anywhere.**

[📖 Documentation](#documentation) • [🚀 Quick Start](#quick-start) • [🏗️ Architecture](#architecture) • [🛠️ Technology Stack](#technology-stack) • [📦 Modules](#modules)

</div>

---

## 📋 Table of Contents

- [About the Project](#about-the-project)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Modules](#modules)
- [Master Password](#master-password)
- [Development](#development)
- [Documentation](#documentation)

---

## 🎯 About the Project

Pocket Star solves the problem when you have an old PC you want to use as a server (file storage, local AI inference, dev environment), but your ISP doesn't give you a public IP. Renting a VPS with the same specs as your old PC would cost a fortune.

Instead, you rent the **cheapest VPS** (~$3/mo) — it acts as a hub with a public IP. Your server and your laptop connect to it via AmneziaWG, forming a virtual local network. Just like that, your server is accessible from anywhere.

### Key Features

✨ **For users:**
- 🔌 Site-to-site VPN in minutes — connect your old PC and laptop into one network
- 🔒 AES-256-GCM encrypted configs with master password
- 🐳 Docker Compose support for hub and server nodes
- 🖥️ Native WireGuard GUI support — import the config and connect
- 🔧 Modular design — extend with plugins

---

## 🏗️ Architecture

### Network topology

```mermaid
graph TB
    subgraph "Internet"
        VPS[VPS Hub<br/>Public IP<br/>10.0.0.1]
    end
    
    subgraph "Home Network"
        SRV[Old PC / Server<br/>10.0.0.2]
    end
    
    subgraph "Anywhere"
        LAP[Laptop / Client<br/>10.0.0.3]
    end
    
    VPS <-->|AmneziaWG<br/>Tunnel| SRV
    VPS <-->|AmneziaWG<br/>Tunnel| LAP
    SRV <-.->|SSH / Ping / SCP<br/>via VPN| LAP
    
    style VPS fill:#4ecdc4,color:#fff
    style SRV fill:#ff6b6b,color:#fff
    style LAP fill:#45b7d1,color:#fff
```

### Project structure

```mermaid
graph LR
    subgraph "CLI"
        PSTAR[pstar<br/>Entry point]
    end
    
    subgraph "Core Packages"
        CLI[CLI Commands]
        VPN[VPN Engine<br/>AmneziaWG]
        CRYPT[Crypto<br/>AES-256-GCM]
        MOD[Module System<br/>+ Manager]
        DOCK[Docker<br/>Compose]
        PULSE[Pulse<br/>Client]
    end
    
    subgraph "Modules"
        CORE[Core Module]
        TUI[TUI Module<br/>Bubbletea]
        MON[Monitoring<br/>Checker + Alerts]
        SFTP[SFTP Module<br/>Planned]
        FWD[Port Forward<br/>Planned]
        MEDIA[Media Storage<br/>Planned]
        CLOUD[Cloud Sync<br/>Planned]
    end
    
    subgraph "Distribution"
        HUB_IMG[Hub Docker Image]
        NODE_IMG[Node Docker Image]
        SCRIPTS[Setup Scripts]
    end
    
    PSTAR --> CLI
    PSTAR --> VPN
    PSTAR --> CRYPT
    PSTAR --> MOD
    PSTAR --> DOCK
    MOD --> CORE
    MOD --> TUI
    MOD --> MON
    MOD -.-> SFTP
    MOD -.-> FWD
    MOD -.-> MEDIA
    MOD -.-> CLOUD
    DOCK --> HUB_IMG
    DOCK --> NODE_IMG
    CLI --> SCRIPTS
    
    style PSTAR fill:#4ecdc4,color:#fff
    style CLI fill:#96ceb4,color:#fff
    style VPN fill:#ff9ff3,color:#fff
    style CRYPT fill:#feca57,color:#000
    style MOD fill:#45b7d1,color:#fff
    style PULSE fill:#96ceb4,color:#fff
    style CORE fill:#ff6b6b,color:#fff
    style TUI fill:#ff9ff3,color:#fff
    style MON fill:#4ecdc4,color:#fff
    style SFTP fill:#a0a0a0,color:#fff,stroke-dasharray: 6 4,stroke-width:2px
    style FWD fill:#a0a0a0,color:#fff,stroke-dasharray: 6 4,stroke-width:2px
    style MEDIA fill:#a0a0a0,color:#fff,stroke-dasharray: 6 4,stroke-width:2px
    style CLOUD fill:#a0a0a0,color:#fff,stroke-dasharray: 6 4,stroke-width:2px
```

### Connection sequence

```mermaid
sequenceDiagram
    participant VPS as 🖥️ VPS Hub
    participant SRV as 🖧 Old PC Server
    participant LAP as 💻 Laptop Client
    
    Note over VPS,LAP: Step 1: Hub initialization
    VPS->>VPS: pstar hub init --docker --pass
    VPS->>VPS: Generate keys, configs, docker-compose
    VPS->>VPS: docker compose up -d
    
    Note over VPS,LAP: Step 2: Server connects
    SRV->>VPS: pstar node join --hub IP:51820
    VPS-->>SRV: Assign IP 10.0.0.2
    SRV->>SRV: docker compose up -d
    
    Note over VPS,LAP: Step 3: Client connects
    LAP->>VPS: pstar node join --hub IP:51820
    VPS-->>LAP: Assign IP 10.0.0.3
    
    Note over VPS,LAP: Step 4: Verify connectivity
    LAP->>SRV: ping 10.0.0.2
    SRV-->>LAP: ✅ pong
    LAP->>SRV: ssh user@10.0.0.2
    SRV-->>LAP: ✅ Connected
```

---

## 🛠️ Technology Stack

### Core
- **Go 1.25+** — single binary, cross-platform compilation
- **AmneziaWG** — VPN with obfuscation, WireGuard-compatible
- **cobra** — CLI framework

### Security
- **Curve25519** — cryptographic key exchange
- **AES-256-GCM** — config encryption with master password

### Infrastructure
- **Docker & Docker Compose** — containerization for hub and server
- **Alpine Linux** — base image for Docker images

### VPN Protocols
- **AmneziaWG** — default VPN with DPI obfuscation
- **WireGuard** — fully compatible, import .conf into any WG client

---

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose (for hub and server)
- A cheap VPS with a public IP (any provider)
- Linux / macOS on the client machine

### Installation

```bash
curl -sSL https://github.com/reneget/Pocket-Star/raw/main/scripts/quick-start.sh | bash
```

### 1. Initialize Hub on VPS

```bash
pstar hub init --docker --pass
# → Creates: hub.conf, docker-compose.yml, encrypted client configs
# → Set a master password (protects node configs)

docker compose -f pstar-data/docker-compose.yml up -d
```

### 2. Connect Server (Old PC)

```bash
pstar node join --hub <YOUR_VPS_IP>:51820 --docker --name server
# → Enter the master password you set on the hub

docker compose -f pstar-data/docker-compose.yml up -d
```

**No Docker?** Copy `pstar-data/clients/server.conf` from the hub and import it into WireGuard / AmneziaWG GUI.

### 3. Connect Client (Laptop)

```bash
pstar node join --hub <YOUR_VPS_IP>:51820 --name work-laptop
# Enter master password → connected.
```

**Or via GUI:** copy `work-laptop.conf` → import into WireGuard/Amnezia → Connect.

### 4. Verify

```bash
ping 10.0.0.2               # Ping the server
ssh user@10.0.0.2           # SSH into it
scp file user@10.0.0.2:~/  # Copy files directly
```

---

## ⌨️ Commands

| Command | Description |
|---------|-------------|
| `pstar hub init --docker --pass` | Initialize hub on VPS |
| `pstar hub status` | List connected peers |
| `pstar node join --hub X [--docker]` | Connect this machine as a node |
| `pstar node status` | VPN connection status |
| `pstar module list` | List available modules |
| `pstar module enable <name>` | Enable a module |
| `pstar module disable <name>` | Disable a module |
| `pstar decrypt <file>` | Decrypt a config with master password |
| `pstar doctor` | Run system diagnostics |
| `pstar version` | Print version |

---

## 📦 Modules

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

### ✅ Available

| Module | Description | Status |
|--------|-------------|--------|
| **core** | VPN + routing | Done |
| **tui** | Terminal UI (bubbletea, mouse, tabs, commands) | Done |
| **monitor** | System/Docker/peer health checks, alerts, ring buffer store | Done |
| **pulse** | Pulse REST API client for remote metrics | Done |

### 🔲 Planned

| Module | Description |
|--------|-------------|
| **sftp** | File sharing over SFTP |
| **port-fwd** | Forward ports through the hub |
| **dashboard** | Web UI |
| **media** | Media storage server (Jellyfin/Immich) |
| **cloud** | Cloud storage & file sync (Nextcloud) |

---

## 🔒 Master Password

Node configs are encrypted with **AES-256-GCM** using a master password. This protects private keys during transfer over SCP or any other channel.

```bash
# On hub: configs are saved as *.conf.enc
# On node: run `pstar decrypt server.conf.enc` and enter the password
# Or just: `pstar node join` — it prompts for the password automatically
```

---

## 🔧 Development

### Build

```bash
go build -o pstar ./cmd/pstar
```

### Test

```bash
go test ./...
```

### Cross-compile

```bash
GOOS=linux GOARCH=arm64 go build -o pstar-arm64 ./cmd/pstar
GOOS=linux GOARCH=amd64 go build -o pstar-linux ./cmd/pstar
GOOS=darwin GOARCH=amd64 go build -o pstar-macos ./cmd/pstar
```

### Project structure

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

---

## 📊 Project Stats

```mermaid
pie title Code Distribution by Package
    "VPN Engine (AmneziaWG)" : 30
    "CLI Commands" : 25
    "Crypto (AES-256-GCM)" : 15
    "Module System" : 15
    "Docker / Scripts" : 10
    "Tests" : 5
```

---

## 📖 Documentation

- [README.ru.md](./README.ru.md) — Russian version of this document

---

## 🤝 Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

---

## 👥 Authors

- **[reneget](https://github.com/reneget)**

---

## ⭐ Star History

<div align="center">

[![Star History Chart](https://api.star-history.com/svg?repos=reneget/Pocket-Star&type=Date)](https://star-history.com/#reneget/Pocket-Star&Date)

</div>

---

<div align="center">

**Made with ❤️ to turn old PCs into servers**

⭐ If this project helped you, give it a star!

</div>
