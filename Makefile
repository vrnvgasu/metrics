
.PHONY: build
build:
	@echo "build server"
	@go build -o cmd/server/server cmd/server/main.go cmd/server/flags.go
	@echo "build agent"
	@go build -o cmd/agent/agent cmd/agent/main.go cmd/agent/flags.go

.PHONY: test
test:
	@echo "test"
	@go test ./... -v