APP_NAME := go-todo-service
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

.PHONY: all build run dev test fmt tidy docker-up docker-down clean

all: build

$(BIN):
	@mkdir -p $(BIN_DIR)
	GO111MODULE=on go build -o $(BIN) ./cmd/api

build: $(BIN)

run: build
	./$(BIN)

dev:
	go run ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*' -not -path './bin/*')

tidy:
	go mod tidy

docker-up:
	docker compose up --build

docker-down:
	docker compose down

clean:
	rm -rf $(BIN_DIR)
