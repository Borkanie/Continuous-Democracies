.PHONY: generate build run seed test test-unit up down snapshot restore

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

up:
	docker compose up -d

down:
	docker compose down

snapshot:
	mongodump --uri=$${MONGODB_URI:-mongodb://localhost:27017} --archive=snapshot.gz --gzip

restore:
	mongorestore --uri=$${MONGODB_URI:-mongodb://localhost:27017} --archive=snapshot.gz --gzip --drop
