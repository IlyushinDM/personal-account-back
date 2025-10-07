package http

import (
	"errors"
	"strings"

	"lk/internal/models"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

const adminCtx = "admin"

// adminIdentity - middleware, которое проверяет токен администратора
// и сохраняет его модель в контекст запроса.
func (h *Handler) adminIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.Error(services.NewUnauthorizedError("empty auth header", nil))
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.Error(services.NewUnauthorizedError("invalid auth header", nil))
		c.Abort()
		return
	}

	adminID, err := h.services.Admin.ParseAdminToken(headerParts[1])
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	admin, err := h.adminRepo.GetByID(c.Request.Context(), adminID)
	if err != nil {
		c.Error(services.NewUnauthorizedError("admin not found for this token", err))
		c.Abort()
		return
	}

	c.Set(adminCtx, admin)
}

// getAdmin - вспомогательная функция для извлечения модели администратора из контекста.
func getAdmin(c *gin.Context) (models.Admin, error) {
	adminVal, ok := c.Get(adminCtx)
	if !ok {
		return models.Admin{}, errors.New("admin not found in context")
	}
	admin, ok := adminVal.(models.Admin)
	if !ok {
		return models.Admin{}, errors.New("admin in context has invalid type")
	}
	return admin, nil
}
