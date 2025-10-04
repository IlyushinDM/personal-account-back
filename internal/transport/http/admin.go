package http

import (
	"net/http"
	"strconv"

	"lk/internal/models"
	"lk/internal/services"

	"github.com/gin-gonic/gin"
)

// --- Аутентификация ---

type adminLoginInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Вход для администратора
// @Tags         Admin Auth
// @Description  Аутентифицирует администратора и возвращает JWT.
// @Id           admin-login
// @Accept       json
// @Produce      json
// @Param        input body adminLoginInput true "Учетные данные"
// @Success      200 {object} map[string]string "accessToken"
// @Failure      400,401,500 {object} errorResponse
// @Router       /admin/login [post]
func (h *Handler) adminLogin(c *gin.Context) {
	var input adminLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	token, err := h.services.Admin.Login(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, token)
}

// --- Dashboard ---

// @Summary      Получить статистику для дашборда
// @Security     ApiKeyAuth
// @Tags         Admin Dashboard
// @Description  Возвращает ключевую статистику для главной страницы админ-панели.
// @Id           get-admin-dashboard
// @Produce      json
// @Success      200 {object} models.AdminDashboardStats
// @Failure      401,500 {object} errorResponse
// @Router       /admin/dashboard [get]
func (h *Handler) getAdminDashboard(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	stats, err := h.services.Admin.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// --- User (Пациент) ---

// @Summary      Получить список всех пациентов
// @Security     ApiKeyAuth
// @Tags         Admin Users
// @Description  Возвращает пагинированный список всех пользователей-пациентов.
// @Id           admin-get-all-users
// @Produce      json
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество на странице" default(10)
// @Success      200 {object} map[string]interface{} "items, total"
// @Failure      401,500 {object} errorResponse
// @Router       /admin/users [get]
func (h *Handler) adminGetAllUsers(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	params := models.PaginationParams{Page: page, Limit: limit}

	users, total, err := h.services.Admin.GetAllUsers(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": users, "total": total})
}

// @Summary      Получить детальную информацию о пациенте
// @Security     ApiKeyAuth
// @Tags         Admin Users
// @Id           admin-get-user-by-id
// @Produce      json
// @Param        id path int true "ID Пациента"
// @Success      200 {object} map[string]interface{} "user, profile"
// @Failure      400,401,404,500 {object} errorResponse
// @Router       /admin/users/{id} [get]
func (h *Handler) adminGetUserByID(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid user ID", err))
		return
	}

	user, profile, err := h.services.Admin.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"profile": profile,
	})
}

// @Summary      Обновить данные пациента
// @Security     ApiKeyAuth
// @Tags         Admin Users
// @Id           admin-update-user
// @Accept       json
// @Produce      json
// @Param        id path int true "ID Пациента"
// @Param        input body services.UpdateUserInput true "Обновляемые данные"
// @Success      200 {object} statusResponse
// @Failure      400,401,404,500 {object} errorResponse
// @Router       /admin/users/{id} [patch]
func (h *Handler) adminUpdateUser(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid user ID", err))
		return
	}
	var input services.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}
	if err := h.services.Admin.UpdateUser(c.Request.Context(), userID, input); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "user updated successfully"})
}

