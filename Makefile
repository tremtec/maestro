.PHONY: build install test clean

APP_NAME = maestro
BIN_DIR = bin

build:
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) .

install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install .

test:
	@go test -v -race ./...

clean:
	@rm -rf $(BIN_DIR)

lint:
	@golangci-lint run ./...
