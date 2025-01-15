package employees

import (
	"fmt"
	"leetify-test/database"

	"gorm.io/gorm"
)

type Repository interface {
	CreateEmployee(dataEmp *database.Employee) error
	UpdateEmployee(dataEmp *database.Employee) error
	DeleteEmployee(employeeID string) error
	GetEmployeeByID(employeeID string) (*database.Employee, error)
	GetEmployees() ([]database.Employee, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateEmployee(dataEmp *database.Employee) error {
	if err := r.db.Create(dataEmp).Error; err != nil {
		fmt.Println("Failed to create employee:", err)
		return err
	}
	return nil
}

func (r *repository) UpdateEmployee(dataEmp *database.Employee) error {
	var existingEmployee database.Employee
	if err := r.db.Where("employee_id = ?", dataEmp.EmployeeID).First(&existingEmployee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("employee dengan id %s tidak ditemukan", dataEmp.EmployeeID)
		}
		return fmt.Errorf("gagal mencari employee: %w", err)
	}

	existingEmployee.Name = dataEmp.Name
	existingEmployee.Address = dataEmp.Address
	if dataEmp.DepartmentID > 1 {
		existingEmployee.DepartmentID = dataEmp.DepartmentID
	}

	if err := r.db.Save(&existingEmployee).Error; err != nil {
		fmt.Println("Failed to update employee:", err)
		return err
	}

	return nil
}
func (r *repository) DeleteEmployee(employeeID string) error {
	var employee database.Employee

	if err := r.db.Where("employee_id = ?", employeeID).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("employee dengan id %s tidak ditemukan", employeeID)
		}
		return fmt.Errorf("gagal mencari employee: %v", err)
	}

	if err := r.db.Delete(&employee).Error; err != nil {
		return fmt.Errorf("gagal menghapus employee dengan id %s: %v", employeeID, err)
	}

	return nil
}

func (r *repository) GetEmployeeByID(employeeID string) (*database.Employee, error) {
	var employee database.Employee
	if err := r.db.Preload("Department").Where("employee_id = ?", employeeID).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee dengan id %s tidak ditemukan", employeeID)
		}
		return nil, fmt.Errorf("gagal mencari employee: %v", err)
	}
	return &employee, nil
}

func (r *repository) GetEmployees() ([]database.Employee, error) {
	var employees []database.Employee
	if err := r.db.Preload("Department").Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("gagal mencari employee: %v", err)
	}
	return employees, nil
}
