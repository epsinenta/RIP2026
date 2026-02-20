package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(1)
	isDraft, err := h.Repository.IsDraftDepartmentApplication(id, creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isDraft {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	items, err := h.Repository.GetDepartmentApplication(id, creatorID)
	if err != nil {
		logrus.Error(err)
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	var totalSalary float64
	for _, item := range items {
		totalSalary += item.Salary
	}

	ctx.HTML(http.StatusOK, "department_application.html", gin.H{
		"department_application": items,
		"department_application_id": id,
		"totalSalary":            totalSalary,
	})
}

func (h *Handler) AddToDepartmentApplication(ctx *gin.Context) {
	departmentIDStr := ctx.PostForm("department_id")
	departmentID, err := strconv.Atoi(departmentIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(1)

	err = h.Repository.AddDepartment(uint(departmentID), creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}

func (h *Handler) DeleteDepartmentApplication(ctx *gin.Context) {
	appIDStr := ctx.PostForm("department_application_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteDepartmentApplication(uint(appID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) UpdateRole(ctx *gin.Context) {
	appIDStr := ctx.PostForm("department_application_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	departmentIDStr := ctx.PostForm("department_id")
	departmentID, err := strconv.Atoi(departmentIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	role := ctx.PostForm("role")
	if role == "" {
		h.errorHandler(ctx, http.StatusBadRequest, nil)
		return
	}

	err = h.Repository.UpdateRole(uint(appID), uint(departmentID), role)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/department_application/"+appIDStr)
}
