#!/usr/bin/env bash
set -euo pipefail

WG_IFACE="wg0"
WG_CONF="/etc/wireguard/${WG_IFACE}.conf"
RECONNECT_PID=""
SERVER_IP=""

log() { echo "[$(date -u '+%Y-%m-%dT%H:%M:%SZ')] [protonvpn] $*"; }

# ── Cleanup ───────────────────────────────────────────────────────────────────

cleanup() {
    log "Shutting down..."
    [ -n "$RECONNECT_PID" ] && kill "$RECONNECT_PID" 2>/dev/null || true
    teardown_killswitch
    bring_down_wg
    exit 0
}
trap cleanup SIGTERM SIGINT SIGQUIT

# ── Config selection ──────────────────────────────────────────────────────────

pick_config() {
    local dir="${WG_CONFIGS_DIR:-/etc/wireguard/configs}"

    [ -d "$dir" ] || { log "ERROR: config dir '${dir}' not found — mount your .conf files there"; exit 1; }

    local configs=()
    while IFS= read -r -d '' f; do
        configs+=("$f")
    done < <(find "$dir" -maxdepth 1 -name "*.conf" -print0 2>/dev/null)

    [ "${#configs[@]}" -gt 0 ] || { log "ERROR: no .conf files found in '${dir}'"; exit 1; }

    local idx=$(( RANDOM % ${#configs[@]} ))
    log "Config: $(basename "${configs[$idx]}")  (${#configs[@]} available)"
    cp "${configs[$idx]}" "$WG_CONF"
    chmod 600 "$WG_CONF"

    # Strip IPv6 from comma-separated fields — ip6tables CONNMARK is unavailable in containers
    sed -i -E \
        -e 's/,\s*[^ ,]*:[^ ,]*(\/[0-9]+)?//g' \
        -e 's/[^ ,]*:[^ ,]*(\/[0-9]+)?\s*,\s*//g' \
        "$WG_CONF"
}

# ── Kill switch (no fwmark / no CONNMARK dependency) ─────────────────────────

setup_killswitch() {
    [ "${VPN_KILL_SWITCH:-1}" != "1" ] && return
    local server_ip="$1"
    log "Kill switch ON"
    iptables -F OUTPUT 2>/dev/null || true
    iptables -A OUTPUT -o lo -j ACCEPT                  # loopback always OK
    iptables -A OUTPUT -d "$server_ip" -j ACCEPT        # must reach VPN endpoint
    iptables -A OUTPUT -o "$WG_IFACE" -j ACCEPT         # tunnel traffic OK
    iptables -A OUTPUT -j REJECT                        # block everything else
}

teardown_killswitch() {
    [ "${VPN_KILL_SWITCH:-1}" != "1" ] && return
    iptables -F OUTPUT 2>/dev/null || true
    log "Kill switch OFF"
}

# ── WireGuard — manual setup, bypasses wg-quick's iptables-restore calls ─────

bring_down_wg() {
    ip -4 route del default dev "$WG_IFACE"         2>/dev/null || true
    [ -n "$SERVER_IP" ] && \
        ip -4 route del "${SERVER_IP}/32"            2>/dev/null || true
    ip link del "$WG_IFACE"                          2>/dev/null || true
}

bring_up_wg() {
    bring_down_wg

    # Parse the wg-quick fields we need to handle ourselves
    local address dns endpoint
    address=$(awk  -F'[[:space:]]*=[[:space:]]*' '/^Address/{print $2; exit}' "$WG_CONF")
    dns=$(awk      -F'[[:space:]]*=[[:space:]]*' '/^DNS/{print $2; exit}'     "$WG_CONF")
    endpoint=$(awk -F'[[:space:]]*=[[:space:]]*' '/^Endpoint/{print $2; exit}' "$WG_CONF")
    SERVER_IP="${endpoint%:*}"

    log "Endpoint: ${endpoint}"

    # Create interface
    ip link add "$WG_IFACE" type wireguard

    # Feed wg only the fields it understands (strip wg-quick extensions)
    wg setconf "$WG_IFACE" \
        <(grep -vE '^\s*(Address|DNS|PostUp|PreDown|Table|MTU)\s*=' "$WG_CONF")

    ip -4 address add "$address" dev "$WG_IFACE"
    ip link set mtu 1420 up dev "$WG_IFACE"

    # Add specific host route to VPN server via physical gateway — prevents routing loop
    local gw dev
    gw=$(ip -4 route show default | awk 'NR==1{print $3}')
    dev=$(ip -4 route show default | awk 'NR==1{print $5}')
    if [ -n "$gw" ] && [ -n "$SERVER_IP" ]; then
        ip -4 route add "${SERVER_IP}/32" via "$gw" dev "$dev" 2>/dev/null || true
    fi

    # Replace Docker's default route with the tunnel
    ip -4 route del default 2>/dev/null || true
    ip -4 route add default dev "$WG_IFACE"

    # DNS — take first entry if comma-separated
    if [ -n "$dns" ]; then
        echo "nameserver ${dns%%,*}" > /etc/resolv.conf
        log "DNS → ${dns%%,*}"
    fi
}

# ── Handshake wait ────────────────────────────────────────────────────────────

wait_for_handshake() {
    local timeout="${CONNECT_TIMEOUT:-30}" elapsed=0
    log "Waiting for handshake (timeout: ${timeout}s)..."
    while [ "$elapsed" -lt "$timeout" ]; do
        local ts
        ts=$(wg show "$WG_IFACE" latest-handshakes 2>/dev/null | awk '{print $2}')
        if [ -n "$ts" ] && [ "$ts" != "0" ]; then
            log "Handshake OK (${elapsed}s)"
            return 0
        fi
        sleep 1
        elapsed=$(( elapsed + 1 ))
    done
    log "WARNING: No handshake within ${timeout}s"
    return 1
}

# ── Connect / reconnect ───────────────────────────────────────────────────────

connect() {
    teardown_killswitch
    pick_config
    bring_up_wg
    setup_killswitch "$SERVER_IP"
    wait_for_handshake || true

    if [ -n "${IP_CHECK_URL:-}" ]; then
        local ip http_code curl_exit
        ip=$(curl -s --max-time 10 -w "\n%{http_code}" "${IP_CHECK_URL}" 2>/dev/null); curl_exit=$?
        http_code=$(echo "$ip" | tail -n1)
        ip=$(echo "$ip" | head -n -1 | jq -r '.ip // .ip_addr // "unknown"' 2>/dev/null || echo "unknown")
        log "External IP: ${ip}  (curl_exit=${curl_exit}, http=${http_code})"
    fi
}

reconnect() { log "Reconnecting..."; connect; }

handle_reconnect_signal() { reconnect; }
trap handle_reconnect_signal USR1

# ── Reconnect scheduler ───────────────────────────────────────────────────────

parse_duration_secs() {
    local val="$1"
    if   [[ "$val" =~ ^([0-9]+)s$ ]]; then echo "${BASH_REMATCH[1]}"
    elif [[ "$val" =~ ^([0-9]+)m$ ]]; then echo $(( BASH_REMATCH[1] * 60 ))
    elif [[ "$val" =~ ^([0-9]+)h$ ]]; then echo $(( BASH_REMATCH[1] * 3600 ))
    elif [[ "$val" =~ ^([0-9]+)d$ ]]; then echo $(( BASH_REMATCH[1] * 86400 ))
    else echo "0"
    fi
}

secs_until_daily_time() {
    local target="$1" now target_epoch
    now=$(date +%s)
    target_epoch=$(date -d "today ${target}" +%s)
    [ "$target_epoch" -le "$now" ] && target_epoch=$(date -d "tomorrow ${target}" +%s)
    echo $(( target_epoch - now ))
}

start_reconnect_scheduler() {
    local schedule="${VPN_RECONNECT}" parent_pid="$$"
    (
        while true; do
            local wait_secs
            if [[ "$schedule" =~ ^[0-9]{1,2}:[0-9]{2}$ ]]; then
                wait_secs=$(secs_until_daily_time "$schedule")
                log "Next reconnect at ${schedule} (in ${wait_secs}s)"
            else
                wait_secs=$(parse_duration_secs "$schedule")
                log "Next reconnect in ${wait_secs}s"
            fi
            [ "$wait_secs" -le 0 ] && wait_secs=3600
            sleep "$wait_secs"
            kill -USR1 "$parent_pid" 2>/dev/null || true
        done
    ) &
    RECONNECT_PID=$!
}

# ── HTTP proxy ────────────────────────────────────────────────────────────────

start_proxy() {
    log "Starting tinyproxy on port 3128..."
    cat > /tmp/tinyproxy.conf <<'TPCONF'
Port 3128
Listen 0.0.0.0
Timeout 600
Allow 0.0.0.0/0
DisableViaHeader Yes
LogFile "/dev/null"
TPCONF
    tinyproxy -c /tmp/tinyproxy.conf
}

# ── Main ──────────────────────────────────────────────────────────────────────

[ "${HTTP_PROXY:-0}" = "1" ] && start_proxy

connect

[ -n "${VPN_RECONNECT:-}" ] && start_reconnect_scheduler

log "VPN up. Monitoring..."
while true; do
    sleep 30 &
    wait $! 2>/dev/null || true   # wakes early on USR1

    if ! ip link show "$WG_IFACE" &>/dev/null 2>&1; then
        log "Interface ${WG_IFACE} lost — reconnecting..."
        reconnect
        continue
    fi

    local_ts=$(wg show "$WG_IFACE" latest-handshakes 2>/dev/null | awk '{print $2}')
    if [ -n "$local_ts" ] && [ "$local_ts" != "0" ]; then
        age=$(( $(date +%s) - local_ts ))
        if [ "$age" -gt 180 ]; then
            log "Handshake stale (${age}s) — reconnecting..."
            reconnect
        fi
    fi
done
