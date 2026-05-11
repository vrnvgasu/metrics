
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

.PHONY: bench
bench:
	@echo "test benchmarks"
	@go test -bench=. -benchmem ./...

.PHONY: cover
cover:
	@echo "coverage"
	@go test -coverprofile=coverage.out \
		-coverpkg=$(shell go list ./... | grep -v "mocks" | paste -sd,) \
		$(shell go list ./... | grep -v "mocks")
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

# Нагрузка на уже запущенный сервер (localhost:8080)
.PHONY: load
load:
	@echo "running load: 4 threads, 20 connections, 15s"
	@wrk -t4 -c20 -d15s -s scripts/wrk_updates.lua http://localhost:8080/updates

# Запуск сервера, нагрузка, снятие base.pprof, остановка сервера
.PHONY: profile-base
profile-base:
	@mkdir -p profiles
	@echo "starting server..."
	@go run ./cmd/server &> /tmp/metrics-server.log & echo $$! > /tmp/metrics-server.pid
	@sleep 3
	@echo "running load..."
	@wrk -t4 -c20 -d15s -s scripts/wrk_updates.lua http://localhost:8080/updates
	@echo "saving base.pprof..."
	@curl -s http://localhost:6060/debug/pprof/heap > profiles/base.pprof
	@kill $$(cat /tmp/metrics-server.pid) 2>/dev/null; rm -f /tmp/metrics-server.pid
	@echo "done: profiles/base.pprof"

# Запуск сервера, нагрузка, снятие result.pprof, остановка сервера
.PHONY: profile-result
profile-result:
	@mkdir -p profiles
	@echo "starting server..."
	@go run ./cmd/server &> /tmp/metrics-server.log & echo $$! > /tmp/metrics-server.pid
	@sleep 3
	@echo "running load..."
	@wrk -t4 -c20 -d15s -s scripts/wrk_updates.lua http://localhost:8080/updates
	@echo "saving result.pprof..."
	@curl -s http://localhost:6060/debug/pprof/heap > profiles/result.pprof
	@kill $$(cat /tmp/metrics-server.pid) 2>/dev/null; rm -f /tmp/metrics-server.pid
	@echo "done: profiles/result.pprof"

# Сравнение профилей
.PHONY: profile-diff
profile-diff:
	@go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof