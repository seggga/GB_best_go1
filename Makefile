.PHONY: run
run:
	go run ./main.go

.PHONY: test
test:
	go test -v ./...
	
.DEFAULT_GOAL := run