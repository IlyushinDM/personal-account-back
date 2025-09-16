package handlers

import (
	"github.com/gin-gonic/gin"
)

// errorResponse - структура для ответа с ошибкой.
type errorResponse struct {
	Message string `json:"message"`
}

// statusResponse - структура для ответа с сообщением о статусе.
// type statusResponse struct {
// 	Status string `json:"status"`
// }

// newErrorResponse отправляет стандартизированный ответ об ошибке.
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
