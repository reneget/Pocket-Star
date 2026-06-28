FROM alpine:3.20

RUN apk add --no-cache wireguard-tools iptables iproute2 bash

COPY pstar-node /usr/local/bin/pstar

VOLUME ["/data"]

ENTRYPOINT ["pstar"]
CMD ["node", "status"]
