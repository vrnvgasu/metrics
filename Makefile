
.PHONY: build
build:
	@echo "build server"
	@go build -o cmd/server/server cmd/server/main.go
	@echo "build agent"
	@go build -o cmd/agent/agent cmd/agent/main.go

.PHONY: test
test:
	@echo "test"
	@go test ./... -v