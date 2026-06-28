FROM alpine:3.20

RUN apk add --no-cache wireguard-tools iptables iproute2 bash

COPY pstar-hub /usr/local/bin/pstar

VOLUME ["/data"]
EXPOSE 51820/udp

ENTRYPOINT ["pstar"]
CMD ["hub", "status"]
