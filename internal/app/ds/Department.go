package ds

type Department struct {
	DepartmentID     uint   `gorm:"primaryKey;column:department_id"`
	Title            string `gorm:"type:varchar(255);not null"`
	Description      string `gorm:"type:varchar(1000);not null"`
	IsDeleted        bool   `gorm:"type:boolean;not null;default:false"`
	Photo            string  `gorm:"column:photo_url;type:varchar(255)"`
	EmployeeCount    int    `gorm:"not null;default:0"`
	Head             string `gorm:"type:varchar(255)"`
	ReportsTo        string `gorm:"type:varchar(255)"`
	Video               string `gorm:"type:varchar(255)"`
	ShortDescription    string `gorm:"type:varchar(500)"`
	ShortDescriptionEN  string `gorm:"column:short_description_en;type:varchar(500)"`
}

func (Department) TableName() string {
	return "departments"
}
