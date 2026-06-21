# ── Stage 1: wireproxy binary ────────────────────────────────────────────────
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

# ── Stage 2: Go build ─────────────────────────────────────────────────────────
FROM golang:1.23-bookworm AS build

# CGO is required for gosseract (Tesseract bindings)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libtesseract-dev \
    libleptonica-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY . .
# Use vendored dependencies to avoid network access during build
RUN CGO_ENABLED=1 GOOS=linux go build -mod=vendor -trimpath -ldflags="-s -w" -o /app/scraper ./cmd/scraper

# ── Stage 3: runtime ──────────────────────────────────────────────────────────
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    tesseract-ocr \
    tesseract-ocr-ron \
    poppler-utils \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=wireproxy-bin /tmp/wireproxy /usr/local/bin/wireproxy
COPY --from=build /app/scraper /scraper
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
