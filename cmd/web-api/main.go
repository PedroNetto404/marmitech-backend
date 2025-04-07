package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/PedroNetto404/marmitech-backend/cmd/web-api/routers"
	"github.com/PedroNetto404/marmitech-backend/internal/app/usecase"
	"github.com/PedroNetto404/marmitech-backend/internal/config"
	"github.com/PedroNetto404/marmitech-backend/internal/infra/files"
	"github.com/PedroNetto404/marmitech-backend/internal/infra/respositories"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/middleware"
)

func main() {
	config.LoadEnvs()

	engine := gin.New()

	engine.Use(middleware.HttpLogger())
	engine.Use(middleware.GlobalErrorHandler())

	timeout, err := strconv.Atoi(config.Env.ApiRouteTimeout)
	if err == nil {
		engine.Use(middleware.RouteTimeout(timeout))
	}

	baseGroup := engine.Group(config.Env.ApiBasePath)

	// Services
	fileStorage := files.NewCloudBlockStorage()

	// Repositories
	db, err := database.New()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("❌ Failed to close database connection: %v", err)
		}
	}()
	restaurantRepository := respositories.NewRestaurantRepository(db)

	// Use Cases
	restaurantUseCase := usecase.NewRestaurantUseCase(restaurantRepository, fileStorage)

	// Routers
	routers.RegisterRestaurantRoutes(baseGroup, restaurantUseCase)

	addr := fmt.Sprintf("%s:%s", config.Env.ApiHost, config.Env.ApiPort)

	log.Printf("🚀 Starting server on %s", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
