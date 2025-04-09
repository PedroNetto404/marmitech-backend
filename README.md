# Marmitech Backend

Backend da plataforma Marmitech, um sistema completo para gestão de restaurantes e delivery.

## 🚀 Visão Geral

O Marmitech é uma plataforma robusta para gestão de restaurantes, oferecendo funcionalidades como:
- Controle de vendas e pedidos
- Gestão de estoque
- Cadastro e fidelização de clientes
- Controle de entregas
- Gestão financeira
- Relatórios e análises

## 🛠 Tecnologias

- **Linguagem**: Go (Golang)
- **Banco de Dados**: PostgreSQL
- **ORM**: GORM
- **Autenticação**: JWT
- **Documentação**: Swagger/OpenAPI
- **Testes**: Go testing package
- **Qualidade de Código**: SonarQube, golangci-lint
- **Versionamento**: Git

## 📁 Estrutura do Projeto

```
marmitech-backend/
├── cmd/                    # Ponto de entrada da aplicação
│   ├── web-api/           # API HTTP
│   └── migrate/           # Migrações do banco de dados
├── internal/              # Código interno da aplicação
│   ├── app/              # Casos de uso e regras de negócio
│   ├── config/           # Configurações da aplicação
│   ├── domain/           # Entidades e interfaces do domínio
│   └── infra/            # Implementações de infraestrutura
└── pkg/                  # Código compartilhado
    ├── database/         # Configurações do banco de dados
    ├── enums/            # Enumeradores
    ├── logger/           # Configurações de log
    ├── middleware/       # Middlewares HTTP
    └── types/            # Tipos compartilhados
```

## ⚙️ Configuração

### Variáveis de Ambiente

```bash
# Configurações do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=marmitech
DB_SSL_MODE=disable

# Configurações da API
API_PORT=8080
API_ENV=development

# Configurações de Autenticação
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION=24h

# Configurações de Log
LOG_LEVEL=info
LOG_FORMAT=json
```

### Instalação e Uso

O projeto utiliza Make para automatizar tarefas comuns. Para ver todos os comandos disponíveis:

```bash
make help
```

#### Comandos Principais

1. **Configuração Inicial**:
```bash
make setup
```

2. **Desenvolvimento**:
```bash
# Compilar o projeto
make build

# Executar a aplicação
make run

# Executar testes
make test

# Executar análise estática de código
make lint
```

3. **Banco de Dados**:
```bash
# Executar migrações
make migrate
```

4. **Qualidade de Código**:
```bash
# Executar análise SonarQube
make sonar

# Gerar documentação Swagger
make swagger
```

5. **Docker**:
```bash
# Construir imagem Docker
make docker-build

# Executar container Docker
make docker-run
```

6. **Limpeza**:
```bash
# Limpar arquivos gerados
make clean
```

## 🧪 Testes

O projeto utiliza o pacote de testes padrão do Go. Para executar os testes:

```bash
# Executar todos os testes
make test

# O comando acima irá:
# 1. Executar os testes com verbose
# 2. Gerar relatório de cobertura
# 3. Criar arquivo HTML com visualização da cobertura
```

## 📚 Documentação

A documentação da API está disponível via Swagger UI. Para gerar/atualizar a documentação:

```bash
make swagger
```

Após iniciar a aplicação, a documentação estará disponível em `http://localhost:8080/swagger/index.html`.

## 🔍 Qualidade de Código

O projeto utiliza várias ferramentas para garantir a qualidade do código:

- **golangci-lint**: Para análise estática de código
- **SonarQube**: Para análise de qualidade e cobertura
- **EditorConfig**: Para padronização de estilo

Para executar as análises:
```bash
make lint    # Executa golangci-lint
make sonar   # Executa análise SonarQube
```

## 🤝 Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request
