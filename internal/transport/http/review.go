package http

import (
	"net/http"

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
