package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/PedroNetto404/marmitech-backend/cmd/web-api/routers"
	"github.com/PedroNetto404/marmitech-backend/internal/app/usecase"
	"github.com/PedroNetto404/marmitech-backend/internal/config"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/internal/infra/files"
	"github.com/PedroNetto404/marmitech-backend/internal/infra/respositories"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/middleware"
)

var (
)

func main() {
	config.LoadEnvs()

	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		if c.Request.URL.Path[len(c.Request.URL.Path)-1] == '/' {
			c.Request.URL.Path = c.Request.URL.Path[:len(c.Request.URL.Path)-1]
		}
		c.Next()			
	})

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.Use(gin.ErrorLogger())
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
	}))
	engine.Use(middleware.HttpLogger())
	engine.Use(middleware.GlobalErrorHandler())

	timeout, err := strconv.Atoi(config.Env.ApiRouteTimeout)
	if err == nil {
		engine.Use(middleware.RouteTimeout(timeout))
	}

	baseGroup := engine.Group(config.Env.ApiBasePath)

	// Services
	var blockStorage ports.IBlockStorage
	if config.Env.DiskStoragePath == "" {
		blockStorage = files.NewCloudBlockStorage()
	} else {
		blockStorage = files.NewDiskStorage(config.Env.DiskStoragePath)
	}

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
	dishRepository := respositories.NewDishRepository(db)
	categoryRepository := respositories.NewCategoryRepository(db)
	productRepository := respositories.NewProductRepository(db)
	
	// Use Cases
	restaurantUseCase := usecase.NewRestaurantUseCase(restaurantRepository, blockStorage)
	dishUseCase := usecase.NewDishUseCase(dishRepository, restaurantRepository, blockStorage)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository, restaurantRepository, blockStorage)
	productUseCase := usecase.NewProductUseCase(productRepository, categoryRepository, restaurantRepository, blockStorage)
	
	// Routers
	routers.RegisterRestaurantRoutes(baseGroup, restaurantUseCase)
	routers.RegisterDishRoutes(baseGroup, dishUseCase)
	routers.RegisterCategoryRoutes(baseGroup, categoryUseCase)
	routers.RegisterProductRoutes(baseGroup, productUseCase)

	addr := fmt.Sprintf("%s:%s", config.Env.ApiHost, config.Env.ApiPort)

	log.Printf("🚀 Starting server on %s", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
