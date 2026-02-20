package handler

import (
	"net/http"
	"strconv"
	"web_backend/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetDepartments(ctx *gin.Context) {
	var departments []ds.Department
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		departments, err = h.Repository.GetDepartments()
	} else {
		departments, err = h.Repository.GetDepartmentsByTitle(searchQuery)
	}
	if err != nil {
		logrus.Error(err)
	}

	creatorID := uint(1)
	appCount := h.Repository.GetDepartmentApplicationCount(creatorID)
	activeAppID := h.Repository.GetActiveDepartmentApplicationID(creatorID)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"departments":               departments,
		"query":                     searchQuery,
		"department_application_count": appCount,
		"department_application_id":   activeAppID,
	})
}

func (h *Handler) GetDepartment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	department, err := h.Repository.GetDepartment(id)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "department.html", gin.H{
		"department": department,
	})
}
