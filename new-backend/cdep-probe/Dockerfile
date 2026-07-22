# ── Stage 1: wireproxy binary ─────────────────────────────────────────────────
FROM debian:bookworm-slim AS wireproxy-bin
ARG TARGETARCH
ARG WIREPROXY_VERSION=v1.1.2
RUN apt-get update \
    && apt-get install -y --no-install-recommends curl ca-certificates \
    && curl -fsSL \
        "https://github.com/pufferffish/wireproxy/releases/download/${WIREPROXY_VERSION}/wireproxy_linux_${TARGETARCH}.tar.gz" \
        -o /tmp/wp.tar.gz \
    && tar xzf /tmp/wp.tar.gz -C /tmp \
    && chmod +x /tmp/wireproxy

# ── Stage 2: runtime ──────────────────────────────────────────────────────────
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    gettext-base \
    && rm -rf /var/lib/apt/lists/*

COPY --from=wireproxy-bin /tmp/wireproxy /usr/local/bin/wireproxy
COPY wg.conf.tmpl /etc/wireguard/wg.conf.tmpl
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# WG_PRIVATE_KEY must be supplied at runtime via -e
ENV WG_PRIVATE_KEY="CKAajx/fJ35V/KkmyByBzZAN4EN+4vkkZGz9zpTjPE0="
ENV PROXY_BIND="127.0.0.1:25344"
ENV VPN_TIMEOUT="45"
ENV CDEP_TARGET="https://www.cdep.ro"

ENTRYPOINT ["/entrypoint.sh"]
