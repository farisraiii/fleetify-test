package attendance

import (
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

func (h *Controller) CheckInAttendance(c *gin.Context) {
	var request AttendanceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response := helpers.APIResponse("Invalid request", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
	}

	err := h.service.CheckInAttendance(request)
	if err != nil {
		response := helpers.APIResponse("Failed to abcense", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success check in", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) CheckOutAttendance(c *gin.Context) {
	var request AttendanceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response := helpers.APIResponse("Invalid request", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
	}

	err := h.service.CheckOutAttendance(request)
	if err != nil {
		response := helpers.APIResponse("Failed to abcense", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success check out", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *Controller) GetLogHistoryAttendance(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	departmentID := c.Query("department_id")
	output, err := h.service.GetLogHistoryAttendance(startDate, endDate, cast.ToUint(departmentID))
	if err != nil {
		response := helpers.APIResponse("Failed to abcense", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.APIResponse("Success check out", http.StatusOK, "success", output)
	c.JSON(http.StatusOK, response)
}
