package handlers

import (
	"net/http"
	"strconv"

	"lk/internal/models"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

// findSpecialistsQuery - структура для валидации query-параметров при поиске и фильтрации врачей.
type findSpecialistsQuery struct {
	Query       string `form:"q"`
	Service     string `form:"service"`
	SpecialtyID uint32 `form:"specialtyID"`
	Page        int    `form:"page,default=1"`
	Limit       int    `form:"limit,default=10"`
	SortBy      string `form:"sortBy,default=rating"`
	SortOrder   string `form:"sortOrder,default=desc"`
}

// @Summary      Найти или отфильтровать специалистов
// @Security     ApiKeyAuth
// @Tags         specialists
// @Description  Находит специалистов по различным критериям. Используйте ОДИН из следующих параметров: 'q' для поиска по ФИО/специальности, 'service' для поиска по услуге, или 'specialtyID' для фильтрации с пагинацией и сортировкой.
// @Id           find-specialists
// @Produce      json
// @Param        q query string false "Поисковый запрос (ФИО или название специальности)"
// @Param        service query string false "Название медицинской услуги"
// @Param        specialtyID query int false "ID специальности для фильтрации"
// @Param        page query int false "Номер страницы (используется с 'specialtyID')" default(1)
// @Param        limit query int false "Количество элементов на странице (используется с 'specialtyID')" default(10)
// @Param        sortBy query string false "Поле для сортировки: 'rating', 'experience', 'name' (используется с 'specialtyID')" Enums(rating, experience, name) default(rating)
// @Param        sortOrder query string false "Порядок сортировки: 'asc', 'desc' (используется с 'specialtyID')" Enums(asc, desc) default(desc)
// @Success      200 {object} models.PaginatedDoctorsResponse "Возвращается при использовании 'specialtyID'"
// @Success      200 {array} models.Doctor "Возвращается при использовании 'q' или 'service'"
// @Failure      400,401,500 {object} errorResponse
// @Router       /specialists [get]
func (h *Handler) findSpecialists(c *gin.Context) {
	var queryParams findSpecialistsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.Error(services.NewBadRequestError("Invalid query parameters", err))
		return
	}

	ctx := c.Request.Context()

	// Логика приоритетов: q > service > specialtyID
	if queryParams.Query != "" {
		doctors, err := h.services.Doctor.SearchDoctors(ctx, queryParams.Query)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, doctors)
		return
	}

	if queryParams.Service != "" {
		doctors, err := h.services.Doctor.SearchDoctorsByService(ctx, queryParams.Service)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, doctors)
		return
	}

	if queryParams.SpecialtyID > 0 {
		paginationParams := models.PaginationParams{
			Page:      queryParams.Page,
			Limit:     queryParams.Limit,
			SortBy:    queryParams.SortBy,
			SortOrder: queryParams.SortOrder,
		}
		paginatedResponse, err := h.services.Doctor.GetDoctorsBySpecialty(ctx, queryParams.SpecialtyID, paginationParams)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, paginatedResponse)
		return
	}

	c.Error(services.NewBadRequestError("At least one query parameter ('q', 'service', or 'specialtyID') is required", nil))
}

// @Summary      Получить специалиста по ID
// @Security     ApiKeyAuth
// @Tags         specialists
// @Description  Получает детальную информацию о специалисте по его ID.
// @Id           get-specialist-by-id
// @Produce      json
// @Param        id path int true "ID Специалиста"
// @Success      200 {object} models.Doctor
// @Failure      400,401,404,500 {object} errorResponse
// @Router       /specialists/{id} [get]
func (h *Handler) getSpecialistByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("Invalid specialist ID format", err))
		return
	}

	doctor, err := h.services.Doctor.GetDoctorByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, doctor)
}

// @Summary      Получить рекомендации специалиста
// @Security     ApiKeyAuth
// @Tags         specialists
// @Description  Получает общие рекомендации по подготовке к приему у конкретного специалиста.
// @Id           get-specialist-recommendations
// @Produce      json
// @Param        id path int true "ID Специалиста"
// @Success      200 {object} models.Recommendation
// @Failure      400,401,404,500 {object} errorResponse
// @Router       /specialists/{id}/recommendations [get]
func (h *Handler) getSpecialistRecommendations(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("Invalid specialist ID format", err))
		return
	}

	recommendations, err := h.services.Doctor.GetSpecialistRecommendations(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// @Summary      Получить рекомендации для услуги
// @Security     ApiKeyAuth
// @Tags         services
// @Description  Получает общие рекомендации по подготовке к конкретной медицинской услуге.
// @Id           get-service-recommendations
// @Produce      json
// @Param        id path int true "ID Услуги"
// @Success      200 {object} models.Recommendation
// @Failure      400,401,500 {object} errorResponse
// @Router       /services/{id}/recommendations [get]
func (h *Handler) getServiceRecommendations(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("Invalid service ID format", err))
		return
	}

	recommendations, err := h.services.Info.GetServiceRecommendations(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, recommendations)
}
