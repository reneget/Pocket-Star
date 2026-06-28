FROM alpine:3.20

RUN apk add --no-cache wireguard-tools iptables iproute2 bash

COPY pstar-hub /usr/local/bin/pstar
COPY logo.svg /usr/local/share/pocket-star/logo.svg

VOLUME ["/data"]
EXPOSE 51820/udp

LABEL org.opencontainers.image.title="Pocket Star Hub" \
      org.opencontainers.image.description="Site-to-site VPN hub (AmneziaWG)" \
      org.opencontainers.image.version="0.1.0" \
      org.opencontainers.image.url="https://github.com/reneget/Pocket-Star"

ENTRYPOINT ["pstar"]
CMD ["hub", "status"]
