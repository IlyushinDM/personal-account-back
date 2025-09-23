package http

import (
	"net/http"
	"strconv"

	"lk/internal/services"

	_ "lk/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary      Получить дерево отделений
// @Security     ApiKeyAuth
// @Tags         directories
// @Description  Возвращает иерархический список всех отделений и вложенных в них специальностей.
// @Id           get-departments-tree
// @Produce      json
// @Success      200 {array} models.DepartmentWithSpecialties
// @Failure      401,500 {object} errorResponse
// @Router       /departments [get]
func (h *Handler) getDepartmentsTree(c *gin.Context) {
	departmentsTree, err := h.services.Directory.GetAllDepartmentsWithSpecialties(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, departmentsTree)
}

// @Summary      Получить список специальностей
// @Security     ApiKeyAuth
// @Tags         directories
// @Description  Возвращает список всех врачебных специальностей. Можно отфильтровать по ID отделения.
// @Id           get-specialties
// @Produce      json
// @Param        departmentID query int false "ID Отделения для фильтрации"
// @Success      200 {array} models.Specialty
// @Failure      400,401,500 {object} errorResponse
// @Router       /specialties [get]
func (h *Handler) getSpecialties(c *gin.Context) {
	var departmentID *uint32
	deptIDStr := c.Query("departmentID")

	// Если параметр был передан, парсим его
	if deptIDStr != "" {
		id, err := strconv.ParseUint(deptIDStr, 10, 32)
		if err != nil {
			c.Error(services.NewBadRequestError("Invalid department ID format", err))
			return
		}
		id32 := uint32(id)
		departmentID = &id32
	}

	// Вызываем сервис, передавая nil, если параметр не был указан
	specialties, err := h.services.Directory.GetSpecialties(c.Request.Context(), departmentID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, specialties)
}
