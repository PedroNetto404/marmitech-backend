package middleware

import (
	"io"
	"net/http"

	"github.com/PedroNetto404/marmitech-backend/pkg/types"
	"github.com/gin-gonic/gin"
)

func SingleFileMiddleware(
	allowedMimeTypes []string,
	maxSize int64,
	handler func(c *gin.Context, file *types.FilePayload),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
			c.Abort()
			return
		}

		if fileHeader.Size > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
			c.Abort()
			return
		}

		if !isValidMimeType(fileHeader.Header.Get("Content-Type"), allowedMimeTypes) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			c.Abort()
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
			c.Abort()
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file"})
			c.Abort()
			return
		}

		filePayload := &types.FilePayload{
			Content:     content,
			ContentType: fileHeader.Header.Get("Content-Type"),
		}

		handler(c, filePayload)
	}
}

func isValidMimeType(mimeType string, allowedMimeTypes []string) bool {
	for _, allowed := range allowedMimeTypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}
