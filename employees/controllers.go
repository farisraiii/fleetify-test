package employees

import (
	"leetify-test/database"
	"leetify-test/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}

func (h *Controller) CreateEmployee(c *gin.Context) {
	var dataEmp database.Employee
	if err := c.ShouldBindJSON(&dataEmp); err != nil {
		response := helpers.APIResponse("Invalid request", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.service.CreateEmployee(&dataEmp)
	if err != nil {
		response := helpers.APIResponse("Failed to create user", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var output ViewEmployee
	output.EmployeeID = dataEmp.EmployeeID
	output.DepartmentID = dataEmp.DepartmentID
	output.Name = dataEmp.Name
	output.Address = dataEmp.Address

	response := helpers.APIResponse("Success", http.StatusOK, "success", output)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) UpdateEmployee(c *gin.Context) {
	var dataEmp database.Employee
	if err := c.ShouldBindJSON(&dataEmp); err != nil {
		response := helpers.APIResponse("Invalid request", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.service.UpdateEmployee(&dataEmp)
	if err != nil {
		response := helpers.APIResponse("Failed to update user", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) DeleteEmployee(c *gin.Context) {
	employeeID := c.Query("employee_id")
	err := h.service.DeleteEmployee(employeeID)
	if err != nil {
		response := helpers.APIResponse("Failed to delete user", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) GetEmployees(c *gin.Context) {
	employees, err := h.service.GetEmployees()
	if err != nil {
		response := helpers.APIResponse("Failed to get employees", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", employees)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) GetEmployeeByID(c *gin.Context) {
	employeeID := c.Query("employee_id")
	employee, err := h.service.GetEmployeeByID(employeeID)
	if err != nil {
		response := helpers.APIResponse("Failed to get employee", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", employee)
	c.JSON(http.StatusOK, response)
}
