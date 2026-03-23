.PHONY: deps run-api run-admin run-worker fmt

deps:
	go mod tidy

run-api:
	go run ./cmd/api

run-admin:
	go run ./cmd/admin

run-worker:
	go run ./cmd/worker

fmt:
	gofmt -w ./cmd ./internal
