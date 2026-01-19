package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goNiki/Nerves/internal/dto"
)

func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Missing X-API-Key header",
			))
			return
		}

		if key != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Invalid API key",
			))
			return
		}

		c.Next()
	}
}
