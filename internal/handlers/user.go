package handlers

import (
	"net/http"

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

// @Summary      Обновить профиль пользователя
// @Security     ApiKeyAuth
// @Tags         profile
// @Description  Обновляет изменяемые поля профиля текущего пользователя.
// @Id           update-profile
// @Accept       json
// @Produce      json
// @Param        input body map[string]interface{} true "Обновляемые поля профиля"
// @Success      200 {object} statusResponse "Статус операции"
// @Failure      400,401 {object} errorResponse "Ошибка валидации или пользователь не авторизован"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /profile [patch]
func (h *Handler) updateProfile(c *gin.Context) {
	// 1. Получить userID из контекста
	// 2. Считать JSON с обновляемыми полями
	// 3. Вызвать соответствующий метод сервиса
	// 4. Вернуть обновленный профиль или статус OK
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented yet"})
}
