package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить дерево отделений
// @Security     ApiKeyAuth
// @Tags         directories
// @Description  Возвращает иерархический список всех отделений и вложенных в них специальностей.
// @Id           get-departments-tree
// @Produce      json
// @Success      200 {array} models.DepartmentWithSpecialties
// @Failure      500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router       /departments [get]
func (h *Handler) getDepartmentsTree(c *gin.Context) {
	departmentsTree, err := h.services.Directory.GetAllDepartmentsWithSpecialties(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,
			"Failed to get departments: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, departmentsTree)
}
