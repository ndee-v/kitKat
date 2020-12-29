PHONY: build
-B build:
			go build -o ./build/tcpServer -v ./cmd/tcpServer/
PHONY: test
test:
			go test -v -timeout 30s ./...

.DEFAULT_GOAL := build