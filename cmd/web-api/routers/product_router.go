package routers

import (
	"net/http"

	"github.com/PedroNetto404/marmitech-backend/internal/app/usecase"
	"github.com/PedroNetto404/marmitech-backend/pkg/middleware"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(
	group *gin.RouterGroup,
	productUseCase usecase.IProductUseCase,
) {
	productGroup := group.Group("/products")
	
	productGroup.GET("/", createProduct(productUseCase))
	productGroup.GET("/:id", getProduct(productUseCase))
	productGroup.POST("/", createProduct(productUseCase))
	productGroup.PUT("/:id", updateProduct(productUseCase))
	productGroup.DELETE("/:id", deleteProduct(productUseCase))
	productGroup.PATCH("/:id/picture", updateProductPicture(productUseCase))
	productGroup.DELETE("/:id/picture", deleteProductPicture(productUseCase))
}

func createProduct(productUseCase usecase.IProductUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload usecase.ProductPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		created, err := productUseCase.Create(&payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(http.StatusCreated, created)
	}
}

func getProduct(productUseCase usecase.IProductUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		product, err := productUseCase.FindById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func updateProduct(productUseCase usecase.IProductUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var payload usecase.ProductPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		product, err := productUseCase.Update(id, &payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func deleteProduct(productUseCase usecase.IProductUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		err := productUseCase.Delete(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func updateProductPicture(productUseCase usecase.IProductUseCase) gin.HandlerFunc {
	return middleware.SingleFileMiddleware(
		[]string{"image/jpeg", "image/png"},
		5*1024*1024,
		func(c *gin.Context, file *types.FilePayload) {
			id := c.Param("id")
			if id == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
				return
			}

			product, err := productUseCase.SetPicture(id, file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product picture"})
				return
			}

			c.JSON(http.StatusOK, product)
		},
	)	
}

func deleteProductPicture(productUseCase usecase.IProductUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		product, err := productUseCase.DeletePicture(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product picture"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}
