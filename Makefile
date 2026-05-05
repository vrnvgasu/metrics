
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

.PHONY: cover
cover:
	@echo "coverage"
	@go test -coverprofile=coverage.out -coverpkg=./... $(shell go list ./... | grep -v "internal/repository/mocks")
	@go tool cover -func=coverage.out | grep "^total:"

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