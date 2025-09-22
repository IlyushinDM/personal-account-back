package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"lk/internal/models"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

// createAppointmentInput - структура для валидации JSON-тела при создании записи.
type createAppointmentInput struct {
	DoctorID        uint64    `json:"doctorID" binding:"required"`
	ServiceID       uint64    `json:"serviceID" binding:"required"`
	ClinicID        uint64    `json:"clinicID" binding:"required"`
	AppointmentDate time.Time `json:"appointmentDate" binding:"required"` // Формат: "2025-09-15T10:00:00Z"
	AppointmentTime string    `json:"appointmentTime" binding:"required"` // Формат: "10:00"
	PriceAtBooking  float64   `json:"priceAtBooking" binding:"required"`
	IsDMS           bool      `json:"isDms"`
}

// @Summary      Создать запись на приём
// @Security     ApiKeyAuth
// @Tags         appointments
// @Description  Создает новую запись на прием для текущего пользователя.
// @ID           create-appointment
// @Accept       json
// @Produce      json
// @Param        input body createAppointmentInput true "Информация о записи"
// @Success      201 {object} map[string]interface{} "message, appointmentID"
// @Failure      400 {object} errorResponse "Неверные входные данные"
// @Failure      401 {object} errorResponse "Пользователь не авторизован"
// @Failure      404 {object} errorResponse "Врач с указанным ID не найден"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /appointments [post]
func (h *Handler) createAppointment(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	var input createAppointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input body: "+err.Error())
		return
	}

	appointment := models.Appointment{
		UserID:          userProfile.UserID,
		DoctorID:        input.DoctorID,
		ServiceID:       input.ServiceID,
		ClinicID:        input.ClinicID,
		AppointmentDate: input.AppointmentDate,
		AppointmentTime: input.AppointmentTime,
		PriceAtBooking:  input.PriceAtBooking,
		IsDMS:           input.IsDMS,
		StatusID:        1, // Например, что 1 - это статус "Запланировано"
	}

	appointmentID, err := h.services.Appointment.CreateAppointment(c.Request.Context(), appointment)
	if err != nil {
		if errors.Is(err, services.ErrDoctorNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "Failed to create appointment: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Appointment created successfully",
		"appointmentID": appointmentID,
	})
}

// @Summary      Получить записи на приём пользователя
// @Security     ApiKeyAuth
// @Tags         appointments
// @Description  Получает список всех записей на прием для текущего пользователя.
// @ID           get-user-appointments
// @Produce      json
// @Success      200 {array} models.Appointment
// @Failure      401 {object} errorResponse "Пользователь не авторизован"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /appointments [get]
func (h *Handler) getUserAppointments(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	appointments, err := h.services.Appointment.GetUserAppointments(c.Request.Context(), userProfile.UserID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to get appointments: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// @Summary      Отменить запись на приём
// @Security     ApiKeyAuth
// @Tags         appointments
// @Description  Отменяет существующую запись на прием по ее ID.
// @ID           cancel-appointment
// @Produce      json
// @Param        id path int true "ID Записи"
// @Success      200 {object} statusResponse "Статус операции"
// @Failure      400 {object} errorResponse "Неверный формат ID"
// @Failure      401 {object} errorResponse "Пользователь не авторизован"
// @Failure      403 {object} errorResponse "Нет прав на отмену этой записи"
// @Failure      404 {object} errorResponse "Запись не найдена"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /appointments/{id} [delete]
func (h *Handler) cancelAppointment(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	idStr := c.Param("id")
	appointmentID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid appointment ID format")
		return
	}

	err = h.services.Appointment.CancelAppointment(c.Request.Context(), userProfile.UserID, appointmentID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAppointmentNotFound):
			newErrorResponse(c, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrForbidden):
			newErrorResponse(c, http.StatusForbidden, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, "Failed to cancel appointment: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "appointment cancelled successfully"})
}

// availableDatesQuery - структура для валидации query-параметров.
type availableDatesQuery struct {
	SpecialistID uint64 `form:"specialistId" binding:"required"`
	ServiceID    uint64 `form:"serviceId" binding:"required"`
	Month        string `form:"month" binding:"required"` // Формат: YYYY-MM
}

// @Summary      Получить доступные дни для записи
// @Security     ApiKeyAuth
// @Tags         appointments
// @Description  Возвращает список дней в указанном месяце, в которые у специалиста есть свободные слоты.
// @Id           get-available-dates
// @Produce      json
// @Param        specialistId query int true "ID Специалиста"
// @Param        serviceId query int true "ID Услуги"
// @Param        month query string true "Месяц в формате YYYY-MM"
// @Success      200 {object} models.AvailableDatesResponse
// @Failure      400 {object} errorResponse "Неверные параметры запроса"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /appointments/available-dates [get]
func (h *Handler) getAvailableDates(c *gin.Context) {
	var queryParams availableDatesQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid query parameters: "+err.Error())
		return
	}

	dates, err := h.services.Appointment.GetAvailableDates(c.Request.Context(),
		queryParams.SpecialistID, queryParams.ServiceID, queryParams.Month)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to get available dates: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, dates)
}

// availableSlotsQuery - структура для валидации query-параметров.
type availableSlotsQuery struct {
	SpecialistID uint64 `form:"specialistId" binding:"required"`
	ServiceID    uint64 `form:"serviceId" binding:"required"`
	Date         string `form:"date" binding:"required"` // Формат: YYYY-MM-DD
}

// @Summary      Получить доступные слоты времени
// @Security     ApiKeyAuth
// @Tags         appointments
// @Description  Возвращает список свободных временных слотов у специалиста на указанную дату.
// @Id           get-available-slots
// @Produce      json
// @Param        specialistId query int true "ID Специалиста"
// @Param        serviceId query int true "ID Услуги"
// @Param        date query string true "Дата в формате YYYY-MM-DD"
// @Success      200 {object} models.AvailableSlotsResponse
// @Failure      400 {object} errorResponse "Неверные параметры запроса"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /appointments/available-slots [get]
func (h *Handler) getAvailableSlots(c *gin.Context) {
	var queryParams availableSlotsQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid query parameters: "+err.Error())
		return
	}

	slots, err := h.services.Appointment.GetAvailableSlots(c.Request.Context(),
		queryParams.SpecialistID, queryParams.ServiceID, queryParams.Date)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to get available slots: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, slots)
}

// @Summary      Получить предстоящие записи пользователя
// @Security     ApiKeyAuth
// @Tags         appointments
// @Description  Получает отсортированный список предстоящих записей для текущего пользователя.
// @Id           get-upcoming-appointments
// @Produce      json
// @Success      200 {array} models.Appointment
// @Failure      401 {object} errorResponse "Пользователь не авторизован"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /appointments/upcoming [get]
func (h *Handler) getUpcomingAppointments(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	appointments, err := h.services.Appointment.GetUpcomingForUser(c.Request.Context(), userProfile.UserID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to get upcoming appointments: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, appointments)
}
