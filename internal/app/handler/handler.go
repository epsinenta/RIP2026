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

	applications, err := h.Repository.GetApplications()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"departments":  departments,
		"query":        searchQuery,
		"applications": applications,
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

	appDep, err := h.Repository.GetApplicationForDepartment(id)
	hasManager := err == nil && appDep != nil

	ctx.HTML(http.StatusOK, "department.html", gin.H{
		"department": department,
		"appDep":     appDep,
		"hasManager": hasManager,
	})
}

func (h *Handler) GetApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	app, err := h.Repository.GetApplication(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "application.html", gin.H{
		"app":            app,
		"totalSalary":    fmt.Sprintf("%.0f", app.TotalSalary),
	})
}
