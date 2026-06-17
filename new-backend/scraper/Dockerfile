FROM golang:1.23-bookworm AS build

# CGO is required for gosseract (Tesseract bindings)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libtesseract-dev \
    libleptonica-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
# Copy all source first so go mod tidy can resolve imports
COPY . .
RUN go mod tidy && go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/scraper ./cmd/scraper

# Runtime image: needs Tesseract + Poppler at runtime
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    tesseract-ocr \
    tesseract-ocr-ron \
    poppler-utils \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/scraper /scraper

ENTRYPOINT ["/scraper"]
