package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"lk/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить историю посещений
// @Security     ApiKeyAuth
// @Tags         medical-card
// @Description  Возвращает пагинированный список завершенных визитов.
// @Id           get-visits
// @Produce      json
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество на странице" default(10)
// @Success      200 {object} models.PaginatedVisitsResponse
// @Failure      401,500 {object} errorResponse
// @Router       /medical-card/visits [get]
func (h *Handler) getVisits(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	params := models.PaginationParams{Page: page, Limit: limit, SortBy: "date", SortOrder: "desc"}

	visits, err := h.services.MedicalCard.GetVisits(c.Request.Context(), userProfile.UserID, params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get visits: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, visits)
}

// @Summary      Получить список анализов
// @Security     ApiKeyAuth
// @Tags         medical-card
// @Description  Возвращает список анализов пользователя.
// @Id           get-analyses
// @Produce      json
// @Param        status query string false "Статус для фильтрации (e.g., 'completed', 'in_progress')"
// @Success      200 {array} models.LabAnalysis
// @Failure      401,500 {object} errorResponse
// @Router       /medical-card/analyses [get]
func (h *Handler) getAnalyses(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	status := c.Query("status")
	analyses, err := h.services.MedicalCard.GetAnalyses(c.Request.Context(), userProfile.UserID, &status)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get analyses: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, analyses)
}

// @Summary      Получить архив назначений
// @Security     ApiKeyAuth
// @Tags         medical-card
// @Description  Возвращает список выполненных или отмененных назначений.
// @Id           get-archived-prescriptions
// @Produce      json
// @Success      200 {array} models.Prescription
// @Failure      401,500 {object} errorResponse
// @Router       /medical-card/archive/prescriptions [get]
func (h *Handler) getArchivedPrescriptions(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	prescriptions, err := h.services.MedicalCard.GetArchivedPrescriptions(c.Request.Context(), userProfile.UserID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			"failed to get archived prescriptions: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, prescriptions)
}

// @Summary      Получить сводку по медкарте
// @Security     ApiKeyAuth
// @Tags         medical-card
// @Description  Возвращает сжатую сводку для главного экрана.
// @Id           get-summary
// @Produce      json
// @Success      200 {object} models.MedicalCardSummary
// @Failure      401,500 {object} errorResponse
// @Router       /medical-card/summary [get]
func (h *Handler) getSummary(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	summary, err := h.services.MedicalCard.GetSummary(c.Request.Context(), userProfile.UserID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get summary: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, summary)
}

// @Summary      Скачать файл
// @Security     ApiKeyAuth
// @Tags         medical-card
// @Description  Скачивает файл результата анализа или визита.
// @Id           download-file
// @Produce      application/octet-stream
// @Param        id path int true "ID файла"
// @Success      200 {file} file
// @Failure      401,404,500 {object} errorResponse
// @Router       /files/{id} [get]
func (h *Handler) downloadFile(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid file id")
		return
	}

	fileBytes, fileName, err := h.services.MedicalCard.DownloadFile(c.Request.Context(), userProfile.UserID, fileID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "file not found: "+err.Error())
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "application/octet-stream", fileBytes)
}

type archiveInput struct {
	PrescriptionID uint64 `json:"prescriptionId" binding:"required"`
}

// @Summary      Архивировать назначение из медкарты
// @Security     ApiKeyAuth
// @Tags         medical-card
// @Description  Перемещает активное назначение в архив.
// @Id           archive-prescription-from-card
// @Accept       json
// @Produce      json
// @Param        input body archiveInput true "ID Назначения"
// @Success      200 {object} statusResponse
// @Failure      400,401,404,500 {object} errorResponse
// @Router       /medical-card/archive/prescriptions [post]
func (h *Handler) archivePrescriptionFromCard(c *gin.Context) {
	userProfile, err := getUserProfile(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var input archiveInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input: "+err.Error())
		return
	}

	err = h.services.MedicalCard.ArchivePrescription(c.Request.Context(), userProfile.UserID,
		input.PrescriptionID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			"failed to archive prescription: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "prescription archived successfully"})
}
