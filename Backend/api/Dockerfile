# The Go binary is built on the host (via `make docker-build`), not inside
# this image: this network TLS-intercepts outbound HTTPS, and the container
# doesn't trust that proxy's root CA the way the host does (macOS keychain),
# so `go mod download` inside a build stage fails here. This Dockerfile just
# packages a pre-built binary - no Go toolchain, no network access needed.
FROM alpine:3.21

RUN apk add --no-cache ca-certificates curl \
    && adduser -D -u 10001 apiuser

COPY server /usr/local/bin/server
COPY seed   /usr/local/bin/seed

USER apiuser

ENV PORT=8090 \
    LOG_LEVEL=info

EXPOSE 8090

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD curl -fsS "http://localhost:${PORT}/health" || exit 1

ENTRYPOINT ["/usr/local/bin/server"]
