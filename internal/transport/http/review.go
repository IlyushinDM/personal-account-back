package http

import (
	"database/sql"
	"net/http"
	"strconv"

	"lk/internal/models"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

// getReviewsQuery - структура для валидации query-параметров при получении отзывов.
type getReviewsQuery struct {
	DoctorID  uint64 `form:"doctor_id" binding:"required"`
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=10"`
	SortBy    string `form:"sort,default=rating"`
	SortOrder string `form:"order,default=desc"`
}

// @Summary      Получить отзывы врача
// @Security     ApiKeyAuth
// @Tags         reviews
// @Description  Получает пагинированный список модерированных отзывов для конкретного врача с возможностью сортировки.
// @Id           get-reviews
// @Produce      json
// @Param        doctor_id query int true "ID врача"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество элементов на странице (максимум 100)" default(10)
// @Param        sort query string false "Поле для сортировки: 'rating', 'created_at', 'date'" Enums(rating, created_at, date) default(rating)
// @Param        order query string false "Порядок сортировки: 'asc', 'desc'" Enums(asc, desc) default(desc)
// @Success      200 {object} models.PaginatedReviewsResponse "Список отзывов с пагинацией"
// @Failure      400,401,500 {object} errorResponse
// @Router       /reviews [get]
func (h *Handler) getReviews(c *gin.Context) {
	var queryParams getReviewsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.Error(services.NewBadRequestError("Некорректные параметры запроса", err))
		return
	}

	// Валидация параметров
	if queryParams.DoctorID == 0 {
		c.Error(services.NewBadRequestError("ID врача обязателен", nil))
		return
	}

	if queryParams.Page < 1 {
		queryParams.Page = 1
	}

	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		queryParams.Limit = 10
	}

	// Валидация параметров сортировки
	allowedSortFields := map[string]bool{
		"rating":     true,
		"created_at": true,
		"date":       true,
	}
	if !allowedSortFields[queryParams.SortBy] {
		queryParams.SortBy = "rating" // Значение по умолчанию
	}

	allowedSortOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}
	if !allowedSortOrders[queryParams.SortOrder] {
		queryParams.SortOrder = "desc" // Значение по умолчанию
	}

	// Создаем параметры пагинации
	paginationParams := models.PaginationParams{
		Page:      queryParams.Page,
		Limit:     queryParams.Limit,
		SortBy:    queryParams.SortBy,
		SortOrder: queryParams.SortOrder,
	}

	ctx := c.Request.Context()

	// Получаем отзывы через сервис
	reviewsResponse, err := h.services.Review.GetReviewsByDoctorID(ctx, queryParams.DoctorID, paginationParams)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reviewsResponse)
}

// @Summary      Получить детали отзыва
// @Security     ApiKeyAuth
// @Tags         reviews
// @Description  Получает детальную информацию о конкретном отзыве по его ID. Возвращает только модерированные отзывы.
// @Id           get-review-by-id
// @Produce      json
// @Param        id path int true "ID отзыва"
// @Success      200 {object} models.Review "Детали отзыва"
// @Failure      400,401,404,500 {object} errorResponse
// @Router       /reviews/{id} [get]
func (h *Handler) getReviewByID(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("Некорректный формат ID отзыва", err))
		return
	}

	ctx := c.Request.Context()

	// Получаем отзыв через сервис
	review, err := h.services.Review.GetReviewByID(ctx, reviewID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, review)
}

