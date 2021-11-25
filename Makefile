.PHONY: run
run:
	go run ./cmd/crawler/main.go

.PHONY: test
test:
	go test -v ./...
	
.DEFAULT_GOAL := run