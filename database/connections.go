package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(config DBConnection) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("✅ Connected to database")
	return db, nil
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Department{})
	db.AutoMigrate(&Employee{})
	db.AutoMigrate(&Attendance{})
	db.AutoMigrate(&LogAttendance{})

	db.Exec("ALTER TABLE departments MODIFY COLUMN max_clock_in_time TIME NOT NULL")
	db.Exec("ALTER TABLE departments MODIFY COLUMN max_clock_out_time TIME NOT NULL")

	db.Exec(`
	ALTER TABLE employees
	ADD CONSTRAINT fk_employees_department_id
	FOREIGN KEY (department_id) REFERENCES departments(id)
	ON DELETE SET NULL;`)

	db.Exec(`
	ALTER TABLE attendances
	ADD CONSTRAINT fk_attendances_employee_id
	FOREIGN KEY (employee_id) REFERENCES employees(employee_id)
	ON DELETE CASCADE;`)

	db.Exec(`
    ALTER TABLE log_attendances
    ADD CONSTRAINT fk_log_attendances_attendance_id
    FOREIGN KEY (attendance_id) REFERENCES attendances(attendance_id)
    ON DELETE CASCADE;
`)

	db.Exec(`
	ADD CONSTRAINT fk_log_attendance_employee_id
    FOREIGN KEY (employee_id) REFERENCES employees(employee_id)
    ON DELETE CASCADE;`)

}
