package departments

import (
	"fmt"
	"leetify-test/database"

	"gorm.io/gorm"
)

type Repository interface {
	CreateDepartment(dataDep *database.Department) error
	UpdateDepartment(departmentID uint, dataDep *database.Department) error
	DeleteDepartment(departmentID uint) error
	GetDepartmentByID(departmentID uint) (*database.Department, error)
	GetDepartments() ([]database.Department, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	if db == nil {
		panic("Received nil DB instance in repository constructor")
	}
	return &repository{db}
}

func (r *repository) CreateDepartment(dataDep *database.Department) error {
	if err := r.db.Create(dataDep).Error; err != nil {
		fmt.Println("Failed to create department:", err)
		return err
	}

	return nil
}

func (r *repository) UpdateDepartment(departmentID uint, dataDep *database.Department) error {
	var department database.Department
	if err := r.db.First(&department, departmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("department with ID %d not found", departmentID)
		}
		return fmt.Errorf("failed to find department: %v", err)
	}

	if dataDep.MaxClockInTime != "" && dataDep.MaxClockOutTime != "" {
		if dataDep.MaxClockInTime > dataDep.MaxClockOutTime {
			return fmt.Errorf("waktu maksimal check in tidak boleh lebih besar dari waktu maksimal check out")
		}
	} else if dataDep.MaxClockInTime == "" && dataDep.MaxClockOutTime != "" {
		if department.MaxClockInTime > dataDep.MaxClockOutTime {
			return fmt.Errorf("waktu maksimal check in tidak boleh lebih besar dari waktu maksimal check out")
		}
	} else if dataDep.MaxClockInTime != "" && dataDep.MaxClockOutTime == "" {
		if dataDep.MaxClockInTime > department.MaxClockOutTime {
			return fmt.Errorf("waktu maksimal check in tidak boleh lebih besar dari waktu maksimal check out")
		}
	}

	if err := r.db.Model(&department).Updates(dataDep).Error; err != nil {
		return fmt.Errorf("failed to update department: %v", err)
	}

	return nil
}

func (r *repository) DeleteDepartment(departmentID uint) error {
	var department database.Department
	if err := r.db.First(&department, departmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("department with ID %d not found", departmentID)
		}
		return fmt.Errorf("failed to find department: %v", err)
	}
	if err := r.db.Delete(&department).Error; err != nil {
		return fmt.Errorf("failed to delete department: %v", err)
	}

	return nil
}

func (r *repository) GetDepartmentByID(departmentID uint) (*database.Department, error) {
	var department database.Department
	if err := r.db.First(&department, departmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("department with ID %d not found", departmentID)
		}
		return nil, fmt.Errorf("failed to get department: %v", err)
	}
	return &department, nil
}

func (r *repository) GetDepartments() ([]database.Department, error) {
	var departments []database.Department
	if err := r.db.Find(&departments).Error; err != nil {
		return nil, fmt.Errorf("failed to get departments: %v", err)
	}
	return departments, nil
}
