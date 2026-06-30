#!/usr/bin/env bash
# If WireGuard configs are mounted, start wireproxy and set CDEP_PROXY_URL.
# Tries each config in random order until one establishes a working tunnel.
# Falls back to direct internet if no configs are found.
set -euo pipefail

WG_CONFIGS_DIR="${WG_CONFIGS_DIR:-/etc/wireguard/configs}"
PROXY_BIND="${PROXY_BIND:-127.0.0.1:25344}"
VPN_TIMEOUT="${VPN_TIMEOUT:-45}"
# Probe a reliable host to confirm the tunnel is up, not the scrape target itself.
VPN_PROBE_URL="${VPN_PROBE_URL:-https://1.1.1.1}"

log() { echo "[$(date -u '+%Y-%m-%dT%H:%M:%SZ')] [entrypoint] $*"; }

WP_PID=""
cleanup() { [ -n "$WP_PID" ] && kill -TERM "$WP_PID" 2>/dev/null || true; }
trap cleanup EXIT TERM INT

try_config() {
    local config_file="$1"
    log "Trying config: $(basename "$config_file")"

    local wpconf
    wpconf=$(mktemp --suffix=.conf)
    cat "$config_file" > "$wpconf"
    printf '\n[Socks5]\nBindAddress = %s\n' "$PROXY_BIND" >> "$wpconf"

    wireproxy -c "$wpconf" &
    WP_PID=$!

    local deadline=$(( $(date +%s) + VPN_TIMEOUT ))
    while ! curl -s --max-time 5 --proxy "socks5h://${PROXY_BIND}" "$VPN_PROBE_URL" -o /dev/null 2>/dev/null; do
        if [ "$(date +%s)" -ge "$deadline" ]; then
            log "Config $(basename "$config_file") did not connect within ${VPN_TIMEOUT}s"
            kill -TERM "$WP_PID" 2>/dev/null || true
            WP_PID=""
            rm -f "$wpconf"
            sleep 1  # let the port free before trying next config
            return 1
        fi
        sleep 2
    done

    rm -f "$wpconf"
    return 0
}

mapfile -d '' configs < <(find "$WG_CONFIGS_DIR" -maxdepth 1 -name '*.conf' -print0 2>/dev/null)

if [ "${#configs[@]}" -gt 0 ]; then
    # Start at a random index, then try all configs in order from there.
    start=$(( RANDOM % ${#configs[@]} ))
    connected=false

    for (( offset = 0; offset < ${#configs[@]}; offset++ )); do
        index=$(( (start + offset) % ${#configs[@]} ))
        if try_config "${configs[$index]}"; then
            connected=true
            break
        fi
    done

    if [ "$connected" = false ]; then
        log "ERROR: all ${#configs[@]} WireGuard config(s) failed — check ProtonVPN server availability"
        exit 1
    fi

    export CDEP_PROXY_URL="socks5h://${PROXY_BIND}"
    log "Tunnel verified, SOCKS5 proxy ready at ${PROXY_BIND}"
else
    log "No WireGuard configs found in $WG_CONFIGS_DIR — connecting directly"
fi

exec /scraper "$@"
