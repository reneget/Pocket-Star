FROM alpine:3.20

RUN apk add --no-cache wireguard-tools iptables iproute2 bash

COPY mps-node /usr/local/bin/mps

VOLUME ["/data"]

ENTRYPOINT ["mps"]
CMD ["node", "status"]
