package serializer

import (
	"time"

	"web_backend/internal/app/ds"
)

type DepartmentApplicationJSON struct {
	ID             uint      `json:"department_application_id"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatorLogin   string     `json:"creator_login"`
	ModeratorLogin *string    `json:"moderator_login"`
	FormingDate    *time.Time `json:"forming_date"`
	FinishDate     *time.Time `json:"finish_date"`
	Title          *string   `json:"title"`
}

func DepartmentApplicationToJSON(app ds.DepartmentApplication, creatorLogin, moderatorLogin string) DepartmentApplicationJSON {
	var mLogin *string
	if moderatorLogin != "" {
		mLogin = &moderatorLogin
	}
	var formingDate, finishDate *time.Time
	if app.FormingDate != nil {
		formingDate = app.FormingDate
	}
	if app.FinishDate.Valid {
		finishDate = &app.FinishDate.Time
	}
	return DepartmentApplicationJSON{
		ID:             app.DepartmentApplicationID,
		Status:         app.Status,
		CreatedAt:      app.CreatedAt,
		CreatorLogin:   creatorLogin,
		ModeratorLogin: mLogin,
		FormingDate:    formingDate,
		FinishDate:     finishDate,
		Title:          app.Title,
	}
}

func DepartmentApplicationFromJSON(j DepartmentApplicationJSON) ds.DepartmentApplication {
	return ds.DepartmentApplication{
		Title: j.Title,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}
