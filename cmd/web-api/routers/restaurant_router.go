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

	group.POST("/", create(useCase))
	group.PUT("/:id", update(useCase))
	group.GET("/:id", getById(useCase))
	group.GET("/", getAll(useCase))
	group.DELETE("/:id", delete(useCase))
	group.POST("/:id/images", setImages(useCase))
}

func create(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
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

func update(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
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

func getById(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
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

func getAll(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.FindArgs
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		restaurants, err := useCase.GetAll(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, restaurants)
	}
}

func delete(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := useCase.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func setImages(useCase usecase.IRestaurantUseCase) gin.HandlerFunc {
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

