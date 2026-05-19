FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        wireguard-tools \
        openresolv \
        curl \
        jq \
        iptables \
        iproute2 \
        iputils-ping \
        tinyproxy \
        ca-certificates \
        procps && \
    rm -rf /var/lib/apt/lists/* && \
    # Debian bookworm defaults iptables to the nftables backend, which is blocked
    # inside Docker containers. Switch to legacy so wg-quick can manage rules.
    update-alternatives --set iptables  /usr/sbin/iptables-legacy

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
# ── Environment variable defaults ─────────────────────────────────────────────

# Directory where WireGuard .conf files are mounted (one per server).
# The entrypoint picks one at random on each connect / reconnect.
# Mount your configs: -v /path/to/configs:/etc/wireguard/configs:ro
ENV WG_CONFIGS_DIR=/etc/wireguard/configs

# Seconds to wait for a successful WireGuard handshake after bringing the interface up
ENV CONNECT_TIMEOUT=30

# URL polled after connecting to log the new external IP. Unset to disable.
ENV IP_CHECK_URL=https://ifconfig.co/json

# Schedule automatic reconnection (picks a new random config each time).
# Either HH:MM for a daily reconnect or a duration: 30m, 6h, 1d
ENV VPN_RECONNECT=""

# Block all non-VPN traffic while the tunnel is down (kill switch)
ENV VPN_KILL_SWITCH=1

# Start tinyproxy HTTP proxy on port 3128 inside the container
ENV HTTP_PROXY=0

ENTRYPOINT ["/entrypoint.sh"]
