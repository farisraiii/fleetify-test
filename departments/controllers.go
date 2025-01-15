package departments

import (
	"leetify-test/database"
	"leetify-test/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}

func (h *Controller) CreateDepartment(c *gin.Context) {
	var dataDep database.Department
	if err := c.ShouldBindJSON(&dataDep); err != nil {
		response := helpers.APIResponse("Invalid request", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
	}

	err := h.service.CreateDepartment(&dataDep)
	if err != nil {
		response := helpers.APIResponse("Failed to create department", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", dataDep)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) UpdateDepartment(c *gin.Context) {
	var dataDep database.Department
	if err := c.ShouldBindJSON(&dataDep); err != nil {
		response := helpers.APIResponse("Invalid request", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
	}

	err := h.service.UpdateDepartment(dataDep.ID, &dataDep)
	if err != nil {
		response := helpers.APIResponse("Failed to update department", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", dataDep)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) DeleteDepartment(c *gin.Context) {
	departmentID := c.Query("department_id")
	err := h.service.DeleteDepartment(cast.ToUint(departmentID))
	if err != nil {
		response := helpers.APIResponse("Failed to delete department", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) GetDepartmentByID(c *gin.Context) {
	departmentID := c.Query("department_id")
	department, err := h.service.GetDepartmentByID(cast.ToUint(departmentID))
	if err != nil {
		response := helpers.APIResponse("Failed to get department", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", department)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) GetDepartments(c *gin.Context) {
	departments, err := h.service.GetDepartments()
	if err != nil {
		response := helpers.APIResponse("Failed to get department", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success", http.StatusOK, "success", departments)
	c.JSON(http.StatusOK, response)
}
