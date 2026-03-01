package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"web_backend/internal/app/ds"
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

func (h *Handler) getMinioURL() string {
	if url := os.Getenv("MINIO_URL"); url != "" {
		return url
	}
	return "http://localhost:9000/test"
}

func (h *Handler) DepartmentPage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	department, err := h.Repository.GetDepartment(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "department.html", gin.H{
		"department": *department,
		"minioUrl":  h.getMinioURL(),
	})
}

func (h *Handler) DepartmentApplicationPage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	_, err = h.Repository.GetSingleDepartmentApplication(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) || errors.Is(err, repository.ErrNotAllowed) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	items, err := h.Repository.GetDepartmentApplicationItems(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "department_application.html", gin.H{
		"department_application_id": id,
		"department_application":   items,
		"minioUrl":                 h.getMinioURL(),
	})
}

func (h *Handler) Index(ctx *gin.Context) {
	query := ctx.Query("query")
	var departments []ds.Department
	var err error
	if query == "" {
		departments, err = h.Repository.GetDepartments()
	} else {
		departments, err = h.Repository.GetDepartmentsByTitle(query)
	}
	if err != nil {
		departments = []ds.Department{}
	}

	creatorID := uint(h.Repository.GetCreatorID())
	count := int(h.Repository.GetDepartmentApplicationCount(creatorID))
	appID := h.Repository.GetActiveDepartmentApplicationID(creatorID)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"query":                     query,
		"departments":                departments,
		"department_application_count": count,
		"department_application_id":   appID,
		"minioUrl":                   h.getMinioURL(),
	})
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.Index)
	router.GET("/department/:id", h.DepartmentPage)
	router.GET("/department_application/:id", h.DepartmentApplicationPage)
	router.POST("/department_application/add", h.AddToDepartmentApplicationFromForm)
	router.POST("/department_application/delete", h.DeleteDepartmentApplicationFromForm)
	router.POST("/department_application/move", h.MoveDepartmentInApplicationFromForm)
	router.POST("/department_application/update_role", h.UpdateRoleFromForm)
	router.GET("/api/departments", h.GetDepartments)
	router.GET("/api/department/:id", h.GetDepartment)
	router.POST("/api/department/create-department", h.CreateDepartment)
	router.POST("/api/department/:id/add-to-department_application", h.AddToDepartmentApplication)

	router.GET("/api/department_application/department_application-cart", h.GetDepartmentApplicationCart)
	router.GET("/api/department_application/all-department_applications", h.GetAllDepartmentApplications)
	router.GET("/api/department_application/:id", h.GetDepartmentApplication)
	router.PUT("/api/department_application/:id/edit-department_application", h.EditDepartmentApplication)
	router.PUT("/api/department_application/:id/form-department_application", h.FormDepartmentApplication)
	router.PUT("/api/department_application/:id/finish-department_application", h.FinishDepartmentApplication)
	router.DELETE("/api/department_application/:id/delete-department_application", h.DeleteDepartmentApplication)

	router.DELETE("/api/dep_app_dep/:department_id/:department_application_id", h.DeleteDepartmentFromDepartmentApplication)
	router.PUT("/api/dep_app_dep/:department_id/:department_application_id", h.EditDepartmentFromDepartmentApplication)

	router.POST("/api/users/signup", h.CreateUser)
	router.POST("/api/users/signin", h.SignIn)
	router.POST("/api/users/signout", h.SignOut)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}
