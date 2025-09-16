package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// authorizationHeader - ключ заголовка для передачи токена.
	authorizationHeader = "Authorization"
	// userCtx - ключ, по которому в контексте запроса хранится ID пользователя.
	userCtx = "userID"
)

// userIdentity - это middleware для проверки JWT и идентификации пользователя.
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userID, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Записываем ID пользователя в контекст Gin.
	// Теперь этот ID будет доступен во всех последующих обработчиках этого запроса.
	c.Set(userCtx, userID)
}

// getUserID - вспомогательная функция для извлечения ID пользователя из контекста.
func getUserID(c *gin.Context) (uint64, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}

	idUint64, ok := id.(uint64)
	if !ok {
		return 0, errors.New("user ID is of invalid type")
	}

	return idUint64, nil
}
