package handlers

import (
	"errors"
	"strings"

	"lk/internal/models"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	// userProfileCtx - ключ, по которому в контексте хранится профиль пользователя.
	userProfileCtx = "userProfile"
)

// userIdentity - это middleware для проверки JWT и загрузки профиля пользователя в контекст.
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.Error(services.NewUnauthorizedError("empty auth header", nil))
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.Error(services.NewUnauthorizedError("invalid auth header", nil))
		c.Abort()
		return
	}

	if len(headerParts[1]) == 0 {
		c.Error(services.NewUnauthorizedError("token is empty", nil))
		c.Abort()
		return
	}

	userID, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		c.Error(err) // ParseToken возвращает типизированную ошибку
		c.Abort()
		return
	}

	// После успешной валидации токена, загружаем профиль пользователя
	userProfile, err := h.userRepo.GetUserProfileByUserID(c.Request.Context(), userID)
	if err != nil {
		c.Error(services.NewUnauthorizedError("user profile not found for this token", err))
		c.Abort()
		return
	}

	// Записываем весь профиль в контекст Gin.
	c.Set(userProfileCtx, userProfile)
}

// getUserProfile - вспомогательная функция для извлечения профиля пользователя из контекста.
func getUserProfile(c *gin.Context) (models.UserProfile, error) {
	profile, ok := c.Get(userProfileCtx)
	if !ok {
		return models.UserProfile{}, errors.New("user profile not found in context (middleware error)")
	}

	userProfile, ok := profile.(models.UserProfile)
	if !ok {
		return models.UserProfile{}, errors.New("user profile is of invalid type in context")
	}

	return userProfile, nil
}
