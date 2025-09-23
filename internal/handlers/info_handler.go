package handlers

import (
	"net/http"

	_ "lk/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить информацию о клинике
// @Security     ApiKeyAuth
// @Tags         info
// @Description  Возвращает контактную информацию, адреса и часы работы клиники.
// @Id           get-clinic-info
// @Produce      json
// @Success      200 {object} models.ClinicInfo
// @Failure      401,404,500 {object} errorResponse
// @Router       /clinic-info [get]
func (h *Handler) getClinicInfo(c *gin.Context) {
	info, err := h.services.Info.GetClinicInfo(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, info)
}

// @Summary      Получить юридические документы
// @Security     ApiKeyAuth
// @Tags         info
// @Description  Возвращает список юридических документов с ссылками и версиями.
// @Id           get-legal-documents
// @Produce      json
// @Success      200 {array} models.LegalDocument
// @Failure      401,500 {object} errorResponse
// @Router       /legal/documents [get]
func (h *Handler) getLegalDocuments(c *gin.Context) {
	docs, err := h.services.Info.GetLegalDocuments(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, docs)
}
