package repository

import (
	"errors"
	"fmt"
	"sort"
	"time"
	"web_backend/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func (r *Repository) GetDepartmentApplication(id int, creatorID uint) ([]ds.DepartmentApplicationDepartment, float64, error) {
	var app ds.DepartmentApplication
	err := r.db.Where("department_application_id = ? AND creator_id = ? AND status != ?",
		id, creatorID, "deleted").First(&app).Error
	if err != nil {
		return nil, 0, err
	}

	var items []ds.DepartmentApplicationDepartment
	err = r.db.Where("department_application_id = ?", id).
		Preload("Department").Preload("MainDepartment").
		Order("sort_order ASC, department_id ASC").Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	var totalSalary float64
	if app.TotalSalary != nil {
		totalSalary = *app.TotalSalary
	}
	return items, totalSalary, nil
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

	if count == 0 {
		var dep ds.Department
		if err := r.db.First(&dep, departmentID).Error; err != nil {
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
		if err := r.CalculateAndSetTotalSalary(app.DepartmentApplicationID); err != nil {
			return err
		}
	}

	return nil
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

func (r *Repository) CalculateAndSetTotalSalary(appID uint) error {
	var total float64
	err := r.db.Model(&ds.DepartmentApplicationDepartment{}).
		Where("department_application_id = ?", appID).
		Select("COALESCE(SUM(salary), 0)").
		Scan(&total).Error
	if err != nil {
		return err
	}
	return r.db.Model(&ds.DepartmentApplication{}).
		Where("department_application_id = ?", appID).
		Update("total_salary", total).Error
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
	if err := r.RecalculateMainDepartments(appID); err != nil {
		return err
	}
	return r.CalculateAndSetTotalSalary(appID)
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
