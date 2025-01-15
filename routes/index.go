package routes

import (
	"leetify-test/attendance"
	"leetify-test/departments"
	"leetify-test/employees"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routing(router *gin.Engine, db *gorm.DB) {
	versioning := router.Group("/api/v1")
	employeeRoute := versioning.Group("/employee")
	deptRoute := versioning.Group("/department")
	attRoute := versioning.Group("/attendance")

	//--Employee--//
	employeeRepository := employees.NewRepository(db)
	employeeService := employees.NewService(employeeRepository)
	employeeController := employees.NewController(employeeService)

	//--Department--//
	deptRepository := departments.NewRepository(db)
	deptService := departments.NewService(deptRepository)
	deptController := departments.NewController(deptService)

	//--Attendance--//
	attRepository := attendance.NewRepository(db)
	attService := attendance.NewService(attRepository)
	attController := attendance.NewController(attService)

	//--Route Employee--//
	employeeRoute.POST("/create", employeeController.CreateEmployee)
	employeeRoute.PUT("/update", employeeController.UpdateEmployee)
	employeeRoute.DELETE("/delete", employeeController.DeleteEmployee)
	employeeRoute.GET("/get-employee", employeeController.GetEmployeeByID)
	employeeRoute.GET("/", employeeController.GetEmployees)

	//--Route Department--//
	deptRoute.POST("/create", deptController.CreateDepartment)
	deptRoute.PUT("/update", deptController.UpdateDepartment)
	deptRoute.DELETE("/delete", deptController.DeleteDepartment)
	deptRoute.GET("/get-department", deptController.GetDepartmentByID)
	deptRoute.GET("/", deptController.GetDepartments)

	//--Route Attendance--//
	attRoute.POST("/check-in", attController.CheckInAttendance)
	attRoute.PUT("/check-out", attController.CheckOutAttendance)
	attRoute.GET("/log-history", attController.GetLogHistoryAttendance)

}
