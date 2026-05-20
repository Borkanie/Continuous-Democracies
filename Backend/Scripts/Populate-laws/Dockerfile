# ── Stage 1: fetch wireproxy binary ──────────────────────────────────────────
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

# ── Stage 2: runtime ─────────────────────────────────────────────────────────
FROM python:3.11-slim

# Build-time args (can be set from CI using --build-arg)
ARG DBName
ARG DBUser
ARG DBPassword
ARG DBServerAddress
ARG DBPort=5432

# Runtime environment variables (populated from args by default)
ENV DBName=${DBName}
ENV DBUser=${DBUser}
ENV DBPassword=${DBPassword}
ENV DBServerAddress=${DBServerAddress}
ENV DBPort=${DBPort}

# Install system dependencies for Poppler and Tesseract + fonts
RUN apt-get update \
     && apt-get install -y --no-install-recommends \
         poppler-utils \
         tesseract-ocr \
         tesseract-ocr-ron \
         libgl1 \
         libsm6 \
         libxrender1 \
         libxext6 \
         libpq-dev \
         build-essential \
         python3-dev \
     && rm -rf /var/lib/apt/lists/*

# Userspace WireGuard → SOCKS5 proxy (used only for cdep.ro calls)
COPY --from=wireproxy-bin /tmp/wireproxy /usr/local/bin/wireproxy

WORKDIR /app

COPY requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

COPY . /app
RUN chmod +x /app/entrypoint.sh \
    && mkdir -p /app/temp

CMD ["/app/entrypoint.sh"]
