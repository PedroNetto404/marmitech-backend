package routers

import (
	"net/http"

	"github.com/PedroNetto404/marmitech-backend/internal/app/usecase"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
	"github.com/gin-gonic/gin"
)

func RegisterDishRoutes(r *gin.RouterGroup, useCase usecase.IDishUseCase) {
	group := r.Group("/dishes")

	group.POST("/", createDish(useCase))
	group.PUT("/:id", updateDish(useCase))
	group.GET("/:id", getDishById(useCase))
	group.GET("/", getDishes(useCase))
	group.DELETE("/:id", deleteDish(useCase))
	group.POST("/:id/picture", setDishImage(useCase))
}

func createDish(useCase usecase.IDishUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload usecase.DishPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		dish, err := useCase.Create(&payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, dish)
	}
}

func updateDish(useCase usecase.IDishUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var payload usecase.DishPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		dish, err := useCase.Update(id, &payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, dish)
	}
}

func getDishById(useCase usecase.IDishUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		dish, err := useCase.FindById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		if dish == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		c.JSON(http.StatusOK, dish)
	}
}

func getDishes(useCase usecase.IDishUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		findArgs := types.NewDefaultFindArgs()
		if err := c.ShouldBindQuery(&findArgs); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		dishes, err := useCase.Find(findArgs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, dishes)
	}
}

func deleteDish(useCase usecase.IDishUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := useCase.Delete(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func setDishImage(useCase usecase.IDishUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")


		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		var content []byte
		if _, err := file.Read(content); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		payload := types.FilePayload{
			Content: content,
			ContentType: fileHeader.Header.Get("Content-Type"),
		}

		dish, err := useCase.SetPicture(id, &payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, dish)
	}
}
