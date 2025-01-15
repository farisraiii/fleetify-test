package helpers

import (
	"fmt"
	"leetify-test/database"

	"gorm.io/gorm"
)

func CheckExistingUser(employeeID string, db *gorm.DB) bool {
	var count int64
	err := db.Model(&database.Employee{}).Where("employee_id = ?", employeeID).Count(&count).Error
	if err != nil {
		fmt.Println("Error checking user existence:", err)
		return false
	}

	if count < 1 {
		return false
	}
	return true
}
