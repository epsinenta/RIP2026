package serializer

import "web_backend/internal/app/ds"

type DepartmentJSON struct {
	DepartmentID       uint    `json:"department_id"`
	IsDeleted         bool    `json:"is_deleted"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	PhotoURL          string  `json:"photo_url"`
	EmployeeCount     int     `json:"employee_count"`
	Head              string  `json:"head"`
	ReportsTo         string  `json:"reports_to"`
	Video                string  `json:"video"`
	ShortDescription     string  `json:"short_description"`
	ShortDescriptionEN   string  `json:"short_description_en"`
}

func DepartmentToJSON(department ds.Department) DepartmentJSON {
	return DepartmentJSON{
		DepartmentID:      department.DepartmentID,
		IsDeleted:         department.IsDeleted,
		Title:             department.Title,
		Description:       department.Description,
		PhotoURL:         department.Photo,
		EmployeeCount:    department.EmployeeCount,
		Head:             department.Head,
		ReportsTo:        department.ReportsTo,
		Video:                department.Video,
		ShortDescription:     department.ShortDescription,
		ShortDescriptionEN:   department.ShortDescriptionEN,
	}
}

func DepartmentFromJSON(j DepartmentJSON) ds.Department {
	return ds.Department{
		Title:            j.Title,
		Description:      j.Description,
		EmployeeCount:    j.EmployeeCount,
		Head:             j.Head,
		ReportsTo:        j.ReportsTo,
		Video:                j.Video,
		ShortDescription:     j.ShortDescription,
		ShortDescriptionEN:   j.ShortDescriptionEN,
	}
}
