package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"web_backend/internal/app/ds"
	"web_backend/internal/app/serializer"
)

func roleToLevel(role string) int {
	switch role {
	case "Головной":
		return 3
	case "Руководящий":
		return 2
	default:
		return 1
	}
}

func (r *Repository) RecalculateMainDepartments(appID uint) error {
	var items []ds.DepartmentApplicationDepartment
	if err := r.db.Where("department_application_id = ?", appID).Find(&items).Error; err != nil {
		return err
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].SortOrder != items[j].SortOrder {
			return items[i].SortOrder < items[j].SortOrder
		}
		return roleToLevel(items[i].Role) > roleToLevel(items[j].Role)
	})
	for i := range items {
		var mainID uint
		if roleToLevel(items[i].Role) == 3 {
			mainID = items[i].DepartmentID
		} else {
			for j := i - 1; j >= 0; j-- {
				if roleToLevel(items[j].Role) > roleToLevel(items[i].Role) {
					mainID = items[j].DepartmentID
					break
				}
			}
			if mainID == 0 {
				mainID = items[i].DepartmentID
			}
		}
		if err := r.db.Model(&ds.DepartmentApplicationDepartment{}).
			Where("department_application_id = ? AND department_id = ?", appID, items[i].DepartmentID).
			Update("main_department_id", mainID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetDepartmentApplicationCount(creatorID uint) int64 {
	var appID uint
	err := r.db.Model(&ds.DepartmentApplication{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("department_application_id").First(&appID).Error
	if err != nil {
		return 0
	}

	var count int64
	err = r.db.Model(&ds.DepartmentApplicationDepartment{}).
		Where("department_application_id = ?", appID).Count(&count).Error
	if err != nil {
		logrus.Error("Error counting department_application_departments:", err)
	}
	return count
}

func (r *Repository) GetActiveDepartmentApplicationID(creatorID uint) uint {
	var appID uint
	err := r.db.Model(&ds.DepartmentApplication{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("department_application_id").First(&appID).Error
	if err != nil {
		return 0
	}
	return appID
}

func (r *Repository) GetDepartmentApplication(id int, creatorID uint) ([]ds.DepartmentApplicationDepartment, error) {
	var app ds.DepartmentApplication
	err := r.db.Where("department_application_id = ? AND creator_id = ? AND status != ?",
		id, creatorID, "deleted").First(&app).Error
	if err != nil {
		return nil, err
	}

	var items []ds.DepartmentApplicationDepartment
	err = r.db.Where("department_application_id = ?", id).
		Preload("Department").Preload("MainDepartment").
		Order("sort_order ASC, department_id ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) AddDepartment(departmentID uint, creatorID uint) error {
	var app ds.DepartmentApplication

	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").
		First(&app).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		app = ds.DepartmentApplication{
			Status:    "draft",
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
		if err := r.db.Create(&app).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var count int64
	r.db.Model(&ds.DepartmentApplicationDepartment{}).
		Where("department_application_id = ? AND department_id = ?", app.DepartmentApplicationID, departmentID).
		Count(&count)

	if count > 0 {
		return fmt.Errorf("%w: отдел %d уже в заявке %d", ErrAlreadyExists, departmentID, app.DepartmentApplicationID)
	}

	{
		var dep ds.Department
		if err := r.db.First(&dep, departmentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: отдел с id %d", ErrNotFound, departmentID)
			}
			return err
		}

		var maxOrder int
		r.db.Model(&ds.DepartmentApplicationDepartment{}).
			Where("department_application_id = ?", app.DepartmentApplicationID).
			Select("COALESCE(MAX(sort_order), -1)").Scan(&maxOrder)

		amount := 1
		baseSalary := roleToBaseSalary("Головной")
		salary := baseSalary + float64(dep.EmployeeCount)*5000

		item := ds.DepartmentApplicationDepartment{
			DepartmentApplicationID: app.DepartmentApplicationID,
			DepartmentID:            departmentID,
			Amount:                  &amount,
			MainDepartmentID:        &departmentID,
			SortOrder:               maxOrder + 1,
			Role:                    "Головной",
			Salary:                  &salary,
		}
		if err := r.db.Create(&item).Error; err != nil {
			return err
		}
		if err := r.RecalculateMainDepartments(app.DepartmentApplicationID); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetDepartmentApplicationDraft(creatorID uint) (ds.DepartmentApplication, bool, error) {
	app, err := r.CheckCurrentDepartmentApplicationDraft(creatorID)
	if errors.Is(err, ErrNoDraft) {
		app = ds.DepartmentApplication{
			Status:    "draft",
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
		if err := r.db.Create(&app).Error; err != nil {
			return ds.DepartmentApplication{}, false, err
		}
		return app, true, nil
	}
	if err != nil {
		return ds.DepartmentApplication{}, false, err
	}
	return app, false, nil
}

func (r *Repository) CheckCurrentDepartmentApplicationDraft(creatorID uint) (ds.DepartmentApplication, error) {
	if creatorID == 0 {
		return ds.DepartmentApplication{}, ErrNotAllowed
	}
	var app ds.DepartmentApplication
	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&app)
	if res.Error != nil {
		return ds.DepartmentApplication{}, res.Error
	}
	if res.RowsAffected == 0 {
		return ds.DepartmentApplication{}, ErrNoDraft
	}
	return app, nil
}

func (r *Repository) GetModeratorAndCreatorLogin(app ds.DepartmentApplication) (string, string, error) {
	var creator ds.Users
	if err := r.db.Where("user_id = ?", app.CreatorID).First(&creator).Error; err != nil {
		return "", "", err
	}
	var moderatorLogin string
	if app.ModeratorID != nil && *app.ModeratorID != 0 {
		var moderator ds.Users
		if err := r.db.Where("user_id = ?", *app.ModeratorID).First(&moderator).Error; err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}
	return creator.Login, moderatorLogin, nil
}

func (r *Repository) GetAllDepartmentApplications(from, to time.Time, status string) ([]ds.DepartmentApplication, error) {
	var apps []ds.DepartmentApplication
	sub := r.db.Where("status != ? AND status != ?", "deleted", "draft")
	if !from.IsZero() {
		sub = sub.Where("forming_date > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("forming_date < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("department_application_id").Find(&apps).Error
	return apps, err
}

func (r *Repository) GetSingleDepartmentApplication(id int) (ds.DepartmentApplication, error) {
	if id < 0 {
		return ds.DepartmentApplication{}, errors.New("неверное id")
	}
	var app ds.DepartmentApplication
	err := r.db.Where("department_application_id = ?", id).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.DepartmentApplication{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.DepartmentApplication{}, err
	}
	if app.Status == "deleted" {
		return ds.DepartmentApplication{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return app, nil
}

func (r *Repository) GetDepartmentApplicationWithDepartments(id int) ([]ds.Department, ds.DepartmentApplication, error) {
	app, err := r.GetSingleDepartmentApplication(id)
	if err != nil {
		return nil, ds.DepartmentApplication{}, err
	}
	var deps []ds.Department
	sub := r.db.Table("department_application_departments").Where("department_application_id = ?", app.DepartmentApplicationID)
	err = r.db.Where("department_id IN (?) AND is_deleted = ?", sub.Select("department_id"), false).Find(&deps).Error
	if err != nil {
		return nil, ds.DepartmentApplication{}, err
	}
	return deps, app, nil
}

func (r *Repository) FormDepartmentApplication(id int, status string) (ds.DepartmentApplication, error) {
	app, err := r.GetSingleDepartmentApplication(id)
	if err != nil {
		return ds.DepartmentApplication{}, err
	}
	if app.Status != "draft" {
		return ds.DepartmentApplication{}, fmt.Errorf("эта заявка не может быть %s", status)
	}
	if status != "deleted" {
		if app.Title == nil || *app.Title == "" {
			defaultTitle := fmt.Sprintf("Заявка №%d", id)
			app.Title = &defaultTitle
			r.db.Model(&app).Update("title", defaultTitle)
		}
		items, _ := r.GetDepartmentApplicationItems(int(app.DepartmentApplicationID))
		for _, item := range items {
			if item.Role == "" {
				return ds.DepartmentApplication{}, errors.New("укажите роль для каждого отдела")
			}
		}
		for _, item := range items {
			var dep ds.Department
			if err := r.db.First(&dep, item.DepartmentID).Error; err != nil {
				return ds.DepartmentApplication{}, err
			}
			salary := roleToBaseSalary(item.Role) + float64(dep.EmployeeCount)*5000
			if err := r.db.Model(&ds.DepartmentApplicationDepartment{}).
				Where("department_application_id = ? AND department_id = ?", app.DepartmentApplicationID, item.DepartmentID).
				Update("salary", salary).Error; err != nil {
				return ds.DepartmentApplication{}, err
			}
		}
	}
	formingDate := time.Now()
	err = r.db.Model(&app).Updates(map[string]interface{}{
		"status":       status,
		"forming_date": formingDate,
	}).Error
	if err != nil {
		return ds.DepartmentApplication{}, err
	}
	app.Status = status
	app.FormingDate = &formingDate
	return app, nil
}

func (r *Repository) GetDepartmentApplicationItems(appID int) ([]ds.DepartmentApplicationDepartment, error) {
	var items []ds.DepartmentApplicationDepartment
	err := r.db.Where("department_application_id = ?", appID).
		Preload("Department").Preload("MainDepartment").
		Order("sort_order ASC, department_id ASC").Find(&items).Error
	return items, err
}

func (r *Repository) EditDepartmentApplication(id int, j serializer.DepartmentApplicationJSON) (ds.DepartmentApplication, error) {
	var app ds.DepartmentApplication
	if id < 0 {
		return ds.DepartmentApplication{}, errors.New("неправильное id")
	}
	err := r.db.Where("department_application_id = ? AND status != ?", id, "deleted").First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.DepartmentApplication{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.DepartmentApplication{}, err
	}
	updates := serializer.DepartmentApplicationFromJSON(j)
	err = r.db.Model(&app).Updates(updates).Error
	if err != nil {
		return ds.DepartmentApplication{}, err
	}
	r.db.Where("department_application_id = ?", id).First(&app)
	return app, nil
}

func (r *Repository) FinishDepartmentApplication(id int, status string) (ds.DepartmentApplication, error) {
	if status != "completed" && status != "rejected" {
		return ds.DepartmentApplication{}, errors.New("неверный статус")
	}
	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.DepartmentApplication{}, err
	}
	if !user.IsModerator {
		return ds.DepartmentApplication{}, fmt.Errorf("%w: вы не модератор", ErrNotAllowed)
	}
	app, err := r.GetSingleDepartmentApplication(id)
	if err != nil {
		return ds.DepartmentApplication{}, err
	}
	if app.Status != "formed" {
		return ds.DepartmentApplication{}, fmt.Errorf("эта заявка не может быть %s", status)
	}
	finishDate := time.Now()
	err = r.db.Model(&app).Updates(map[string]interface{}{
		"status":      status,
		"finish_date": finishDate,
		"moderator_id": user.UserID,
	}).Error
	if err != nil {
		return ds.DepartmentApplication{}, err
	}
	app.Status = status
	app.FinishDate = sql.NullTime{Time: finishDate, Valid: true}
	if app.ModeratorID == nil {
		uid := user.UserID
		app.ModeratorID = &uid
	} else {
		*app.ModeratorID = user.UserID
	}
	if status == "completed" {
		items, err := r.GetDepartmentApplicationItems(int(app.DepartmentApplicationID))
		if err != nil {
			return ds.DepartmentApplication{}, err
		}
		for _, item := range items {
			var dep ds.Department
			if err := r.db.First(&dep, item.DepartmentID).Error; err != nil {
				return ds.DepartmentApplication{}, err
			}
			salary := roleToBaseSalary(item.Role) + float64(dep.EmployeeCount)*5000
			if item.Salary == nil {
				item.Salary = &salary
			} else {
				*item.Salary = salary
			}
			r.db.Model(&item).Update("salary", salary)
		}
	}
	return app, nil
}

func (r *Repository) DeleteDepartmentApplication(appID uint) error {
	query := `
		UPDATE department_applications 
		SET status = 'deleted'
		WHERE department_application_id = $1;
	`
	result := r.db.Exec(query, appID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("department_application with id %d not found", appID)
	}
	return nil
}

func roleToBaseSalary(role string) float64 {
	switch role {
	case "Головной":
		return 200000
	case "Руководящий":
		return 150000
	default:
		return 100000
	}
}

func (r *Repository) UpdateRole(appID, departmentID uint, role string) error {
	var dep ds.Department
	if err := r.db.First(&dep, departmentID).Error; err != nil {
		return err
	}
	baseSalary := roleToBaseSalary(role)
	salary := baseSalary + float64(dep.EmployeeCount)*5000
	if err := r.db.Model(&ds.DepartmentApplicationDepartment{}).
		Where("department_application_id = ? AND department_id = ?", appID, departmentID).
		Updates(map[string]interface{}{"role": role, "salary": salary}).Error; err != nil {
		return err
	}
	return r.RecalculateMainDepartments(appID)
}

func (r *Repository) MoveDepartmentInApplication(appID, departmentID uint, direction int) error {
	var items []ds.DepartmentApplicationDepartment
	if err := r.db.Where("department_application_id = ?", appID).
		Order("sort_order ASC, department_id ASC").Find(&items).Error; err != nil {
		return err
	}
	idx := -1
	for i := range items {
		if items[i].DepartmentID == departmentID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("department %d not in application %d", departmentID, appID)
	}
	newIdx := idx + direction
	if newIdx < 0 || newIdx >= len(items) {
		return nil
	}
	items[idx], items[newIdx] = items[newIdx], items[idx]
	for i := range items {
		if err := r.db.Model(&ds.DepartmentApplicationDepartment{}).
			Where("department_application_id = ? AND department_id = ?", appID, items[i].DepartmentID).
			Update("sort_order", i).Error; err != nil {
			return err
		}
	}
	return r.RecalculateMainDepartments(appID)
}

func (r *Repository) IsDraftDepartmentApplication(appID int, creatorID uint) (bool, error) {
	var app ds.DepartmentApplication
	err := r.db.Select("status").Where("department_application_id = ? AND creator_id = ?",
		appID, creatorID).First(&app).Error
	if err != nil {
		return false, err
	}
	return app.Status == "draft", nil
}
