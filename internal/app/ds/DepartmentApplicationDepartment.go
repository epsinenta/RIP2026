package ds

type DepartmentApplicationDepartment struct {
	DepartmentApplicationID uint    `gorm:"primaryKey;column:department_application_id"`
	DepartmentID            uint    `gorm:"primaryKey;column:department_id"`
	Amount                  int     `gorm:"not null;default:1"`
	IsMain                  bool    `gorm:"column:is_main;default:false"`
	Role                    string  `gorm:"type:varchar(100)"`
	Salary                  float64 `gorm:"type:numeric(12,2)"`

	Department           Department           `gorm:"foreignKey:DepartmentID"`
	DepartmentApplication DepartmentApplication `gorm:"foreignKey:DepartmentApplicationID"`
}

func (DepartmentApplicationDepartment) TableName() string {
	return "department_application_departments"
}
