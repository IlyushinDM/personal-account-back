package handlers

import (
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
// @Description  Создает новый аккаунт пользователя и его профиль, возвращает пару токенов для автоматического входа.
// @Id           create-account
// @Accept       json
// @Produce      json
// @Param        input body signUpInput true "Информация для регистрации"
// @Success      201 {object} map[string]string "Возвращает accessToken и refreshToken"
// @Failure      400,409,500 {object} errorResponse
// @Router       /auth/register [post]
func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	tokens, err := h.services.Authorization.CreateUser(c.Request.Context(), input.Phone,
		input.Password, input.FullName, input.Gender, input.BirthDate, input.CityID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

// signInInput - структура для валидации данных при входе.
type signInInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Вход в систему (аутентификация)
// @Tags         auth
// @Description  Аутентифицирует пользователя по номеру телефона и паролю, возвращает пару JWT (access и refresh).
// @Id           login
// @Accept       json
// @Produce      json
// @Param        input body signInInput true "Учетные данные для входа"
// @Success      200 {object} map[string]string "Возвращает accessToken и refreshToken"
// @Failure      400,401,500 {object} errorResponse
// @Router       /auth/login [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	tokens, err := h.services.Authorization.GenerateToken(c.Request.Context(),
		input.Phone, input.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

type refreshInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// @Summary      Обновление токенов
// @Tags         auth
// @Description  Получает новую пару access и refresh токенов по валидному refresh токену.
// @Id           refresh-token
// @Accept       json
// @Produce      json
// @Param        input body refreshInput true "Refresh токен"
// @Success      200 {object} map[string]string "Возвращает accessToken и refreshToken"
// @Failure      400,401 {object} errorResponse
// @Router       /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var input refreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	tokens, err := h.services.Authorization.RefreshToken(c.Request.Context(), input.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary      Выход из системы
// @Tags         auth
// @Description  Инвалидирует refresh токен на сервере, завершая сессию.
// @Id           logout
// @Accept       json
// @Produce      json
// @Param        input body refreshInput true "Refresh токен"
// @Success      200 {object} statusResponse "Статус операции"
// @Failure      400,401,500 {object} errorResponse
// @Router       /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	var input refreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	if err := h.services.Authorization.Logout(c.Request.Context(), input.RefreshToken); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "you have been logged out"})
}

type forgotPasswordInput struct {
	Phone string `json:"phone" binding:"required"`
}

// @Summary      Восстановление пароля (шаг 1: запрос кода)
// @Tags         auth
// @Description  Отправляет код подтверждения для сброса пароля (в текущей реализации код выводится в лог сервера).
// @Id           forgot-password
// @Accept       json
// @Produce      json
// @Param        input body forgotPasswordInput true "Номер телефона"
// @Success      200 {object} statusResponse
// @Failure      400,500 {object} errorResponse
// @Router       /auth/forgot-password [post]
func (h *Handler) forgotPassword(c *gin.Context) {
	var input forgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	if err := h.services.Authorization.ForgotPassword(c.Request.Context(), input.Phone); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "confirmation code has been sent"})
}

type resetPasswordInput struct {
	Phone       string `json:"phone" binding:"required"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// @Summary      Восстановление пароля (шаг 2: сброс)
// @Tags         auth
// @Description  Устанавливает новый пароль, используя код подтверждения.
// @Id           reset-password
// @Accept       json
// @Produce      json
// @Param        input body resetPasswordInput true "Данные для сброса"
// @Success      200 {object} statusResponse
// @Failure      400,401,500 {object} errorResponse
// @Router       /auth/reset-password [post]
func (h *Handler) resetPassword(c *gin.Context) {
	var input resetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	err := h.services.Authorization.ResetPassword(c.Request.Context(), input.Phone, input.Code,
		input.NewPassword)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "password has been reset successfully"})
}

// @Summary      Вход через Госуслуги (заглушка)
// @Tags         auth
// @Description  Перенаправляет пользователя на страницу авторизации Госуслуг. (Не реализовано)
// @Id           gosuslugi-login
// @Router       /auth/gosuslugi [get]
func (h *Handler) gosuslugiLogin(c *gin.Context) {
	// TODO: FR-1.3 - Реализовать OAuth2-аутентификацию через Госуслуги.
	// 1. Создать OAuth2-конфигурацию с ClientID, ClientSecret, RedirectURL и эндпоинтами Госуслуг.
	// 2. Сгенерировать URL для редиректа пользователя на страницу авторизации Госуслуг.
	// 3. Временный state-параметр для защиты от CSRF-атак нужно сохранить в сессии/кэше.
	// 4. Выполнить редирект: c.Redirect(http.StatusTemporaryRedirect, url)
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

// @Summary      Callback от Госуслуг (заглушка)
// @Tags         auth
// @Description  Обрабатывает ответ от Госуслуг. (Не реализовано)
// @Id           gosuslugi-callback
// @Router       /auth/gosuslugi/callback [post]
func (h *Handler) gosuslugiCallback(c *gin.Context) {
	// TODO: FR-1.4 - Реализовать обработку callback от Госуслуг.
	// 1. Получить 'code' и 'state' из query-параметров.
	// 2. Проверить 'state' на соответствие сохраненному для защиты от CSRF.
	// 3. Обменять 'code' на access-токен Госуслуг.
	// 4. С помощью токена получить данные пользователя (СНИЛС, ФИО и т.д.) из API Госуслуг.
	// 5. Вызвать новый метод в Authorization-сервисе, например, `AuthorizeGosuslugi(user_data)`.
	// 6. Сервис должен найти пользователя по СНИЛС, или создать нового, если он не найден.
	// 7. В случае успеха сгенерировать пару JWT-токенов (access/refresh) и вернуть их клиенту.
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}
