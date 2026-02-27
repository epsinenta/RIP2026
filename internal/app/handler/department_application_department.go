package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

func (h *Handler) DeleteDepartmentFromDepartmentApplication(ctx *gin.Context) {
	departmentID, err := strconv.Atoi(ctx.Param("department_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	departmentApplicationID, err := strconv.Atoi(ctx.Param("department_application_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	app, err := h.Repository.DeleteDepartmentFromDepartmentApplication(departmentApplicationID, departmentID)
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
	ctx.JSON(http.StatusOK, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin))
}

func (h *Handler) EditDepartmentFromDepartmentApplication(ctx *gin.Context) {
	departmentID, err := strconv.Atoi(ctx.Param("department_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	departmentApplicationID, err := strconv.Atoi(ctx.Param("department_application_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.DepartmentApplicationDepartmentJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	item, err := h.Repository.EditDepartmentFromDepartmentApplication(departmentApplicationID, departmentID, j)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.DepartmentApplicationDepartmentToJSON(item))
}
