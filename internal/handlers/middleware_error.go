package handlers

import (
	"errors"
	"net/http"

	"lk/internal/logger"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware - это middleware для централизованной обработки ошибок.
// Является обычной функцией, так как не требует доступа к Handler.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *services.AppError
		if errors.As(err, &appErr) {
			c.AbortWithStatusJSON(appErr.StatusCode, errorResponse{Message: appErr.Message})
		} else {
			logger.Default().WithField("module", "GIN_ERROR").WithError(err).Error("Unhandled error occurred")
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{Message: "Internal Server Error"})
		}
	}
}
