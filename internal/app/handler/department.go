package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/ds"
	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

// GetDepartments godoc
// @Summary Получить список отделов
// @Description Возвращает все отделы или фильтрует по названию
// @Tags departments
// @Produce json
// @Param Title query string false "Название отдела для поиска"
// @Success 200 {array} serializer.DepartmentJSON "Список отделов"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /departments [get]
func (h *Handler) GetDepartments(ctx *gin.Context) {
	var departments []ds.Department
	var err error
	searchQuery := ctx.Query("Title")
	if searchQuery == "" {
		departments, err = h.Repository.GetDepartments()
	} else {
		departments, err = h.Repository.GetDepartmentsByTitle(searchQuery)
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]serializer.DepartmentJSON, 0, len(departments))
	for _, d := range departments {
		resp = append(resp, serializer.DepartmentToJSON(d))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetDepartment godoc
// @Summary Получить отдел по ID
// @Description Возвращает информацию об отделе по идентификатору
// @Tags departments
// @Produce json
// @Param id path int true "ID отдела"
// @Success 200 {object} serializer.DepartmentJSON "Данные отдела"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Отдел не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /department/{id} [get]
func (h *Handler) GetDepartment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	department, err := h.Repository.GetDepartment(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.DepartmentToJSON(*department))
}

// CreateDepartment godoc
// @Summary Создать отдел
// @Description Создает новый отдел
// @Tags departments
// @Accept json
// @Produce json
// @Param department body serializer.DepartmentJSON true "Данные нового отдела"
// @Success 201 {object} serializer.DepartmentJSON "Созданный отдел"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /department/create-department [post]
func (h *Handler) CreateDepartment(ctx *gin.Context) {
	var j serializer.DepartmentJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	department, err := h.Repository.CreateDepartment(j)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Location", fmt.Sprintf("/api/department/%d", department.DepartmentID))
	ctx.JSON(http.StatusCreated, serializer.DepartmentToJSON(department))
}

// AddToDepartmentApplication godoc
// @Summary Добавить отдел в заявку
// @Description Добавляет отдел в заявку-черновик пользователя
// @Tags department_application_departments
// @Produce json
// @Param department_id path int true "ID отдела"
// @Success 200 {object} serializer.DepartmentApplicationJSON "Заявка с добавленным отделом"
// @Success 201 {object} serializer.DepartmentApplicationJSON "Создана новая заявка"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Отдел не найден"
// @Failure 409 {object} map[string]string "Отдел уже в заявке"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /dep_app_dep/add/{department_id} [post]
func (h *Handler) AddToDepartmentApplication(ctx *gin.Context) {
	departmentIDStr := ctx.Param("department_id")
	departmentID, err := strconv.Atoi(departmentIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	app, created, err := h.Repository.GetDepartmentApplicationDraft(uint(creatorID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	err = h.Repository.AddDepartment(uint(departmentID), uint(creatorID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
	incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
	status := http.StatusOK
	if created {
		ctx.Header("Location", fmt.Sprintf("/api/department_application/%d", app.DepartmentApplicationID))
		status = http.StatusCreated
	}
	ctx.JSON(status, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin, incompleteCount))
}
