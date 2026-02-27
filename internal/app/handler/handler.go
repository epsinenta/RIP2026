package handler

import (
	"errors"

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

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/api/departments", h.GetDepartments)
	router.GET("/api/department/:id", h.GetDepartment)
	router.POST("/api/department/create-department", h.CreateDepartment)
	router.PUT("/api/department/:id/edit-department", h.EditDepartment)
	router.DELETE("/api/department/:id/delete-department", h.DeleteDepartment)
	router.POST("/api/department/:id/add-to-department_application", h.AddToDepartmentApplication)
	router.POST("/api/department/:id/add-photo", h.AddPhoto)

	router.GET("/api/department_application/department_application-cart", h.GetDepartmentApplicationCart)
	router.GET("/api/department_application/all-department_applications", h.GetAllDepartmentApplications)
	router.GET("/api/department_application/:id", h.GetDepartmentApplication)
	router.PUT("/api/department_application/:id/edit-department_application", h.EditDepartmentApplication)
	router.PUT("/api/department_application/:id/form-department_application", h.FormDepartmentApplication)
	router.PUT("/api/department_application/:id/finish-department_application", h.FinishDepartmentApplication)
	router.DELETE("/api/department_application/:id/delete-department_application", h.DeleteDepartmentApplication)

	router.DELETE("/api/dep_app/:department_id/:department_application_id", h.DeleteDepartmentFromDepartmentApplication)
	router.PUT("/api/dep_app/:department_id/:department_application_id", h.EditDepartmentFromDepartmentApplication)

	router.POST("/api/users/signup", h.CreateUser)
	router.GET("/api/users/info", h.GetInfo)
	router.PUT("/api/users/info", h.EditInfo)
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
