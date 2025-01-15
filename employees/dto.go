package employees

type ViewEmployee struct {
	EmployeeID     string `json:"employee_id"`
	DepartmentID   uint   `json:"department_id"`
	DepartmentName string `json:"department_name"`
	Name           string `json:"name"`
	Address        string `json:"address"`
}
