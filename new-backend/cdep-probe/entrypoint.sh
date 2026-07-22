#!/usr/bin/env bash
set -euo pipefail

PROXY_BIND="${PROXY_BIND:-127.0.0.1:25344}"
VPN_TIMEOUT="${VPN_TIMEOUT:-45}"
CDEP_TARGET="${CDEP_TARGET:-https://www.cdep.ro}"

log() { echo "[$(date -u '+%Y-%m-%dT%H:%M:%SZ')] $*"; }

if [ -z "${WG_PRIVATE_KEY:-}" ]; then
    log "ERROR: WG_PRIVATE_KEY env var is required"
    exit 1
fi

# Render the WireGuard config, appending the SOCKS5 listener for wireproxy.
WPCONF=$(mktemp --suffix=.conf)
envsubst < /etc/wireguard/wg.conf.tmpl > "$WPCONF"
printf '\n[Socks5]\nBindAddress = %s\n' "$PROXY_BIND" >> "$WPCONF"

cleanup() { kill -TERM "$WP_PID" 2>/dev/null || true; rm -f "$WPCONF"; }
trap cleanup EXIT TERM INT

log "Starting wireproxy (RO#103 endpoint)..."
wireproxy -c "$WPCONF" &
WP_PID=$!

log "Waiting for tunnel (timeout=${VPN_TIMEOUT}s)..."
deadline=$(( $(date +%s) + VPN_TIMEOUT ))
until curl -s --max-time 5 --proxy "socks5h://${PROXY_BIND}" https://1.1.1.1 -o /dev/null 2>/dev/null; do
    if [ "$(date +%s)" -ge "$deadline" ]; then
        log "FAIL: tunnel did not come up within ${VPN_TIMEOUT}s"
        exit 1
    fi
    sleep 2
done
log "Tunnel is up."

log "Probing ${CDEP_TARGET} through SOCKS5 proxy..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    --max-time 30 \
    --proxy "socks5h://${PROXY_BIND}" \
    "${CDEP_TARGET}")

log "HTTP response code: ${HTTP_CODE}"

if [[ "$HTTP_CODE" =~ ^[23] ]]; then
    log "SUCCESS: reached ${CDEP_TARGET} via ProtonVPN RO#103"
    exit 0
else
    log "FAIL: unexpected HTTP code ${HTTP_CODE} from ${CDEP_TARGET}"
    exit 1
fi
