package http

import (
	"net/http"

	"lk/internal/models"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить профиль пользователя
// @Security     ApiKeyAuth
// @Tags         profile
// @Description  Возвращает полную информацию о профиле текущего пользователя и его предстоящие записи.
// @Id           get-profile
// @Produce      json
// @Success      200 {object} map[string]interface{} "profile, appointments"
// @Failure      401,500 {object} errorResponse
// @Router       /profile [get]
func (h *Handler) getProfile(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		c.Error(services.NewInternalServerError("failed to get user from context", err))
		return
	}

	// Профиль уже есть из middleware, нужно только получить записи
	_, appointments, err := h.services.User.GetFullUserProfile(c.Request.Context(), userProfile.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile":      userProfile,
		"appointments": appointments,
	})
}

// updateUserProfileInput - DTO для обновления профиля.
type updateUserProfileInput struct {
	Email  *string `json:"email"`
	CityID *uint32 `json:"cityID"`
}

// @Summary      Обновить профиль пользователя
// @Security     ApiKeyAuth
// @Tags         profile
// @Description  Обновляет изменяемые поля профиля текущего пользователя (например, email, cityID).
// @Id           update-profile
// @Accept       json
// @Produce      json
// @Param        input body updateUserProfileInput true "Обновляемые поля профиля"
// @Success      200 {object} models.UserProfile
// @Failure      400,401,500 {object} errorResponse
// @Router       /profile [patch]
func (h *Handler) updateProfile(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		c.Error(services.NewInternalServerError("failed to get user from context", err))
		return
	}

	var input updateUserProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	// Создаем модель на основе DTO
	userProfileUpdate := models.UserProfile{}
	if input.Email != nil {
		userProfileUpdate.Email.String = *input.Email
		userProfileUpdate.Email.Valid = true
	}
	if input.CityID != nil {
		userProfileUpdate.CityID = *input.CityID
	}

	updatedProfile, err := h.services.User.UpdateUserProfile(c.Request.Context(), userProfile.UserID, userProfileUpdate)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedProfile)
}

// @Summary      Обновить аватар пользователя
// @Security     ApiKeyAuth
// @Tags         profile
// @Description  Загружает новый файл аватара для текущего пользователя.
// @Id           update-avatar
// @Accept       multipart/form-data
// @Produce      json
// @Param        avatar formData file true "Файл изображения для аватара"
// @Success      200 {object} map[string]interface{} "message, avatarURL"
// @Failure      400,401,500 {object} errorResponse
// @Router       /profile/avatar [post]
func (h *Handler) updateAvatar(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		c.Error(services.NewInternalServerError("failed to get user from context", err))
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.Error(services.NewBadRequestError("avatar file is required", err))
		return
	}

	avatarKey, err := h.services.User.UpdateAvatar(c.Request.Context(), userProfile.UserID, file)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "avatar updated successfully",
		"avatarURL": avatarKey, // Возвращаем ключ объекта, а не полный URL
	})
}
