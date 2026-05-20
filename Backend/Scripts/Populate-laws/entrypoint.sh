#!/usr/bin/env bash
# Picks a random WireGuard config, starts wireproxy as a userspace tunnel
# exposing a SOCKS5 listener on PROXY_BIND, then exec's the Python orchestrator.
# Only cdep.ro traffic is routed through the SOCKS proxy (see utils._proxies_for).
set -euo pipefail

WG_CONFIGS_DIR="${WG_CONFIGS_DIR:-/etc/wireguard/configs}"
PROXY_BIND="${PROXY_BIND:-127.0.0.1:25344}"

log() { echo "[$(date -u '+%Y-%m-%dT%H:%M:%SZ')] [entrypoint] $*"; }

mapfile -d '' configs < <(find "$WG_CONFIGS_DIR" -maxdepth 1 -name '*.conf' -print0)
if [ "${#configs[@]}" -eq 0 ]; then
    log "ERROR: no .conf files found in $WG_CONFIGS_DIR — mount your wg configs there"
    exit 1
fi
picked="${configs[RANDOM % ${#configs[@]}]}"
log "Using WG config: $(basename "$picked") (${#configs[@]} available)"

wpconf=$(mktemp --suffix=.conf)
cat "$picked" > "$wpconf"
cat >> "$wpconf" <<EOF

[Socks5]
BindAddress = ${PROXY_BIND}
EOF

wireproxy -c "$wpconf" &
WP_PID=$!
trap 'kill -TERM "$WP_PID" 2>/dev/null || true' EXIT TERM INT

# Brief pause so the SOCKS listener is up before Python starts hitting it
sleep 2

export CDEP_PROXY_URL="socks5h://${PROXY_BIND}"
exec python -u orchestrator.py
