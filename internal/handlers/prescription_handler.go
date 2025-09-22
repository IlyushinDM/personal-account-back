package handlers

import (
	"net/http"
	"strconv"

	"lk/internal/services"

	_ "lk/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить активные назначения
// @Security     ApiKeyAuth
// @Tags         prescriptions
// @Description  Возвращает список активных (неархивированных) назначений для текущего пользователя.
// @Id           get-active-prescriptions
// @Produce      json
// @Success      200 {array} models.Prescription
// @Failure      401,500 {object} errorResponse
// @Router       /prescriptions/active [get]
func (h *Handler) getActivePrescriptions(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		c.Error(services.NewInternalServerError("failed to get user from context", err))
		return
	}

	prescriptions, err := h.services.Prescription.GetActiveForUser(c.Request.Context(), userProfile.UserID)
	if err != nil {
		c.Error(err)
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
// @Failure      400,401,403,404,500 {object} errorResponse
// @Router       /prescriptions/{id}/archive [post]
func (h *Handler) archivePrescription(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		c.Error(services.NewInternalServerError("failed to get user from context", err))
		return
	}

	idStr := c.Param("id")
	prescriptionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid prescription ID format", err))
		return
	}

	err = h.services.Prescription.ArchiveForUser(c.Request.Context(), userProfile.UserID, prescriptionID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "prescription archived successfully"})
}
