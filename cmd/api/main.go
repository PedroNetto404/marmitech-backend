package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/magalu/projects/marmitech-backend/internal/infra/api"
	"github.com/magalu/projects/marmitech-backend/internal/infra/config"
)

func main() {
	// Carrega as configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Configura o modo do Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializa o router
	router := gin.Default()

	// Configura o Swagger
	api.SetupSwagger(router)

	// Configura as rotas da API
	api.SetupRoutes(router)

	// Inicia o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
