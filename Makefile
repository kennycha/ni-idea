.PHONY: build install test clean

BINARY_NAME=ni
BUILD_DIR=./cmd/ni

build:
	go build -o $(BINARY_NAME) $(BUILD_DIR)

install:
	go install $(BUILD_DIR)

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)
	go clean

tidy:
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run
