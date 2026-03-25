package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

// GetDepartmentApplicationCart godoc
// @Summary Получить корзину заявки
// @Description Возвращает информацию о текущей заявке-черновике пользователя
// @Tags department_applications
// @Produce json
// @Success 200 {object} map[string]interface{} "Данные корзины или no_draft"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /department_application/department_application-cart [get]
func (h *Handler) GetDepartmentApplicationCart(ctx *gin.Context) {
	creatorID, err := getUserID(ctx)
	if err != nil || creatorID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"has_draft":             false,
			"departments_count":     0,
			"incomplete_items_count": 0,
		})
		return
	}
	count := h.Repository.GetDepartmentApplicationCount(uint(creatorID))
	if count == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"has_draft":             false,
			"departments_count":     count,
			"incomplete_items_count": 0,
		})
		return
	}
	app, err := h.Repository.CheckCurrentDepartmentApplicationDraft(uint(creatorID))
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, gin.H{
				"has_draft":             false,
				"departments_count":     0,
				"incomplete_items_count": 0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
	ctx.JSON(http.StatusOK, gin.H{
		"id":                    app.DepartmentApplicationID,
		"has_draft":             true,
		"departments_count":     h.Repository.GetDepartmentApplicationCount(uint(creatorID)),
		"incomplete_items_count": incompleteCount,
	})
}

// GetAllDepartmentApplications godoc
// @Summary Получить список заявок
// @Description Возвращает заявки с возможностью фильтрации по датам и статусу
// @Tags department_applications
// @Produce json
// @Param from-date query string false "Начальная дата (YYYY-MM-DD)"
// @Param to-date query string false "Конечная дата (YYYY-MM-DD)"
// @Param status query string false "Статус заявки"
// @Success 200 {array} serializer.DepartmentApplicationJSON "Список заявок"
// @Failure 400 {object} map[string]string "Неверный формат даты"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /department_application/all-department_applications [get]
func (h *Handler) GetAllDepartmentApplications(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from, to time.Time
	if fromDate != "" {
		t, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = t
	}
	toDate := ctx.Query("to-date")
	if toDate != "" {
		t, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = t
	}
	status := ctx.Query("status")
	userID, _ := getUserID(ctx)
	creatorID := uint(0)
	if isModVal, ok := ctx.Get("is_moderator"); ok {
		if isMod, _ := isModVal.(bool); !isMod {
			creatorID = userID
		}
	} else {
		creatorID = userID
	}
	apps, err := h.Repository.GetAllDepartmentApplications(from, to, status, creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]serializer.DepartmentApplicationJSON, 0, len(apps))
	for _, app := range apps {
		creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
		incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
		resp = append(resp, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin, incompleteCount))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetDepartmentApplication godoc
// @Summary Получить заявку по ID
// @Description Возвращает полную информацию о заявке
// @Tags department_applications
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{} "Данные заявки с устройствами"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /department_application/{id} [get]
func (h *Handler) GetDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	app, err := h.Repository.GetSingleDepartmentApplication(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	userID, _ := getUserID(ctx)
	if isModVal, ok := ctx.Get("is_moderator"); ok {
		if isMod, _ := isModVal.(bool); !isMod && app.CreatorID != userID {
			h.errorHandler(ctx, http.StatusForbidden, repository.ErrNotAllowed)
			return
		}
	} else if app.CreatorID != userID {
		h.errorHandler(ctx, http.StatusForbidden, repository.ErrNotAllowed)
		return
	}
	items, err := h.Repository.GetDepartmentApplicationItems(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
	incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
	itemsResp := make([]serializer.DepartmentApplicationDepartmentJSON, 0, len(items))
	for _, item := range items {
		itemsResp = append(itemsResp, serializer.DepartmentApplicationDepartmentToJSON(h.Repository.EnsureItemSalary(&item)))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"department_application": serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin, incompleteCount),
		"items":                  itemsResp,
	})
}

// EditDepartmentApplication godoc
// @Summary Изменить заявку
// @Description Обновляет данные заявки
// @Tags department_applications
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param application body serializer.DepartmentApplicationJSON true "Новые данные заявки"
// @Success 200 {object} serializer.DepartmentApplicationJSON "Обновленная заявка"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /department_application/{id}/edit-department_application [put]
func (h *Handler) EditDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.DepartmentApplicationJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	app, err := h.Repository.EditDepartmentApplication(id, j)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
	incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
	ctx.JSON(http.StatusOK, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin, incompleteCount))
}

// FormDepartmentApplication godoc
// @Summary Сформировать заявку
// @Description Переводит заявку в статус "formed"
// @Tags department_applications
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} serializer.DepartmentApplicationJSON "Сформированная заявка"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /department_application/{id}/form-department_application [put]
func (h *Handler) FormDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	app, err := h.Repository.FormDepartmentApplication(id, "formed")
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
	incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
	ctx.JSON(http.StatusOK, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin, incompleteCount))
}

// FinishDepartmentApplication godoc
// @Summary Завершить заявку
// @Description Изменяет статус заявки (только для модераторов)
// @Tags department_applications
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param status body serializer.StatusJSON true "Новый статус (completed/rejected)"
// @Success 200 {object} serializer.DepartmentApplicationJSON "Результат модерации"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /department_application/{id}/finish-department_application [put]
func (h *Handler) FinishDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	moderatorID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	var statusJSON serializer.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	app, err := h.Repository.FinishDepartmentApplication(id, statusJSON.Status, moderatorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
	incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
	ctx.JSON(http.StatusOK, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin, incompleteCount))
}

// DeleteDepartmentApplication godoc
// @Summary Удалить заявку
// @Description Выполняет логическое удаление заявки
// @Tags department_applications
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /department_application/{id}/delete-department_application [delete]
func (h *Handler) DeleteDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	_, err = h.Repository.FormDepartmentApplication(id, "deleted")
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Department application deleted"})
}
