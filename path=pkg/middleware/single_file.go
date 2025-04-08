package middleware

import (
	"io"
	"net/http"
	"golang.org/x/exp/slices"

	"github.com/PedroNetto404/marmitech-backend/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func isValidMimeType(mimeType string, allowedMimeTypes []string) bool {
	return slices.Contains(allowedMimeTypes, mimeType)
} 