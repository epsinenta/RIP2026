package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"web_backend/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetDepartments(ctx *gin.Context) {
	var departments []repository.Department
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		departments, err = h.Repository.GetDepartments()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		departments, err = h.Repository.GetDepartmentByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	department_applications, err := h.Repository.GetDepartmentApplications()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"departments":             departments,
		"query":                   searchQuery,
		"department_applications": department_applications,
	})
}

func (h *Handler) GetDepartment(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	department, err := h.Repository.GetDepartment(id)
	if err != nil {
		logrus.Error(err)
	}

	appDep, err := h.Repository.GetDepartmentApplicationForDepartment(id)
	hasManager := err == nil && appDep != nil

	ctx.HTML(http.StatusOK, "department.html", gin.H{
		"department": department,
		"appDep":     appDep,
		"hasManager": hasManager,
	})
}

func (h *Handler) GetDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	department_application, err := h.Repository.GetDepartmentApplication(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "department_application.html", gin.H{
		"department_application": department_application,
		"totalSalary":            fmt.Sprintf("%.0f", department_application.TotalSalary),
	})
}
