# build and run the application
.PHONY: run
run:
	go run ./cmd/crawler/main.go

# build binary
.PHONY: build
build: test lint
	go build -o crawler ./cmd/crawler/main.go 

# run tests
.PHONY: test
test:
	go test -v ./...

# run linters 
.PHONY: lint
lint:
	golangci-lint run ./...

.DEFAULT_GOAL := run