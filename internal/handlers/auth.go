package handlers

import (
	"errors"
	"net/http"

	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

// signUpInput - структура для валидации входящих данных при регистрации.
type signUpInput struct {
	Phone     string `json:"phone" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	FullName  string `json:"fullName" binding:"required"`
	Gender    string `json:"gender" binding:"required"`
	BirthDate string `json:"birthDate" binding:"required"` // Ожидаемый формат: "YYYY-MM-DD"
	CityID    uint32 `json:"cityID" binding:"required"`
}

// @Summary      Регистрация пользователя
// @Tags         auth
// @Description  Создает новый аккаунт пользователя и его профиль, возвращает токен для автоматического входа.
// @Id           create-account
// @Accept       json
// @Produce      json
// @Param        input body signUpInput true "Информация для регистрации"
// @Success      201 {object} map[string]interface{} "message, userID, accessToken, tokenType"
// @Failure      400,409 {object} errorResponse "Ошибка валидации или пользователь уже существует"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body: "+err.Error())
		return
	}

	// Вызываем сервис для создания пользователя
	userID, err := h.services.Authorization.CreateUser(c.Request.Context(), input.Phone,
		input.Password,
		input.FullName, input.Gender, input.BirthDate, input.CityID)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Генерируем токен для автоматического входа после регистрации
	token, err := h.services.Authorization.GenerateToken(c.Request.Context(), input.Phone,
		input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			"failed to generate token after registration: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "user successfully registered",
		"userID":      userID,
		"accessToken": token,
		"tokenType":   "Bearer",
	})
}

// signInInput - структура для валидации данных при входе.
type signInInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Вход в систему (аутентификация)
// @Tags         auth
// @Description  Аутентифицирует пользователя по номеру телефона и паролю, возвращает JWT.
// @Id           login
// @Accept       json
// @Produce      json
// @Param        input body signInInput true "Учетные данные для входа"
// @Success      200 {object} map[string]interface{} "accessToken, tokenType"
// @Failure      400,401 {object} errorResponse "Ошибка валидации или неверные учетные данные"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /auth/login [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body: "+err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(c.Request.Context(),
		input.Phone, input.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
		"tokenType":   "Bearer",
	})
}
