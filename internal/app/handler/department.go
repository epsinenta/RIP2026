package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/ds"
	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

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

func (h *Handler) CreateDepartment(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")
	var j serializer.DepartmentJSON
	if strings.HasPrefix(contentType, "application/json") {
		if err := ctx.BindJSON(&j); err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	} else {
		title := ctx.PostForm("title")
		desc := ctx.PostForm("description")
		if title == "" || desc == "" {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("title and description are required"))
			return
		}
		empCount := 0
		if ec := ctx.PostForm("employee_count"); ec != "" {
			fmt.Sscanf(ec, "%d", &empCount)
		}
		j = serializer.DepartmentJSON{
			Title:            title,
			Description:      desc,
			EmployeeCount:    empCount,
			Head:             ctx.PostForm("head"),
			ReportsTo:        ctx.PostForm("reports_to"),
			Video:            ctx.PostForm("video"),
			ShortDescription: ctx.PostForm("short_description"),
		}
	}
	department, err := h.Repository.CreateDepartment(j)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if imageFile, err := ctx.FormFile("image"); err == nil {
		dep, err := h.Repository.AddPhoto(ctx, int(department.DepartmentID), imageFile)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		department = dep
	}
	if videoFile, err := ctx.FormFile("video"); err == nil {
		dep, err := h.Repository.AddVideo(ctx, int(department.DepartmentID), videoFile)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		department = dep
	}
	ctx.Header("Location", fmt.Sprintf("/api/department/%d", department.DepartmentID))
	ctx.JSON(http.StatusCreated, serializer.DepartmentToJSON(department))
}

func (h *Handler) AddToDepartmentApplication(ctx *gin.Context) {
	departmentIDStr := ctx.Param("id")
	departmentID, err := strconv.Atoi(departmentIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID := uint(h.Repository.GetCreatorID())
	app, created, err := h.Repository.GetDepartmentApplicationDraft(creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	err = h.Repository.AddDepartment(uint(departmentID), creatorID)
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

func (h *Handler) AddToDepartmentApplicationForm(ctx *gin.Context) {
	departmentIDStr := ctx.Param("department_id")
	if departmentIDStr == "" {
		departmentIDStr = ctx.PostForm("department_id")
	}
	if departmentIDStr == "" {
		if ctx.GetHeader("Accept") != "" && strings.Contains(ctx.GetHeader("Accept"), "application/json") {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("department_id is required"))
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}
	departmentID, err := strconv.Atoi(departmentIDStr)
	if err != nil {
		if ctx.GetHeader("Accept") != "" && strings.Contains(ctx.GetHeader("Accept"), "application/json") {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}
	creatorID := uint(h.Repository.GetCreatorID())
	app, created, err := h.Repository.GetDepartmentApplicationDraft(creatorID)
	if err != nil {
		if ctx.GetHeader("Accept") != "" && strings.Contains(ctx.GetHeader("Accept"), "application/json") {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}
	err = h.Repository.AddDepartment(uint(departmentID), creatorID)
	if err != nil {
		if ctx.GetHeader("Accept") != "" && strings.Contains(ctx.GetHeader("Accept"), "application/json") {
			if errors.Is(err, repository.ErrNotFound) {
				h.errorHandler(ctx, http.StatusNotFound, err)
			} else if errors.Is(err, repository.ErrAlreadyExists) {
				h.errorHandler(ctx, http.StatusConflict, err)
			} else {
				h.errorHandler(ctx, http.StatusInternalServerError, err)
			}
			return
		}
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}
	if ctx.GetHeader("Accept") != "" && strings.Contains(ctx.GetHeader("Accept"), "application/json") {
		creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
		incompleteCount, _ := h.Repository.GetIncompleteItemsCount(app.DepartmentApplicationID)
		status := http.StatusOK
		if created {
			ctx.Header("Location", fmt.Sprintf("/api/department_application/%d", app.DepartmentApplicationID))
			status = http.StatusCreated
		}
		ctx.JSON(status, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin, incompleteCount))
		return
	}
	redirectTo := ctx.GetHeader("Referer")
	if redirectTo == "" {
		redirectTo = "/"
	}
	ctx.Redirect(http.StatusSeeOther, redirectTo)
}