// @Summary      Удалить пациента
// @Security     ApiKeyAuth
// @Tags         Admin Users
// @Id           admin-delete-user
// @Param        id path int true "ID Пациента"
// @Success      204 "No Content"
// @Failure      400,401,500 {object} errorResponse
// @Router       /admin/users/{id} [delete]
func (h *Handler) adminDeleteUser(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid user ID", err))
		return
	}
	if err := h.services.Admin.DeleteUser(c.Request.Context(), userID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary      Получить все записи пациента (админ)
// @Security     ApiKeyAuth
// @Tags         Admin Users
// @Id           admin-get-user-appointments
// @Produce      json
// @Param        id path int true "ID Пациента"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество на странице" default(10)
// @Success      200 {object} map[string]interface{} "items, total"
// @Failure      400,401,500 {object} errorResponse
// @Router       /admin/users/{id}/appointments [get]
func (h *Handler) adminGetUserAppointments(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid user ID", err))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	params := models.PaginationParams{Page: page, Limit: limit}

	appointments, total, err := h.services.Admin.GetUserAppointments(c.Request.Context(), userID, params)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": appointments, "total": total})
}

// @Summary      Получить все анализы пациента (админ)
// @Security     ApiKeyAuth
// @Tags         Admin Users
// @Id           admin-get-user-analyses
// @Produce      json
// @Param        id path int true "ID Пациента"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество на странице" default(10)
// @Success      200 {object} map[string]interface{} "items, total"
// @Failure      400,401,500 {object} errorResponse
// @Router       /admin/users/{id}/analyses [get]
func (h *Handler) adminGetUserAnalyses(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid user ID", err))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	params := models.PaginationParams{Page: page, Limit: limit}

	analyses, total, err := h.services.Admin.GetUserAnalyses(c.Request.Context(), userID, params)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": analyses, "total": total})
}

// --- Doctor ---

// @Summary      Получить список всех врачей (админ)
// @Security     ApiKeyAuth
// @Tags         Admin Specialists
// @Id           admin-get-all-specialists
// @Produce      json
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество на странице" default(10)
// @Success      200 {object} map[string]interface{} "items, total"
// @Failure      401,500 {object} errorResponse
// @Router       /admin/specialists [get]
func (h *Handler) adminGetAllSpecialists(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	params := models.PaginationParams{Page: page, Limit: limit}

	doctors, total, err := h.services.Admin.GetAllSpecialists(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": doctors, "total": total})
}

// @Summary      Создать нового врача
// @Security     ApiKeyAuth
// @Tags         Admin Specialists
// @Id           admin-create-specialist
// @Accept       json
// @Produce      json
// @Param        input body services.CreateDoctorInput true "Данные нового врача"
// @Success      201 {object} map[string]uint64 "id"
// @Failure      400,401,500 {object} errorResponse
// @Router       /admin/specialists [post]
func (h *Handler) adminCreateSpecialist(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	var input services.CreateDoctorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}

	doctorID, err := h.services.Admin.CreateSpecialist(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": doctorID})
}

// @Summary      Получить детальную информацию о враче (админ)
// @Security     ApiKeyAuth
// @Tags         Admin Specialists
// @Id           admin-get-specialist-by-id
// @Produce      json
// @Param        id path int true "ID Врача"
// @Success      200 {object} models.Doctor
// @Failure      400,401,404,500 {object} errorResponse
// @Router       /admin/specialists/{id} [get]
func (h *Handler) adminGetSpecialistByID(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	doctorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid specialist ID", err))
		return
	}
	doctor, err := h.services.Admin.GetSpecialistByID(c.Request.Context(), doctorID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, doctor)
}

// @Summary      Обновить данные врача
// @Security     ApiKeyAuth
// @Tags         Admin Specialists
// @Id           admin-update-specialist
// @Accept       json
// @Produce      json
// @Param        id path int true "ID Врача"
// @Param        input body services.UpdateDoctorInput true "Обновляемые данные"
// @Success      200 {object} statusResponse
// @Failure      400,401,404,500 {object} errorResponse
// @Router       /admin/specialists/{id} [put]
func (h *Handler) adminUpdateSpecialist(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	doctorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid specialist ID", err))
		return
	}
	var input services.UpdateDoctorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}
	if err := h.services.Admin.UpdateSpecialist(c.Request.Context(), doctorID, input); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "specialist updated successfully"})
}

// @Summary      Удалить врача
// @Security     ApiKeyAuth
// @Tags         Admin Specialists
// @Id           admin-delete-specialist
// @Param        id path int true "ID Врача"
// @Success      204 "No Content"
// @Failure      400,401,500 {object} errorResponse
// @Router       /admin/specialists/{id} [delete]
func (h *Handler) adminDeleteSpecialist(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	doctorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid specialist ID", err))
		return
	}
	if err := h.services.Admin.DeleteSpecialist(c.Request.Context(), doctorID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary      Получить расписание врача
// @Security     ApiKeyAuth
// @Tags         Admin Specialists
// @Id           admin-get-specialist-schedule
// @Produce      json
// @Param        id path int true "ID Врача"
// @Success      200 {array} models.Schedule
// @Failure      400,401,500 {object} errorResponse
// @Router       /admin/specialists/{id}/schedule [get]
func (h *Handler) adminGetSpecialistSchedule(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	doctorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid specialist ID", err))
		return
	}
	schedule, err := h.services.Admin.GetSpecialistSchedule(c.Request.Context(), doctorID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, schedule)
}

// @Summary      Обновить расписание врача
// @Security     ApiKeyAuth
// @Tags         Admin Specialists
// @Id           admin-update-specialist-schedule
// @Accept       json
// @Produce      json
// @Param        id path int true "ID Врача"
// @Param        input body services.UpdateScheduleInput true "Новое расписание"
// @Success      200 {object} statusResponse
// @Failure      400,401,500 {object} errorResponse
// @Router       /admin/specialists/{id}/schedule [post]
func (h *Handler) adminUpdateSpecialistSchedule(c *gin.Context) {
	if _, err := getAdmin(c); err != nil {
		c.Error(services.NewInternalServerError("failed to identify admin from context", err))
		return
	}

	doctorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(services.NewBadRequestError("invalid specialist ID", err))
		return
	}
	var input services.UpdateScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(services.NewBadRequestError("invalid input body", err))
		return
	}
	if err := h.services.Admin.UpdateSpecialistSchedule(c.Request.Context(), doctorID, input); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "schedule updated successfully"})
}

// TODO: Реализовать
// --- Заглушки для других обработчиков админа ---

func (h *Handler) adminGetAllAppointments(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetAppointmentStats(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetAppointmentDetails(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminUpdateAppointmentStatus(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminDeleteAppointment(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetAllServices(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminCreateService(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminUpdateService(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminDeleteService(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetAllDepartments(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminCreateDepartment(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminUpdateDepartment(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminDeleteDepartment(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetAllAnalyses(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminCreateAnalysisResult(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminUpdateAnalysis(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminDeleteAnalysis(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetAllPrescriptions(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminCreatePrescription(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetFamilyRelations(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminDeleteFamilyRelation(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetAuditLogs(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetClinicSettings(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminUpdateClinicSettings(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetLegalDocs(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminCreateLegalDoc(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminUpdateLegalDoc(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminCreateBackup(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminGetBackupList(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}

func (h *Handler) adminRestoreFromBackup(c *gin.Context) {
	c.Error(services.NewInternalServerError("Not implemented yet", nil))
}
