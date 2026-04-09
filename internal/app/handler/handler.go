package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	router.Use(CORSMiddleware())

	api := router.Group("/api")

	unauthorized := api.Group("/")
	unauthorized.POST("/users/signup", h.CreateUser)
	unauthorized.POST("/users/signin", h.SignIn)
	unauthorized.GET("/departments", h.GetDepartments)
	unauthorized.GET("/department/:id", h.GetDepartment)

	optionalauthorized := api.Group("/")
	optionalauthorized.Use(h.WithOptionalAuthCheck())
	optionalauthorized.GET("/department_application/department_application-cart", h.GetDepartmentApplicationCart)

	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))
	authorized.POST("/department/create-department", h.CreateDepartment)
	authorized.GET("/department_application/all-department_applications", h.GetAllDepartmentApplications)
	authorized.GET("/department_application/:id", h.GetDepartmentApplication)
	authorized.PUT("/department_application/:id/edit-department_application", h.EditDepartmentApplication)
	authorized.PUT("/department_application/:id/form-department_application", h.FormDepartmentApplication)
	authorized.DELETE("/department_application/:id/delete-department_application", h.DeleteDepartmentApplication)
	authorized.POST("/dep_app_dep/add/:department_id", h.AddToDepartmentApplication)
	authorized.DELETE("/dep_app_dep/:department_id/:department_application_id", h.DeleteDepartmentFromDepartmentApplication)
	authorized.PUT("/dep_app_dep/:department_id/:department_application_id", h.EditDepartmentFromDepartmentApplication)
	authorized.POST("/users/signout", h.SignOut)

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	moderator.PUT("/department_application/:id/finish-department_application", h.FinishDepartmentApplication)

	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
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
		"description": errorMessage,
	})
}
