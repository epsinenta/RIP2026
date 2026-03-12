package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/ds"
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
	incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
	ctx.JSON(http.StatusOK, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin, incompleteCount))
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
	j.DepartmentID = uint(departmentID)
	j.DepartmentApplicationID = uint(departmentApplicationID)

	if j.Direction == "up" || j.Direction == "down" {
		direction := 1
		if j.Direction == "up" {
			direction = -1
		}
		err = h.Repository.MoveDepartmentInApplication(uint(departmentApplicationID), uint(departmentID), direction)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				h.errorHandler(ctx, http.StatusNotFound, err)
			} else {
				h.errorHandler(ctx, http.StatusInternalServerError, err)
			}
			return
		}
		items, err := h.Repository.GetDepartmentApplicationItems(departmentApplicationID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		var item *ds.DepartmentApplicationDepartment
		for i := range items {
			if items[i].DepartmentID == uint(departmentID) {
				item = &items[i]
				break
			}
		}
		if item != nil {
			ctx.JSON(http.StatusOK, serializer.DepartmentApplicationDepartmentToJSON(*item))
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "moved"})
		}
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
