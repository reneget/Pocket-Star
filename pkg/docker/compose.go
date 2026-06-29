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

  uptime-kuma:
    image: louislam/uptime-kuma:latest
    container_name: uptime-kuma
    ports:
      - "3001:3001"
    volumes:
      - ./pstar-data/uptime:/app/data
    restart: unless-stopped
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
