# Variáveis configuráveis via CLI
PROJECT ?= web-api
BIN_DIR := bin
CMD_DIR := cmd/$(PROJECT)
APP_NAME := marmitech-$(PROJECT)
DB_URL ?= postgres://postgres:postgres@localhost:5432/marmitech?sslmode=disable

# Targets
.PHONY: build run migrate up down logs fmt lint test tidy clean rebuild

## Builda o projeto indicado em PROJECT (default: web-api)
build:
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

## Executa o binário buildado
run: build
	./$(BIN_DIR)/$(APP_NAME)

## Executa a main de migrations
migrate:
	go run cmd/migrate/main.go

## Sobe containers Docker
up:
	docker-compose up -d

## Para e remove containers Docker
down:
	docker-compose down

## Mostra logs dos serviços
logs:
	docker-compose logs -f

## Formata os arquivos do projeto
fmt:
	go fmt ./...

lint:
	golangci-lint run ./

test:
	go test -v ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)

rebuild: clean tidy build
