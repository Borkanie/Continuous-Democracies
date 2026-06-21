#!/usr/bin/env bash
# If WireGuard configs are mounted, start wireproxy and set CDEP_PROXY_URL.
# Otherwise run the scraper directly (direct internet access assumed).
set -euo pipefail

WG_CONFIGS_DIR="${WG_CONFIGS_DIR:-/etc/wireguard/configs}"
PROXY_BIND="${PROXY_BIND:-127.0.0.1:25344}"

log() { echo "[$(date -u '+%Y-%m-%dT%H:%M:%SZ')] [entrypoint] $*"; }

mapfile -d '' configs < <(find "$WG_CONFIGS_DIR" -maxdepth 1 -name '*.conf' -print0 2>/dev/null)

if [ "${#configs[@]}" -gt 0 ]; then
    picked="${configs[RANDOM % ${#configs[@]}]}"
    log "Starting wireproxy with config: $(basename "$picked")"

    wpconf=$(mktemp --suffix=.conf)
    cat "$picked" > "$wpconf"
    cat >> "$wpconf" <<EOF

[Socks5]
BindAddress = ${PROXY_BIND}
EOF

    wireproxy -c "$wpconf" &
    WP_PID=$!
    trap 'kill -TERM "$WP_PID" 2>/dev/null || true' EXIT TERM INT
    sleep 2

    export CDEP_PROXY_URL="socks5h://${PROXY_BIND}"
    log "SOCKS5 proxy ready at ${PROXY_BIND}"
else
    log "No WireGuard configs found in $WG_CONFIGS_DIR — connecting directly"
fi

exec /scraper "$@"
