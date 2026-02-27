package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

func (h *Handler) GetDepartmentApplicationCart(ctx *gin.Context) {
	creatorID := uint(h.Repository.GetUserID())
	count := h.Repository.GetDepartmentApplicationCount(creatorID)
	if count == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":                    "no_draft",
			"departments_count":         count,
		})
		return
	}
	app, err := h.Repository.CheckCurrentDepartmentApplicationDraft(creatorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, gin.H{
				"status":            "no_draft",
				"departments_count": 0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":                app.DepartmentApplicationID,
		"departments_count": h.Repository.GetDepartmentApplicationCount(creatorID),
	})
}

func (h *Handler) GetAllDepartmentApplications(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from, to time.Time
	if fromDate != "" {
		t, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = t
	}
	toDate := ctx.Query("to-date")
	if toDate != "" {
		t, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = t
	}
	status := ctx.Query("status")
	apps, err := h.Repository.GetAllDepartmentApplications(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]serializer.DepartmentApplicationJSON, 0, len(apps))
	for _, app := range apps {
		creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
		resp = append(resp, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin))
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	deps, app, err := h.Repository.GetDepartmentApplicationWithDepartments(id)
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
	depResp := make([]serializer.DepartmentJSON, 0, len(deps))
	for _, d := range deps {
		depResp = append(depResp, serializer.DepartmentToJSON(d))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"department_application": serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin),
		"departments":             depResp,
	})
}

func (h *Handler) EditDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.DepartmentApplicationJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	app, err := h.Repository.EditDepartmentApplication(id, j)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(app)
	ctx.JSON(http.StatusOK, serializer.DepartmentApplicationToJSON(app, creatorLogin, moderatorLogin))
}

func (h *Handler) FormDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	app, err := h.Repository.FormDepartmentApplication(id, "formed")
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

func (h *Handler) FinishDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var statusJSON serializer.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	app, err := h.Repository.FinishDepartmentApplication(id, statusJSON.Status)
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

func (h *Handler) DeleteDepartmentApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	_, err = h.Repository.FormDepartmentApplication(id, "deleted")
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
	ctx.JSON(http.StatusOK, gin.H{"message": "Department application deleted"})
}
