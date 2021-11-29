.PHONY: run
run:
	go run ./cmd/crawler/main.go

.PHONY: test
test:
	go test -v ./...
	
.PHONY: lint
lint:
	golangci-lint run ./...

.DEFAULT_GOAL := run