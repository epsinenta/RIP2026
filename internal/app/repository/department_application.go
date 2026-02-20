package repository

import (
	"errors"
	"fmt"
	"time"
	"web_backend/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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
		Preload("Department").Find(&items).Error
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
			Title:     "Заявка на объединение департаментов",
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

		const k = 5000
		baseSalary := roleToBaseSalary("Головной")
		salary := baseSalary + float64(dep.EmployeeCount)*k

		item := ds.DepartmentApplicationDepartment{
			DepartmentApplicationID: app.DepartmentApplicationID,
			DepartmentID:            departmentID,
			Amount:                 1,
			IsMain:                 true,
			Role:                   "Головной",
			Salary:                 salary,
		}
		if err := r.db.Create(&item).Error; err != nil {
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
	const k = 5000
	salary := baseSalary + float64(dep.EmployeeCount)*k
	return r.db.Model(&ds.DepartmentApplicationDepartment{}).
		Where("department_application_id = ? AND department_id = ?", appID, departmentID).
		Updates(map[string]interface{}{"role": role, "salary": salary}).Error
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
