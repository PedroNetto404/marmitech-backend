# Variáveis configuráveis via CLI
PROJECT ?= web-api
BIN_DIR := bin
CMD_DIR := cmd/$(PROJECT)
APP_NAME := marmitech-$(PROJECT)
DB_URL ?= postgres://postgres:postgres@localhost:5432/marmitech?sslmode=disable
ENV_FILE ?= .env

# Targets
.PHONY: all build run test clean lint migrate setup help swagger

all: setup build

## Help: Exibe todos os comandos disponíveis
help:
	@echo "Marmitech Backend - Comandos disponíveis:"
	@echo "make setup     - Configura o ambiente de desenvolvimento"
	@echo "make build     - Compila o projeto"
	@echo "make run       - Executa a aplicação"
	@echo "make test      - Executa os testes"
	@echo "make lint      - Executa análise estática de código"
	@echo "make migrate   - Executa as migrações do banco de dados"
	@echo "make clean     - Limpa arquivos gerados"
	@echo "make sonar     - Executa análise SonarQube"
	@echo "make swagger   - Gera documentação Swagger"

## Setup: Configura o ambiente de desenvolvimento
setup:
	@echo "Configurando ambiente de desenvolvimento..."
	@if [ ! -f $(ENV_FILE) ]; then \
		cp .env.example $(ENV_FILE); \
		echo "Arquivo .env criado. Por favor, configure as variáveis de ambiente."; \
	fi
	@go mod download
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Ambiente configurado com sucesso!"

## Build: Compila o projeto
build:
	@echo "Compilando projeto..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	@echo "Projeto compilado com sucesso!"

## Run: Executa a aplicação
run: build
	@echo "Iniciando aplicação..."
	@./$(BIN_DIR)/$(APP_NAME)

## Test: Executa os testes
test:
	@echo "Executando testes..."
	@go test -v -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado em coverage.html"

## Lint: Executa análise estática de código
lint:
	@echo "Executando análise estática de código..."
	@golangci-lint run
	@echo "Análise concluída!"

## Migrate: Executa as migrações do banco de dados
migrate:
	@echo "Executando migrações..."
	@go run cmd/migrate/main.go
	@echo "Migrações concluídas!"

## Sonar: Executa análise SonarQube
sonar:
	@echo "Executando análise SonarQube..."
	@sonar-scanner
	@echo "Análise SonarQube concluída!"

## Clean: Limpa arquivos gerados
clean:
	@echo "Limpando arquivos gerados..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@echo "Limpeza concluída!"

## Docker: Comandos Docker
docker-build:
	@echo "Construindo imagem Docker..."
	@docker build -t $(APP_NAME) .

docker-run:
	@echo "Executando container Docker..."
	@docker run -p 8080:8080 --env-file $(ENV_FILE) $(APP_NAME)

## Swagger: Gera documentação Swagger
swagger:
	@echo "Gerando documentação Swagger..."
	@swag init -g $(CMD_DIR)/main.go -o docs
	@echo "Documentação Swagger gerada!"
