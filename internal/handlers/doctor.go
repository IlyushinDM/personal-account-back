package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// findSpecialistsQuery - структура для валидации query-параметров при поиске врачей.
type findSpecialistsQuery struct {
	// Поиск по ФИО или названию специальности
	Query string `form:"q"`
	// Фильтрация по ID специальности
	SpecialtyID uint32 `form:"specialtyID"`
}

// @Summary Find Specialists
// @Security ApiKeyAuth
// @Tags specialists
// @Description Find specialists by query string (name, specialty) or filter by specialty ID
// @Id find-specialists
// @Produce  json
// @Param q query string false "Search query (full name or specialty name)"
// @Param specialtyID query int false "ID of the specialty to filter by"
// @Success 200 {object} []models.Doctor
// @Failure 400,500 {object} errorResponse
// @Router /specialists [get]
func (h *Handler) findSpecialists(c *gin.Context) {
	var queryParams findSpecialistsQuery

	// Используем BindQuery для парсинга параметров из URL.
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid query parameters: "+err.Error())
		return
	}

	ctx := c.Request.Context()

	// Бизнес-логика выбора метода сервиса
	if queryParams.Query != "" {
		// Если есть поисковый запрос, ищем по нему.
		doctors, err := h.services.Doctor.SearchDoctors(ctx, queryParams.Query)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError,
				"Failed to search doctors: "+err.Error())
			return
		}
		c.JSON(http.StatusOK, doctors)
		return
	}

	if queryParams.SpecialtyID > 0 {
		// Если есть ID специальности, фильтруем по нему.
		doctors, err := h.services.Doctor.GetDoctorsBySpecialty(ctx, queryParams.SpecialtyID)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError,
				"Failed to get doctors by specialty: "+err.Error())
			return
		}
		c.JSON(http.StatusOK, doctors)
		return
	}

	newErrorResponse(c, http.StatusBadRequest,
		"At least one query parameter ('q' or 'specialtyID') is required")
}

// @Summary Get Specialist By ID
// @Security ApiKeyAuth
// @Tags specialists
// @Description Get detailed information about a specialist by their ID
// @Id get-specialist-by-id
// @Produce  json
// @Param id path int true "Specialist ID"
// @Success 200 {object} models.Doctor
// @Failure 400,404,500 {object} errorResponse
// @Router /specialists/{id} [get]
func (h *Handler) getSpecialistByID(c *gin.Context) {
	// Получаем ID из URL-параметра.
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid specialist ID format")
		return
	}

	doctor, err := h.services.Doctor.GetDoctorByID(c.Request.Context(), id)
	if err != nil {
		// TODO: Различать ошибку "не найдено" (404) от других серверных ошибок (500).
		// Это требует от сервиса возвращать кастомные типы ошибок.
		newErrorResponse(c, http.StatusInternalServerError,
			"Failed to get doctor details: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, doctor)
}
