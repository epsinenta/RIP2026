package ds

import (
	"database/sql"
	"time"
)

type DepartmentApplication struct {
	DepartmentApplicationID uint         `gorm:"primaryKey;column:department_application_id"`
	Status                  string       `gorm:"type:varchar(20);not null"`
	CreatedAt               time.Time    `gorm:"not null"`
	CreatorID               uint         `gorm:"not null"`
	FormingDate             *time.Time   `gorm:"column:forming_date"`
	FinishDate              sql.NullTime `gorm:"column:finish_date"`
	ModeratorID             *uint        `gorm:"column:moderator_id"`
	Title                   *string      `gorm:"type:varchar(255)"`
	TotalSalary             *float64     `gorm:"type:numeric(12,2)"`

	Creator   Users  `gorm:"foreignKey:CreatorID"`
	Moderator *Users `gorm:"foreignKey:ModeratorID"`
}

func (DepartmentApplication) TableName() string {
	return "department_applications"
}
