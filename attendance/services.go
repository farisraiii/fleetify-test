package attendance

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

type Service interface {
	CheckInAttendance(request AttendanceRequest) error
	CheckOutAttendance(request AttendanceRequest) error
	GetLogHistoryAttendance(startDate, endDate string, departmentID uint) ([]map[string]interface{}, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) CheckInAttendance(request AttendanceRequest) error {
	return s.repository.CheckInAttendance(request)
}

func (s *service) CheckOutAttendance(request AttendanceRequest) error {
	return s.repository.CheckOutAttendance(request)
}

func (s *service) GetLogHistoryAttendance(startDate, endDate string, departmentID uint) ([]map[string]interface{}, error) {
	result, err := s.repository.GetLogHistoryAttendance(startDate, endDate, departmentID)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		var earlyOut string
		var late string
		timeInParse, err := time.Parse("2006-01-02 15:04:05", cast.ToString(v["check_in"]))
		if err != nil {
			return nil, err
		}
		timeCompareIn := timeInParse.Format("15:04:05")
		if timeCompareIn > cast.ToString(v["max_clock_in_time"]) {
			timeCompareParse, err := time.Parse("15:04:05", timeCompareIn)
			if err != nil {
				fmt.Println("Error parsing check in:", err)
				return nil, err
			}
			clokcInDep, err := time.Parse("15:04:05", cast.ToString(v["max_clock_in_time"]))
			if err != nil {
				fmt.Println("Error parsing max clock in time:", err)
				return nil, err
			}

			timeLate := timeCompareParse.Sub(clokcInDep)

			hours := int(timeLate.Hours())
			minutes := int(timeLate.Minutes()) % 60
			seconds := int(timeLate.Seconds()) % 60

			late = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		}
		v["check_in"] = timeInParse.Format("2006-01-02 15:04:05")

		if v["check_out"] != nil {
			timeOut, err := time.Parse("2006-01-02 15:04:05", cast.ToString(v["check_out"]))
			if err != nil {
				return nil, err
			}
			timeCompareOut := timeOut.Format("15:04:05")

			if timeCompareOut < cast.ToString(v["max_clock_out_time"]) {
				timeCompareParse, err := time.Parse("15:04:05", timeCompareOut)
				if err != nil {
					fmt.Println("Error parsing check out:", err)
					return nil, err
				}
				clokcOutDep, err := time.Parse("15:04:05", cast.ToString(v["max_clock_out_time"]))
				if err != nil {
					fmt.Println("Error parsing max clock out time:", err)
					return nil, err
				}

				timeEarly := clokcOutDep.Sub(timeCompareParse)

				hours := int(timeEarly.Hours())
				minutes := int(timeEarly.Minutes()) % 60
				seconds := int(timeEarly.Seconds()) % 60

				earlyOut = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
			}
			v["check_out"] = timeOut.Format("2006-01-02 15:04:05")
		}

		v["time_early_out"] = earlyOut
		v["time_late"] = late
		if v["check_out"] == nil || v["check_out"] == "" {
			v["check_out"] = "-"
			v["time_early_out"] = "-"
			v["remark_early_out"] = "-"
		}

		if v["time_early_out"] == nil || v["time_early_out"] == "" {
			v["early_out"] = "-"
			v["remark_early_out"] = "-"
			v["time_early_out"] = "-"
		}

		if v["time_late"] == nil || v["time_late"] == "" {
			v["late"] = "-"
			v["remark_late"] = "-"
		}

	}

	return result, nil
}
