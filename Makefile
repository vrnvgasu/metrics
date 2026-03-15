
.PHONY: build
build:
	@echo "build server"
	@go build -o cmd/server/server cmd/server/main.go cmd/server/flags.go cmd/server/app.go
	@echo "build agent"
	@go build -o cmd/agent/agent cmd/agent/main.go cmd/agent/flags.go

.PHONY: test
test:
	@echo "test"
	@go test ./... -v

.PHONY: generate
generate:
	@echo "generate"
	@go generate ./...

.PHONY: infra-up
infra-up:
	@echo "local-up"
	@docker compose -f deployments/docker-compose.yaml up -d

.PHONY: infra-down
infra-down:
	@echo "local-down"
	@docker compose -f deployments/docker-compose.yaml down