package repository

import (
	"fmt"
	"web_backend/internal/app/ds"
)

func (r *Repository) GetDepartments() ([]ds.Department, error) {
	var departments []ds.Department
	err := r.db.Where("is_deleted = ?", false).Find(&departments).Error
	if err != nil {
		return nil, err
	}
	if len(departments) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return departments, nil
}

func (r *Repository) GetDepartment(id int) (ds.Department, error) {
	var department ds.Department
	err := r.db.Where("department_id = ? AND is_deleted = ?", id, false).First(&department).Error
	if err != nil {
		return ds.Department{}, err
	}
	return department, nil
}

func (r *Repository) GetDepartmentsByTitle(title string) ([]ds.Department, error) {
	var departments []ds.Department
	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+title+"%", false).Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}
