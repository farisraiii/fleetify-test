package database

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Department struct {
	ID              uint   `gorm:"primarykey" json:"id"`
	DepartmentName  string `gorm:"type:varchar(255);not null" json:"department_name"`
	MaxClockInTime  string `gorm:"type:time" json:"max_clock_in_time"`
	MaxClockOutTime string `gorm:"type:time" json:"max_clock_out_time"`
}

type Employee struct {
	EmployeeID   string     `gorm:"primarykey;unique;type:varchar(50);not null" json:"employee_id"`
	DepartmentID uint       `gorm:"not null" json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID;references:ID" json:"department"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Address      string     `gorm:"type:text" json:"address"`
}

type Attendance struct {
	EmployeeID   string   `gorm:"type:varchar(50);not null" json:"employee_id"`
	Employee     Employee `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee"`
	AttendanceID string   `gorm:"primarykey;type:varchar(100);not null" json:"attendance_id"`
	ClockIn      *string  `gorm:"type:timestamp" json:"clock_in"`
	ClockOut     *string  `gorm:"type:timestamp" json:"clock_out"`
}

type LogAttendance struct {
	BaseModel
	EmployeeID     string     `gorm:"type:varchar(50);not null" json:"employee_id"`
	Employee       Employee   `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee"`
	AttendanceID   string     `gorm:"type:varchar(100);not null" json:"attendance_id"`
	Attendance     Attendance `gorm:"foreignKey:AttendanceID;references:AttendanceID" json:"attendance"`
	DateAttendance string     `gorm:"type:timestamp;not null" json:"date_attendance"`
	AttendanceType int        `gorm:"type:tinyint(1);not null" json:"attendance_type"`
	Description    string     `gorm:"type:text" json:"description"`
}

func (d *Department) BeforeCreate(tx *gorm.DB) (err error) {
	var lastDepartment Department
	tx.Order("id DESC").Limit(1).Find(&lastDepartment)

	if lastDepartment.ID == 0 {
		d.ID = 1001001
	} else {
		d.ID = lastDepartment.ID + 1
	}

	return nil
}
