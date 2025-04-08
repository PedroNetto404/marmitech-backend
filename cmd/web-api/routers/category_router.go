package routers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/PedroNetto404/marmitech-backend/internal/app/usecase"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterCategoryRoutes(
	routerGroup *gin.RouterGroup,
	categoryUseCase usecase.ICategoryUseCase,
) {
	group := routerGroup.Group("/categories")
	group.POST("/", createcategory(categoryUseCase))
	group.PUT("/:id", updatecategory(categoryUseCase))
	group.GET("/:id", getcategoryById(categoryUseCase))
	group.GET("/", getProductCategories(categoryUseCase))
	group.DELETE("/:id", deletecategory(categoryUseCase))
	group.POST("/:id/picture", setcategoryImage(categoryUseCase))
	group.DELETE("/:id/picture", deletecategoryImage(categoryUseCase))
	group.POST("/:id/activate", activatecategory(categoryUseCase))
	group.POST("/:id/deactivate", deactivatecategory(categoryUseCase))
	group.POST("/:id/reorder/priority/:priority", reordercategory(categoryUseCase))
	group.POST("/:id/reorder/swap/:swapId", reordercategory(categoryUseCase))
}

func createcategory(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload usecase.CategoryPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		category, err := useCase.Create(&payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, category)
	}
}

func updatecategory(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var payload usecase.CategoryPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		category, err := useCase.Update(id, &payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func getcategoryById(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		category, err := useCase.FindById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func getProductCategories(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		findArgs := types.NewDefaultFindArgs()
		if err := c.ShouldBindQuery(&findArgs); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		restaurantId := c.Query("restaurant_id")
		if restaurantId == "" {
			c.JSON(http.StatusBadRequest, "restaurant_id is required")
			return
		}

		findArgs.Filter["restaurant_id"] = restaurantId

		productCategories, err := useCase.Find(findArgs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, productCategories)
	}
}

func deletecategory(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
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

func setcategoryImage(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
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

		content, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		picture := &types.FilePayload{
			Content:     content,
			ContentType: fileHeader.Header.Get("Content-Type"),
		}

		category, err := useCase.SetPicture(id, picture)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func deletecategoryImage(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		category, err := useCase.DeletePicture(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func activatecategory(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		category, err := useCase.Activate(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func deactivatecategory(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		category, err := useCase.Deactivate(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func reordercategory(useCase usecase.ICategoryUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		swapId := c.Param("swapId")
		if swapId != "" {
			if _, err := uuid.Parse(swapId); err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
		}

		priority, err := strconv.Atoi(c.Param("priority"))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		category, err := useCase.Reorder(id, priority, swapId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, category)
	}
}
