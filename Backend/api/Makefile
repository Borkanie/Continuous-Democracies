.PHONY: generate build run seed test test-unit docker-build up down snapshot restore

# Regenerates internal/generated/{types,server,spec}.go from openapi.yaml and
# keeps internal/swagger/openapi.yaml in sync (go:embed cannot reach outside
# its own package directory, so the spec is copied next to the handler).
generate:
	go tool oapi-codegen -config oapi-codegen.types.yaml openapi.yaml
	go tool oapi-codegen -config oapi-codegen.server.yaml openapi.yaml
	go tool oapi-codegen -config oapi-codegen.spec.yaml openapi.yaml
	cp openapi.yaml internal/swagger/openapi.yaml

build:
	go build ./...

run:
	go run ./cmd/server

seed:
	go run ./cmd/seed

test:
	go test ./tests/unit/... ./tests/integration/...

test-unit:
	go test ./tests/unit/...

# Builds the server+seed binaries on the host (so `go mod download` runs
# with the host's trusted CAs, not the container's) then packages them into
# the api image - see the comment at the top of ./Dockerfile.
docker-build:
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o server ./cmd/server
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o seed   ./cmd/seed
	docker build -t backend-api .
	rm -f server seed

up:
	docker compose -f ../docker-compose.yml up -d

down:
	docker compose -f ../docker-compose.yml down

snapshot:
	mongodump --uri=$${MONGODB_URI:-mongodb://localhost:27017} --archive=snapshot.gz --gzip

restore:
	mongorestore --uri=$${MONGODB_URI:-mongodb://localhost:27017} --archive=snapshot.gz --gzip --drop
