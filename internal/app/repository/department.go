package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"web_backend/internal/app/ds"
	minioClient "web_backend/internal/app/minioClient"
	"web_backend/internal/app/serializer"
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

func (r *Repository) GetDepartment(id int) (*ds.Department, error) {
	var department ds.Department
	err := r.db.Where("department_id = ? AND is_deleted = ?", id, false).First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: отдел с id %d", ErrNotFound, id)
		}
		return nil, err
	}
	return &department, nil
}

func (r *Repository) GetDepartmentsByTitle(title string) ([]ds.Department, error) {
	var departments []ds.Department
	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+title+"%", false).Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *Repository) CreateDepartment(j serializer.DepartmentJSON) (ds.Department, error) {
	department := serializer.DepartmentFromJSON(j)
	err := r.db.Create(&department).Scan(&department).Error
	if err != nil {
		return ds.Department{}, err
	}
	return department, nil
}

func (r *Repository) EditDepartment(id int, j serializer.DepartmentJSON) (ds.Department, error) {
	var department ds.Department
	if id < 0 {
		return ds.Department{}, errors.New("id должно быть >= 0")
	}
	err := r.db.Where("department_id = ? AND is_deleted = ?", id, false).First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Department{}, fmt.Errorf("%w: отдел с id %d", ErrNotFound, id)
		}
		return ds.Department{}, err
	}
	updates := serializer.DepartmentFromJSON(j)
	err = r.db.Model(&department).Updates(updates).Error
	if err != nil {
		return ds.Department{}, err
	}
	r.db.Where("department_id = ?", id).First(&department)
	return department, nil
}

func (r *Repository) DeleteDepartment(id int) error {
	var department ds.Department
	if id < 0 {
		return errors.New("id должно быть >= 0")
	}
	err := r.db.Where("department_id = ? AND is_deleted = ?", id, false).First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: отдел с id %d", ErrNotFound, id)
		}
		return err
	}
	if department.Photo != "" {
		if err := minioClient.DeleteObject(context.Background(), r.mc, minioClient.GetImgBucket(), department.Photo); err != nil {
			return err
		}
	}
	return r.db.Model(&ds.Department{}).Where("department_id = ?", id).Update("is_deleted", true).Error
}

func (r *Repository) AddPhoto(ctx *gin.Context, departmentID int, file *multipart.FileHeader) (ds.Department, error) {
	department, err := r.GetDepartment(departmentID)
	if err != nil {
		return ds.Department{}, err
	}
	if department.Photo != "" {
		_ = minioClient.DeleteObject(ctx, r.mc, minioClient.GetImgBucket(), department.Photo)
	}
	fileName, err := minioClient.UploadImage(ctx, r.mc, minioClient.GetImgBucket(), file, *department)
	if err != nil {
		return ds.Department{}, err
	}
	if err := r.db.Model(&ds.Department{}).Where("department_id = ?", departmentID).Update("photo_url", fileName).Error; err != nil {
		return ds.Department{}, err
	}
	department.Photo = fileName
	return *department, nil
}

func (r *Repository) AddVideo(ctx *gin.Context, departmentID int, file *multipart.FileHeader) (ds.Department, error) {
	department, err := r.GetDepartment(departmentID)
	if err != nil {
		return ds.Department{}, err
	}
	if department.Video != "" {
		_ = minioClient.DeleteObject(ctx, r.mc, minioClient.GetImgBucket(), department.Video)
	}
	fileName, err := minioClient.UploadVideo(ctx, r.mc, minioClient.GetImgBucket(), file, department.DepartmentID)
	if err != nil {
		return ds.Department{}, err
	}
	if err := r.db.Model(&ds.Department{}).Where("department_id = ?", departmentID).Update("video", fileName).Error; err != nil {
		return ds.Department{}, err
	}
	department.Video = fileName
	return *department, nil
}
