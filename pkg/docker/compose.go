package docker

import "fmt"

func GenerateHubCompose(hasPassword bool) string {
	return `services:
  pstar-hub:
    image: pstar/hub:latest
    container_name: pstar-hub
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    ports:
      - "51820:51820/udp"
    volumes:
      - ./pstar-data:/data
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - PSTAR_ROLE=hub
      - PSTAR_NETWORK=10.0.0.0/24
    sysctls:
      - net.ipv4.ip_forward=1
      - net.ipv4.conf.all.src_valid_mark=1
    restart: unless-stopped
`
}

func GenerateHubComposeWithMonitor(hasPassword bool) string {
	return `services:
  pstar-hub:
    image: pstar/hub:latest
    container_name: pstar-hub
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    ports:
      - "51820:51820/udp"
    volumes:
      - ./pstar-data:/data
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - PSTAR_ROLE=hub
      - PSTAR_NETWORK=10.0.0.0/24
    sysctls:
      - net.ipv4.ip_forward=1
      - net.ipv4.conf.all.src_valid_mark=1
    restart: unless-stopped

  pulse:
    image: rcourtman/pulse:latest
    container_name: pulse
    restart: unless-stopped
    ports:
      - "7655:7655"
    volumes:
      - ./pstar-data/pulse:/data
    environment:
      - TZ=UTC
      - PULSE_AUTH_USER=admin
      - PULSE_AUTH_PASS=pstar
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:7655/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
`
}

func GenerateNodeCompose(hubAddr string) string {
	return fmt.Sprintf(`services:
  pstar-node:
    image: pstar/node:latest
    container_name: pstar-node
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    volumes:
      - ./pstar-data:/data
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - PSTAR_ROLE=node
      - PSTAR_HUB=%s
    sysctls:
      - net.ipv4.ip_forward=1
      - net.ipv4.conf.all.src_valid_mark=1
    restart: unless-stopped
`, hubAddr)
}

func GenerateNodeComposeWithMonitor(hubAddr string) string {
	return fmt.Sprintf(`services:
  pstar-node:
    image: pstar/node:latest
    container_name: pstar-node
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    volumes:
      - ./pstar-data:/data
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - PSTAR_ROLE=node
      - PSTAR_HUB=%s
    sysctls:
      - net.ipv4.ip_forward=1
      - net.ipv4.conf.all.src_valid_mark=1
    restart: unless-stopped

  pulse-agent:
    image: rcourtman/pulse:latest
    container_name: pulse-agent
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./pstar-data/pulse-agent:/data
    environment:
      - PULSE_SERVER=http://%s:7655
      - TZ=UTC
    command: ["pulse-agent", "--enable-docker"]
`, hubAddr, hubAddr)
}
