package attendance

type AttendanceRequest struct {
	EmployeeID       string `json:"employee_id"`
	RemarkAttendance string `json:"remark_attendance"`
}
