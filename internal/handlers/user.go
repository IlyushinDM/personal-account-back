package handlers

import (
	"database/sql"
	"net/http"

	"lk/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить профиль пользователя
// @Security     ApiKeyAuth
// @Tags         profile
// @Description  Возвращает полную информацию о профиле текущего пользователя и его предстоящие записи.
// @Id           get-profile
// @Produce      json
// @Success      200 {object} map[string]interface{} "profile, appointments"
// @Failure      401 {object} errorResponse "Пользователь не авторизован"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /profile [get]
func (h *Handler) getProfile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	profile, appointments, err := h.services.User.GetFullUserProfile(c.Request.Context(), userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile":      profile,
		"appointments": appointments,
	})
}

// updateUserProfileInput - DTO для обновления профиля.
// Используем указатели, чтобы различать отсутствующие поля и поля с нулевыми значениями.
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
// @Failure      400 {object} errorResponse "Ошибка валидации"
// @Failure      401 {object} errorResponse "Пользователь не авторизован"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /profile [patch]
func (h *Handler) updateProfile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	var input updateUserProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body: "+err.Error())
		return
	}

	// Создаем модель на основе DTO
	userProfileUpdate := models.UserProfile{}
	if input.Email != nil {
		userProfileUpdate.Email = sql.NullString{String: *input.Email, Valid: true}
	}
	if input.CityID != nil {
		userProfileUpdate.CityID = *input.CityID
	}

	updatedProfile, err := h.services.User.UpdateUserProfile(c.Request.Context(), userID, userProfileUpdate)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			"failed to update profile: "+err.Error())
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
// @Failure      400 {object} errorResponse "Файл не был предоставлен"
// @Failure      401 {object} errorResponse "Пользователь не авторизован"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /profile/avatar [post]
func (h *Handler) updateAvatar(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "avatar file is required: "+err.Error())
		return
	}

	avatarURL, err := h.services.User.UpdateAvatar(c.Request.Context(), userID, file)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to update avatar: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "avatar updated successfully",
		"avatarURL": avatarURL,
	})
}
