# ---- build stage ----
#
# KNOWN ENVIRONMENT BLOCKER (needs a human decision — see notes below):
# On a network that TLS-intercepts outbound HTTPS (as this machine does),
# `go mod download` FAILS here with:
#   tls: failed to verify certificate: x509: certificate signed by unknown authority
# The host Go toolchain works because the interceptor's root CA is trusted via
# the macOS keychain; this container has no such CA in its trust store, and no
# Go env var can paper over it. `docker compose build api` therefore cannot
# succeed as-is on this network.
#
# Two viable fixes, both of which are a judgement call for a human:
#   1. Vendor the module graph (`go mod vendor`, then build with
#      `-mod=vendor` and drop the `go mod download` step). Makes the image
#      build hermetic and network-free, but commits a large vendor/ tree —
#      and because the integration tests depend on testcontainers, that tree
#      pulls in the Docker client and OpenTelemetry purely for tests.
#   2. Inject the interceptor's root CA into this stage (COPY it in and set
#      SSL_CERT_FILE / run update-ca-certificates). Smaller change, but bakes
#      an organisation-specific certificate into the repo.
#
# Until one is chosen, run the API on the host against the composed Mongo:
#   docker compose up -d mongo && make seed && make run
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Dependencies first so the module cache layer survives source edits.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/seed   ./cmd/seed

# ---- runtime stage ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates curl \
    && adduser -D -u 10001 apiuser

COPY --from=builder /out/server /usr/local/bin/server
COPY --from=builder /out/seed   /usr/local/bin/seed

USER apiuser

ENV PORT=8090 \
    LOG_LEVEL=info

EXPOSE 8090

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD curl -fsS "http://localhost:${PORT}/health" || exit 1

ENTRYPOINT ["/usr/local/bin/server"]
