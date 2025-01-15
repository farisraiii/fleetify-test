package attendance

import (
	"fmt"
	"leetify-test/database"
	"leetify-test/helpers"
	"time"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Repository interface {
	CheckInAttendance(request AttendanceRequest) error
	CheckOutAttendance(request AttendanceRequest) error
	GetLogHistoryAttendance(startDate, endDate string, departmentID uint) ([]map[string]interface{}, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CheckInAttendance(request AttendanceRequest) error {
	var dataAtt *database.Attendance
	var employee database.Employee

	if ok := helpers.CheckExistingUser(request.EmployeeID, r.db); !ok {
		return fmt.Errorf("employee dengan id %s tidak ditemukan", request.EmployeeID)
	}

	err := r.db.Where("employee_id = ? AND clock_out IS NULL", request.EmployeeID).
		Order("clock_in DESC").
		First(&dataAtt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err := r.db.Preload("Department").
				Where("employee_id = ?", request.EmployeeID).
				First(&employee).Error
			if err != nil {
				return fmt.Errorf("gagal memuat data employee: %w", err)

			}
			clockInTime := time.Now().Format("2006-01-02 15:04:05")
			timeParse, err := time.Parse("2006-01-02 15:04:05", clockInTime)
			if err != nil {
				return err
			}
			hourTime := timeParse.Format("15:04")

			maxTimeParse, err := time.Parse("15:04:05", employee.Department.MaxClockInTime)
			if err != nil {
				return fmt.Errorf("error parsing max clock in time: %w", err)
			}
			maxTime := maxTimeParse.Format("15:04")

			if hourTime > maxTime {
				if request.RemarkAttendance == "" {
					return fmt.Errorf("waktu check in melebihi batas maksimal %s, harap masukkan alasan anda", maxTime)
				}
				if len(request.RemarkAttendance) < 8 {
					return fmt.Errorf("alasan terlalu pendek, harap masukkan alasan minimal 8 karakter")
				}
			}

			dataAtt = &database.Attendance{
				EmployeeID:   request.EmployeeID,
				AttendanceID: "ATT-" + cast.ToString(time.Now().Unix()),
				ClockIn:      &clockInTime,
				ClockOut:     nil,
			}
			err = r.db.Create(dataAtt).Error
			if err != nil {
				return err
			}

			err = r.LogHistoryAttendance(request, &employee, dataAtt, 1)
			if err != nil {
				return err
			}

			return nil
		}
		return err
	}

	return fmt.Errorf("employee dengan id %s sudah melakukan check in, harap melakukan check out terlebih dahulu", request.EmployeeID)
}

func (r *repository) CheckOutAttendance(request AttendanceRequest) error {
	if ok := helpers.CheckExistingUser(request.EmployeeID, r.db); !ok {
		return fmt.Errorf("employee dengan id %s tidak ditemukan", request.EmployeeID)
	}

	var lastAttendance database.Attendance
	var employee database.Employee

	err := r.db.Where("employee_id = ? AND clock_out IS NULL", request.EmployeeID).
		Order("clock_in DESC").
		First(&lastAttendance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("employee dengan id %s belum melakukan check in, harap melakukan check in terlebih dahulu", request.EmployeeID)
		}
		return err
	}

	err = r.db.Preload("Department").
		Where("employee_id = ?", request.EmployeeID).
		First(&employee).Error
	if err != nil {
		return fmt.Errorf("gagal memuat data employee: %w", err)
	}
	clockOutTime := time.Now().Format("2006-01-02 15:04:05")
	timeParse, err := time.Parse("2006-01-02 15:04:05", clockOutTime)
	if err != nil {
		return err
	}
	hourTime := timeParse.Format("15:04")

	minTimeOut, err := time.Parse("15:04:05", employee.Department.MaxClockOutTime)
	if err != nil {
		return fmt.Errorf("error parsing max clock in time: %w", err)
	}
	minTime := minTimeOut.Format("15:04")

	if hourTime < minTime {
		if request.RemarkAttendance == "" {
			return fmt.Errorf("anda melakukan check out sebelum batas waktu %s, harap masukkan alasan anda", minTime)
		}
		if len(request.RemarkAttendance) < 8 {
			return fmt.Errorf("alasan terlalu pendek, harap masukkan alasan minimal 8 karakter")
		}
	}

	lastAttendance.ClockOut = &clockOutTime
	err = r.LogHistoryAttendance(request, &employee, &lastAttendance, 2)
	if err != nil {
		return err
	}

	return r.db.Save(&lastAttendance).Error
}

func (r *repository) LogHistoryAttendance(attRequest AttendanceRequest, empData *database.Employee, attData *database.Attendance, typeAbcense int) error {
	var dataAttHistory *database.LogAttendance
	var timeAttendance string

	if typeAbcense == 1 {
		timeAttendance = *attData.ClockIn
	} else if typeAbcense == 2 {
		timeAttendance = *attData.ClockOut
	}

	if attRequest.RemarkAttendance == "" {
		attRequest.RemarkAttendance = "On Time (Remark By System)"
	}

	dataAttHistory = &database.LogAttendance{
		EmployeeID:     attRequest.EmployeeID,
		AttendanceID:   attData.AttendanceID,
		DateAttendance: timeAttendance,
		AttendanceType: typeAbcense,
		Description:    attRequest.RemarkAttendance,
	}

	return r.db.Create(dataAttHistory).Error
}

func (r *repository) GetLogHistoryAttendance(startDate, endDate string, departmentID uint) ([]map[string]interface{}, error) {
	var whereClauses []string

	if departmentID != 0 {
		whereClauses = append(whereClauses, `department_id = `+cast.ToString(departmentID))
	}

	if startDate != "" && endDate != "" {
		whereClauses = append(whereClauses, `clock_in BETWEEN '`+startDate+`' AND '`+endDate+`'`)
	}

	var query string
	if len(whereClauses) > 0 {
		query = "WHERE " + whereClauses[0]
		for i := 1; i < len(whereClauses); i++ {
			query += " AND " + whereClauses[i]
		}
	}

	sql := `
		SELECT
			a.attendance_id,
			e.employee_id,
			e.name AS employee_name,
			d.department_name,
			d.max_clock_in_time,
			d.max_clock_out_time,
			GROUP_CONCAT(CASE WHEN la.attendance_type = 1 THEN la.date_attendance ELSE NULL END) AS check_in,
			GROUP_CONCAT(CASE WHEN la.attendance_type = 2 THEN la.date_attendance ELSE NULL END) AS check_out,
			GROUP_CONCAT(CASE WHEN la.attendance_type = 1 THEN la.description ELSE NULL END) AS remark_late,
			GROUP_CONCAT(CASE WHEN la.attendance_type = 2 THEN la.description ELSE NULL END) AS remark_early_out
		FROM log_attendances la
		JOIN employees e ON e.employee_id = la.employee_id
		JOIN departments d ON d.id = e.department_id
		JOIN attendances a ON a.attendance_id = la.attendance_id
		` + query + `
		GROUP BY a.attendance_id, e.employee_id, e.name
		ORDER BY check_in ASC, e.employee_id ASC;`

	var results []map[string]interface{}
	err := r.db.Raw(sql).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
