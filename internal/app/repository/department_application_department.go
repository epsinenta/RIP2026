package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"web_backend/internal/app/ds"
	"web_backend/internal/app/serializer"
)

func (r *Repository) DeleteDepartmentFromDepartmentApplication(departmentApplicationID, departmentID int) (ds.DepartmentApplication, error) {
	var app ds.DepartmentApplication
	err := r.db.Where("department_application_id = ?", departmentApplicationID).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.DepartmentApplication{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, departmentApplicationID)
		}
		return ds.DepartmentApplication{}, err
	}
	err = r.db.Where("department_id = ? AND department_application_id = ?", departmentID, departmentApplicationID).
		Delete(&ds.DepartmentApplicationDepartment{}).Error
	if err != nil {
		return ds.DepartmentApplication{}, err
	}
	if err := r.RecalculateMainDepartments(app.DepartmentApplicationID); err != nil {
		return ds.DepartmentApplication{}, err
	}
	return app, nil
}

func (r *Repository) EditDepartmentFromDepartmentApplication(departmentApplicationID, departmentID int, j serializer.DepartmentApplicationDepartmentJSON) (ds.DepartmentApplicationDepartment, error) {
	var item ds.DepartmentApplicationDepartment
	err := r.db.Where("department_id = ? AND department_application_id = ?", departmentID, departmentApplicationID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.DepartmentApplicationDepartment{}, fmt.Errorf("%w: отдел в заявке", ErrNotFound)
		}
		return ds.DepartmentApplicationDepartment{}, err
	}
	updates := serializer.DepartmentApplicationDepartmentFromJSON(j)
	updatesMap := map[string]interface{}{
		"sort_order": updates.SortOrder,
		"role":       updates.Role,
	}
	if updates.Amount != nil {
		updatesMap["amount"] = updates.Amount
	}
	if updates.Salary != nil {
		updatesMap["salary"] = updates.Salary
	}
	err = r.db.Model(&item).Updates(updatesMap).Error
	if err != nil {
		return ds.DepartmentApplicationDepartment{}, err
	}
	if updates.Role != "" {
		var dep ds.Department
		if err := r.db.First(&dep, departmentID).Error; err == nil {
			salary := roleToBaseSalary(updates.Role) + float64(dep.EmployeeCount)*5000
			r.db.Model(&item).Update("salary", salary)
		}
	}
	if err := r.RecalculateMainDepartments(uint(departmentApplicationID)); err != nil {
		return ds.DepartmentApplicationDepartment{}, err
	}
	r.db.Where("department_id = ? AND department_application_id = ?", departmentID, departmentApplicationID).
		Preload("Department").Preload("MainDepartment").First(&item)
	return item, nil
}
