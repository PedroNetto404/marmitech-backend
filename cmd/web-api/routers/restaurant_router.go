package routers

import (
	"net/http"

	"github.com/PedroNetto404/marmitech-backend/internal/app/dtos"
	"github.com/PedroNetto404/marmitech-backend/internal/app/usecase"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
	"github.com/gin-gonic/gin"
)

func RegisterRestaurantRoutes(r *gin.RouterGroup, useCase usecase.IRestaurantUseCase) {
	group := r.Group("/restaurants")

	group.POST("/", createRestaurant(useCase))
	group.PUT("/:id", updateRestaurant(useCase))
	group.GET("/:id", getRestaurantById(useCase))
	group.GET("/", getAllRestaurants(useCase))
	group.DELETE("/:id", deleteRestaurant(useCase))
	group.POST("/:id/images", setRestaurantImages(useCase))
}

func createRestaurant(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.RestaurantPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		restaurant, err := useCase.Create(&payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, restaurant)
	}
}

func updateRestaurant(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var payload dtos.RestaurantPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		restaurant, err := useCase.Update(id, &payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, restaurant)
	}
}

func getRestaurantById(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		restaurant, err := useCase.GetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, restaurant)
	}
}

func getAllRestaurants(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryArgs := types.NewDefaultFindArgs()
		if err := c.ShouldBindQuery(&queryArgs); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		restaurants, err := useCase.GetAll(queryArgs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, restaurants)
	}
}

func deleteRestaurant(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := useCase.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func setRestaurantImages(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var payload dtos.SetRestaurantImagesPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		restaurant, err := useCase.SetImages(id, &payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, restaurant)
	}
}
