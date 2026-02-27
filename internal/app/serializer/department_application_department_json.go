package serializer

import "web_backend/internal/app/ds"

type DepartmentApplicationDepartmentJSON struct {
	DepartmentApplicationID uint     `json:"department_application_id"`
	DepartmentID           uint     `json:"department_id"`
	Amount                 *int     `json:"amount"`
	MainDepartmentID       *uint    `json:"main_department_id"`
	SortOrder              int      `json:"sort_order"`
	Role                   string   `json:"role"`
	Salary                 *float64 `json:"salary"`
}

func DepartmentApplicationDepartmentToJSON(d ds.DepartmentApplicationDepartment) DepartmentApplicationDepartmentJSON {
	return DepartmentApplicationDepartmentJSON{
		DepartmentApplicationID: d.DepartmentApplicationID,
		DepartmentID:           d.DepartmentID,
		Amount:                  d.Amount,
		MainDepartmentID:        d.MainDepartmentID,
		SortOrder:               d.SortOrder,
		Role:                    d.Role,
		Salary:                  d.Salary,
	}
}

func DepartmentApplicationDepartmentFromJSON(j DepartmentApplicationDepartmentJSON) ds.DepartmentApplicationDepartment {
	return ds.DepartmentApplicationDepartment{
		Amount:    j.Amount,
		SortOrder: j.SortOrder,
		Role:      j.Role,
		Salary:    j.Salary,
	}
}
