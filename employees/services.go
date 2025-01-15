package employees

import (
	"leetify-test/database"
	"time"

	"github.com/spf13/cast"
)

type Service interface {
	CreateEmployee(dataEmp *database.Employee) error
	UpdateEmployee(dataEmp *database.Employee) error
	DeleteEmployee(employeeID string) error
	GetEmployeeByID(employeeID string) (employee ViewEmployee, err error)
	GetEmployees() (employee []ViewEmployee, err error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) CreateEmployee(dataEmp *database.Employee) error {
	unixTime := time.Now().Unix()
	dataEmp.EmployeeID = "EMP-" + cast.ToString(unixTime)
	return s.repository.CreateEmployee(dataEmp)
}

func (s *service) UpdateEmployee(dataEmp *database.Employee) error {
	return s.repository.UpdateEmployee(dataEmp)
}

func (s *service) DeleteEmployee(employeeID string) error {
	return s.repository.DeleteEmployee(employeeID)
}

func (s *service) GetEmployeeByID(employeeID string) (employee ViewEmployee, err error) {
	data, err := s.repository.GetEmployeeByID(employeeID)
	if err != nil {
		return employee, err
	}
	employee = ViewEmployee{
		EmployeeID:     data.EmployeeID,
		Name:           data.Name,
		Address:        data.Address,
		DepartmentID:   data.DepartmentID,
		DepartmentName: data.Department.DepartmentName,
	}

	return employee, nil
}

func (s *service) GetEmployees() (employees []ViewEmployee, err error) {
	data, err := s.repository.GetEmployees()
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		employees = append(employees, ViewEmployee{
			EmployeeID:     v.EmployeeID,
			Name:           v.Name,
			Address:        v.Address,
			DepartmentID:   v.DepartmentID,
			DepartmentName: v.Department.DepartmentName,
		})
	}

	return employees, nil
}
