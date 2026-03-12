package ds

type DepartmentApplicationDepartment struct {
	DepartmentApplicationID uint   `gorm:"primaryKey;column:department_application_id"`
	DepartmentID            uint   `gorm:"primaryKey;column:department_id"`
	MainDepartmentID        *uint  `gorm:"column:main_department_id"`
	SortOrder               int    `gorm:"column:sort_order;default:0"`
	Role                    string `gorm:"type:varchar(100)"`
	Salary                  *float64 `gorm:"type:numeric(12,2)"`

	Department            Department            `gorm:"foreignKey:DepartmentID"`
	MainDepartment        *Department           `gorm:"foreignKey:MainDepartmentID"`
	DepartmentApplication DepartmentApplication `gorm:"foreignKey:DepartmentApplicationID"`
}

func (DepartmentApplicationDepartment) TableName() string {
	return "department_application_departments"
}
