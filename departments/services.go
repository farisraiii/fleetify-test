package departments

import (
	"fmt"
	"leetify-test/database"
)

type Service interface {
	CreateDepartment(dataDep *database.Department) error
	UpdateDepartment(departmentID uint, dataDep *database.Department) error
	DeleteDepartment(departmentID uint) error
	GetDepartmentByID(departmentID uint) (*database.Department, error)
	GetDepartments() ([]database.Department, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) CreateDepartment(dataDep *database.Department) error {
	if dataDep.MaxClockInTime > dataDep.MaxClockOutTime {
		return fmt.Errorf("waktu maksimal check in tidak boleh lebih besar dari waktu maksimal check out")
	}
	return s.repository.CreateDepartment(dataDep)
}

func (s *service) UpdateDepartment(departmentID uint, dataDep *database.Department) error {
	return s.repository.UpdateDepartment(departmentID, dataDep)
}

func (s *service) DeleteDepartment(departmentID uint) error {
	return s.repository.DeleteDepartment(departmentID)
}

func (s *service) GetDepartmentByID(departmentID uint) (*database.Department, error) {
	return s.repository.GetDepartmentByID(departmentID)
}

func (s *service) GetDepartments() ([]database.Department, error) {
	return s.repository.GetDepartments()
}
