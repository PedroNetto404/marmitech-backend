# Variáveis configuráveis via CLI
PROJECT ?= web-api
BIN_DIR := bin
CMD_DIR := cmd/$(PROJECT)
APP_NAME := marmitech-$(PROJECT)
DB_URL ?= postgres://postgres:postgres@localhost:5432/marmitech?sslmode=disable
ENV_FILE ?= .env

# Targets
.PHONY: all build run test clean lint migrate setup help swagger

# Cores para output
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m

all: setup build

## Help: Exibe todos os comandos disponíveis
help:
	@echo "$(GREEN)Marmitech Backend - Comandos disponíveis:$(NC)"
	@echo "$(YELLOW)make setup$(NC)     - Configura o ambiente de desenvolvimento"
	@echo "$(YELLOW)make build$(NC)     - Compila o projeto"
	@echo "$(YELLOW)make run$(NC)       - Executa a aplicação"
	@echo "$(YELLOW)make test$(NC)      - Executa os testes"
	@echo "$(YELLOW)make lint$(NC)      - Executa análise estática de código"
	@echo "$(YELLOW)make migrate$(NC)   - Executa as migrações do banco de dados"
	@echo "$(YELLOW)make clean$(NC)     - Limpa arquivos gerados"
	@echo "$(YELLOW)make sonar$(NC)     - Executa análise SonarQube"
	@echo "$(YELLOW)make swagger$(NC)   - Gera documentação Swagger"

## Setup: Configura o ambiente de desenvolvimento
setup:
	@echo "$(GREEN)Configurando ambiente de desenvolvimento...$(NC)"
	@if [ ! -f $(ENV_FILE) ]; then \
		cp .env.example $(ENV_FILE); \
		echo "$(YELLOW)Arquivo .env criado. Por favor, configure as variáveis de ambiente.$(NC)"; \
	fi
	@go mod download
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)Ambiente configurado com sucesso!$(NC)"

## Build: Compila o projeto
build:
	@echo "$(GREEN)Compilando projeto...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	@echo "$(GREEN)Projeto compilado com sucesso!$(NC)"

## Run: Executa a aplicação
run: build
	@echo "$(GREEN)Iniciando aplicação...$(NC)"
	@./$(BIN_DIR)/$(APP_NAME)

## Test: Executa os testes
test:
	@echo "$(GREEN)Executando testes...$(NC)"
	@go test -v -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Relatório de cobertura gerado em coverage.html$(NC)"

## Lint: Executa análise estática de código
lint:
	@echo "$(GREEN)Executando análise estática de código...$(NC)"
	@golangci-lint run
	@echo "$(GREEN)Análise concluída!$(NC)"

## Migrate: Executa as migrações do banco de dados
migrate:
	@echo "$(GREEN)Executando migrações...$(NC)"
	@go run cmd/migrate/main.go
	@echo "$(GREEN)Migrações concluídas!$(NC)"

## Sonar: Executa análise SonarQube
sonar:
	@echo "$(GREEN)Executando análise SonarQube...$(NC)"
	@sonar-scanner
	@echo "$(GREEN)Análise SonarQube concluída!$(NC)"

## Clean: Limpa arquivos gerados
clean:
	@echo "$(GREEN)Limpando arquivos gerados...$(NC)"
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Limpeza concluída!$(NC)"

## Docker: Comandos Docker
docker-build:
	@echo "$(GREEN)Construindo imagem Docker...$(NC)"
	@docker build -t $(APP_NAME) .

docker-run:
	@echo "$(GREEN)Executando container Docker...$(NC)"
	@docker run -p 8080:8080 --env-file $(ENV_FILE) $(APP_NAME)

## Swagger: Gera documentação Swagger
swagger:
	@echo "$(GREEN)Gerando documentação Swagger...$(NC)"
	@swag init -g $(CMD_DIR)/main.go -o docs
	@echo "$(GREEN)Documentação Swagger gerada!$(NC)"
