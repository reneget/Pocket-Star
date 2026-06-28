FROM alpine:3.20

RUN apk add --no-cache wireguard-tools iptables iproute2 bash

COPY mps-hub /usr/local/bin/mps

VOLUME ["/data"]
EXPOSE 51820/udp

ENTRYPOINT ["mps"]
CMD ["hub", "status"]
