FROM golang:1.23-bookworm AS build

WORKDIR /src
# Copy all source first so go mod tidy can resolve imports
COPY . .
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/server ./cmd/server

FROM gcr.io/distroless/static-debian12

COPY --from=build /app/server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
