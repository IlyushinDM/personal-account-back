package handlers

import (
	"net/http"
	"strconv"
	"time"

	"lk/internal/models"

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

// @Summary Create Appointment
// @Security ApiKeyAuth
// @Tags appointments
// @Description Create a new appointment for the current user
// @ID create-appointment
// @Accept  json
// @Produce  json
// @Param input body createAppointmentInput true "Appointment info"
// @Success 201 {object} map[string]interface{}
// @Failure 400,401,500 {object} errorResponse
// @Router /appointments [post]
func (h *Handler) createAppointment(c *gin.Context) {
	userID, err := getUserID(c)
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
		UserID:          userID,
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
		newErrorResponse(c, http.StatusInternalServerError, "Failed to create appointment: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Appointment created successfully",
		"appointmentID": appointmentID,
	})
}

// @Summary Get User Appointments
// @Security ApiKeyAuth
// @Tags appointments
// @Description Get a list of all appointments for the current user
// @ID get-user-appointments
// @Produce  json
// @Success 200 {object} []models.Appointment
// @Failure 401,500 {object} errorResponse
// @Router /appointments [get]
func (h *Handler) getUserAppointments(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	appointments, err := h.services.Appointment.GetUserAppointments(c.Request.Context(), userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to get appointments: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// @Summary Cancel Appointment
// @Security ApiKeyAuth
// @Tags appointments
// @Description Cancel an existing appointment
// @ID cancel-appointment
// @Produce  json
// @Param id path int true "Appointment ID"
// @Success 200 {object} statusResponse
// @Failure 400,401,403,500 {object} errorResponse
// @Router /appointments/{id} [delete]
func (h *Handler) cancelAppointment(c *gin.Context) {
	_, err := getUserID(c) // Проверяем авторизацию
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	idStr := c.Param("id")
	_, err = strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid appointment ID format")
		return
	}

	// TODO: Реализовать логику отмены в сервисе и репозитории.
	// 1. Сервис должен проверить, принадлежит ли эта запись текущему пользователю.
	//    Это защитит от ситуации, когда пользователь A пытается отменить запись пользователя B.
	// 2. Репозиторий должен обновить статус записи на "Отменено".

	// Заглушка
	// err = h.services.Appointment.Cancel(c.Request.Context(), userID, appointmentID)

	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented yet"})
}
