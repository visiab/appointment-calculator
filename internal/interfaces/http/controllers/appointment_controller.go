package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/visiab/appointment-calculator/internal/application/dto"
	"github.com/visiab/appointment-calculator/internal/application/usecases"
)

type AppointmentController struct {
	createUseCase *usecases.CreateAppointmentUseCase
	updateUseCase *usecases.UpdateAppointmentUseCase
}

func NewAppointmentController(
	createUseCase *usecases.CreateAppointmentUseCase,
	updateUseCase *usecases.UpdateAppointmentUseCase,
) *AppointmentController {
	return &AppointmentController{
		createUseCase: createUseCase,
		updateUseCase: updateUseCase,
	}
}

func (c *AppointmentController) CreateAppointment(ctx *gin.Context) {
	var request dto.CreateAppointmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	response, err := c.createUseCase.Execute(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create appointment",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AppointmentController) UpdateAppointment(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Appointment ID is required",
		})
		return
	}

	var request dto.UpdateAppointmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	response, err := c.updateUseCase.Execute(appointmentID, request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update appointment",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AppointmentController) CancelAppointment(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Appointment ID is required",
		})
		return
	}

	err := c.updateUseCase.Cancel(appointmentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to cancel appointment",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Appointment cancelled successfully",
	})
}

func (c *AppointmentController) GetAppointment(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Appointment ID is required",
		})
		return
	}

	// This would require a GetAppointmentUseCase, for now return not implemented
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error": "Get appointment endpoint not implemented yet",
	})
}

func (c *AppointmentController) ListAppointments(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	participantID := ctx.Query("participant_id")

	// This would require a ListAppointmentsUseCase, for now return not implemented
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error": "List appointments endpoint not implemented yet",
		"params": gin.H{
			"page":           page,
			"limit":          limit,
			"participant_id": participantID,
		},
	})
}
