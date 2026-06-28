FROM alpine:3.20

RUN apk add --no-cache wireguard-tools iptables iproute2 bash

COPY pstar-node /usr/local/bin/pstar
COPY logo.svg /usr/local/share/pocket-star/logo.svg

VOLUME ["/data"]

LABEL org.opencontainers.image.title="Pocket Star Node" \
      org.opencontainers.image.description="Site-to-site VPN node (AmneziaWG)" \
      org.opencontainers.image.version="0.1.0" \
      org.opencontainers.image.url="https://github.com/reneget/Pocket-Star"

ENTRYPOINT ["pstar"]
CMD ["node", "status"]
