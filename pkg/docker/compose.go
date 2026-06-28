package docker

import "fmt"

func GenerateHubCompose(hasPassword bool) string {
	return `services:
  mps-hub:
    image: mps/hub:latest
    container_name: mps-hub
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    ports:
      - "51820:51820/udp"
    volumes:
      - ./mps-data:/data
    environment:
      - MPS_ROLE=hub
      - MPS_NETWORK=10.0.0.0/24
    sysctls:
      - net.ipv4.ip_forward=1
      - net.ipv4.conf.all.src_valid_mark=1
    restart: unless-stopped
`
}

func GenerateNodeCompose(hubAddr string) string {
	return fmt.Sprintf(`services:
  mps-node:
    image: mps/node:latest
    container_name: mps-node
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    volumes:
      - ./mps-data:/data
    environment:
      - MPS_ROLE=node
      - MPS_HUB=%s
    sysctls:
      - net.ipv4.ip_forward=1
      - net.ipv4.conf.all.src_valid_mark=1
    restart: unless-stopped
`, hubAddr)
}
