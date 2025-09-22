package handlers

import (
	"errors"
	"net/http"
	"strings"

	"lk/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userProfileCtx      = "userProfile"
)

// userIdentity - это middleware для проверки JWT и загрузки профиля пользователя в контекст.
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

	// После успешной валидации токена, загружаем профиль пользователя
	userProfile, err := h.userRepo.GetUserProfileByUserID(c.Request.Context(), userID)
	if err != nil {
		// Если профиль не найден для валидного токена - это тоже ошибка авторизации
		newErrorResponse(c, http.StatusUnauthorized, "user profile not found")
		return
	}

	// Записываем весь профиль в контекст Gin.
	c.Set(userProfileCtx, userProfile)
}

// getUserProfile - вспомогательная функция для извлечения профиля пользователя из контекста.
func getUserProfile(c *gin.Context) (models.UserProfile, error) {
	profile, ok := c.Get(userProfileCtx)
	if !ok {
		return models.UserProfile{}, errors.New("user profile not found in context")
	}

	userProfile, ok := profile.(models.UserProfile)
	if !ok {
		return models.UserProfile{}, errors.New("user profile is of invalid type in context")
	}

	return userProfile, nil
}
