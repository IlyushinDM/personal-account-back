package handlers

import (
	"errors"
	"net/http"
	"strconv"

	_ "lk/internal/models"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить активные назначения
// @Security     ApiKeyAuth
// @Tags         prescriptions
// @Description  Возвращает список активных (неархивированных) назначений для текущего пользователя.
// @Id           get-active-prescriptions
// @Produce      json
// @Success      200 {array} models.Prescription
// @Failure      401 {object} errorResponse "Пользователь не авторизован"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /prescriptions/active [get]
func (h *Handler) getActivePrescriptions(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	prescriptions, err := h.services.Prescription.GetActiveForUser(c.Request.Context(), userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			"failed to get active prescriptions: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, prescriptions)
}

// @Summary      Архивировать назначение
// @Security     ApiKeyAuth
// @Tags         prescriptions
// @Description  Перемещает назначение в архив.
// @Id           archive-prescription
// @Produce      json
// @Param        id path int true "ID Назначения"
// @Success      200 {object} statusResponse
// @Failure      400 {object} errorResponse "Неверный формат ID"
// @Failure      403 {object} errorResponse "Нет прав на архивацию этого назначения"
// @Failure      404 {object} errorResponse "Назначение не найдено"
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /prescriptions/{id}/archive [post]
func (h *Handler) archivePrescription(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get user ID from context")
		return
	}

	idStr := c.Param("id")
	prescriptionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid prescription ID format")
		return
	}

	err = h.services.Prescription.ArchiveForUser(c.Request.Context(), userID, prescriptionID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPrescriptionNotFound):
			newErrorResponse(c, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrForbidden):
			newErrorResponse(c, http.StatusForbidden, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError,
				"failed to archive prescription: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "prescription archived successfully"})
}