// @Summary      Создать новый отзыв
// @Security     ApiKeyAuth
// @Tags         reviews
// @Description  Создает новый отзыв о враче. Пользователь может оставить только один отзыв на врача.
// @Id           create-review
// @Accept       json
// @Produce      json
// @Param        request body models.CreateReviewRequest true "Данные отзыва"
// @Success      201 {object} models.CreateReviewResponse "Отзыв создан"
// @Failure      400,401,409,500 {object} errorResponse
// @Router       /reviews [post]
func (h *Handler) createReview(c *gin.Context) {
	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(services.NewBadRequestError("Некорректные данные отзыва", err))
		return
	}

	// Получаем ID пользователя из контекста (устанавливается middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(services.NewUnauthorizedError("Пользователь не авторизован", nil))
		return
	}

	// Создаем модель отзыва
	review := models.Review{
		DoctorID: req.DoctorID,
		Rating:   req.Rating,
		Comment:  sql.NullString{String: req.Comment, Valid: req.Comment != ""},
	}

	ctx := c.Request.Context()

	// Создаем отзыв через сервис
	reviewID, err := h.services.Review.CreateReview(ctx, userID.(uint64), review)
	if err != nil {
		c.Error(err)
		return
	}

	response := models.CreateReviewResponse{
		ID:      reviewID,
		Message: "Отзыв создан",
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary      Обновить отзыв
// @Security     ApiKeyAuth
// @Tags         reviews
// @Description  Обновляет существующий отзыв. Пользователь может обновлять только свои отзывы.
// @Id           update-review
// @Accept       json
// @Produce      json
// @Param        id path int true "ID отзыва"
// @Param        request body models.UpdateReviewRequest true "Обновленные данные отзыва"
// @Success      200 {object} models.UpdateReviewResponse "Отзыв обновлен"
// @Failure      400,401,403,404,500 {object} errorResponse
// @Router       /reviews/{id} [put]
func (h *Handler) updateReview(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("Некорректный формат ID отзыва", err))
		return
	}

	var req models.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(services.NewBadRequestError("Некорректные данные отзыва", err))
		return
	}

	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(services.NewUnauthorizedError("Пользователь не авторизован", nil))
		return
	}

	// Создаем модель отзыва для обновления
	review := models.Review{
		Rating:  req.Rating,
		Comment: sql.NullString{String: req.Comment, Valid: req.Comment != ""},
	}

	ctx := c.Request.Context()

	// Обновляем отзыв через сервис
	err = h.services.Review.UpdateReview(ctx, userID.(uint64), reviewID, review)
	if err != nil {
		c.Error(err)
		return
	}

	response := models.UpdateReviewResponse{
		Message: "Отзыв обновлен",
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Удалить отзыв
// @Security     ApiKeyAuth
// @Tags         reviews
// @Description  Удаляет отзыв. Пользователь может удалять только свои отзывы.
// @Id           delete-review
// @Produce      json
// @Param        id path int true "ID отзыва"
// @Success      200 {object} models.DeleteReviewResponse "Отзыв удален"
// @Failure      400,401,403,404,500 {object} errorResponse
// @Router       /reviews/{id} [delete]
func (h *Handler) deleteReview(c *gin.Context) {
	idStr := c.Param("id")
	reviewID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("Некорректный формат ID отзыва", err))
		return
	}

	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(services.NewUnauthorizedError("Пользователь не авторизован", nil))
		return
	}

	ctx := c.Request.Context()

	// Удаляем отзыв через сервис
	err = h.services.Review.DeleteReview(ctx, userID.(uint64), reviewID)
	if err != nil {
		c.Error(err)
		return
	}

	response := models.DeleteReviewResponse{
		Message: "Отзыв удален",
	}

	c.JSON(http.StatusOK, response)
}

// getDoctorReviewsQuery - структура для валидации query-параметров при получении отзывов врача.
type getDoctorReviewsQuery struct {
	Page          int  `form:"page,default=1"`
	Limit         int  `form:"limit,default=5"`
	OnlyModerated bool `form:"only_moderated,default=true"`
}

// @Summary      Получить отзывы врача
// @Security     ApiKeyAuth
// @Tags         doctors
// @Description  Получает отзывы конкретного врача с его средним рейтингом.
// @Id           get-doctor-reviews
// @Produce      json
// @Param        id path int true "ID врача"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество элементов на странице (максимум 100)" default(5)
// @Param        only_moderated query bool false "Показывать только модерированные отзывы" default(true)
// @Success      200 {object} models.DoctorReviewsResponse "Отзывы врача с рейтингом"
// @Failure      400,401,500 {object} errorResponse
// @Router       /doctors/{id}/reviews [get]
func (h *Handler) getDoctorReviews(c *gin.Context) {
	idStr := c.Param("id")
	doctorID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("Некорректный формат ID врача", err))
		return
	}

	var queryParams getDoctorReviewsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.Error(services.NewBadRequestError("Некорректные параметры запроса", err))
		return
	}

	// Валидация параметров
	if queryParams.Page < 1 {
		queryParams.Page = 1
	}
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		queryParams.Limit = 5
	}

	// Создаем параметры пагинации
	paginationParams := models.PaginationParams{
		Page:      queryParams.Page,
		Limit:     queryParams.Limit,
		SortBy:    "rating", // Сортировка по рейтингу по умолчанию
		SortOrder: "desc",
	}

	ctx := c.Request.Context()

	// Получаем отзывы через сервис
	reviewsResponse, err := h.services.Review.GetDoctorReviewsWithRating(
		ctx, doctorID, paginationParams, queryParams.OnlyModerated)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reviewsResponse)
}
