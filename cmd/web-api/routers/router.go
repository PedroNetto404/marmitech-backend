package routers

import (
	"github.com/PedroNetto404/marmitech-backend/internal/app/usecase"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	engine *gin.Engine,
	categoryUseCase usecase.ICategoryUseCase,
	productUseCase usecase.IProductUseCase,
	dishUseCase usecase.IDishUseCase,
	restaurantUseCase usecase.IRestaurantUseCase,
) {
	apiGroup := engine.Group("/api")

	registerV1(apiGroup, categoryUseCase, productUseCase, dishUseCase, restaurantUseCase)
}

func registerV1(
	apiGroup *gin.RouterGroup,
	categoryUseCase usecase.ICategoryUseCase,
	productUseCase usecase.IProductUseCase,
	dishUseCase usecase.IDishUseCase,
	restaurantUseCase usecase.IRestaurantUseCase,
) {
	v1Group := apiGroup.Group("/v1/restaurants")
	RegisterRestaurantRoutes(v1Group, restaurantUseCase)
	
	restaurantGroup := v1Group.Group("/:restaurantId")
	RegisterCategoryRoutes(restaurantGroup, categoryUseCase)
	RegisterProductRoutes(restaurantGroup, productUseCase)
	RegisterDishRoutes(restaurantGroup, dishUseCase)
}
