package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visiab/appointment-calculator/internal/application/dto"
	"github.com/visiab/appointment-calculator/internal/application/usecases"
)

type ScheduleController struct {
	findAvailableTimeSlotsUseCase *usecases.FindAvailableTimeSlotsUseCase
}

func NewScheduleController(
	findAvailableTimeSlotsUseCase *usecases.FindAvailableTimeSlotsUseCase,
) *ScheduleController {
	return &ScheduleController{
		findAvailableTimeSlotsUseCase: findAvailableTimeSlotsUseCase,
	}
}

func (c *ScheduleController) FindAvailableTimeSlots(ctx *gin.Context) {
	var query dto.AvailabilityQuery
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	result, err := c.findAvailableTimeSlotsUseCase.Execute(query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to find available time slots",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *ScheduleController) GetScheduleOverview(ctx *gin.Context) {
	ownerID := ctx.Param("owner_id")
	if ownerID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Owner ID is required",
		})
		return
	}

	// This would require a GetScheduleOverviewUseCase
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error": "Get schedule overview endpoint not implemented yet",
		"owner_id": ownerID,
	})
}

func (c *ScheduleController) GetScheduleDetail(ctx *gin.Context) {
	ownerID := ctx.Param("owner_id")
	if ownerID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Owner ID is required",
		})
		return
	}

	// Parse query parameters for date range
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	// This would require a GetScheduleDetailUseCase
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error": "Get schedule detail endpoint not implemented yet",
		"params": gin.H{
			"owner_id":   ownerID,
			"start_date": startDate,
			"end_date":   endDate,
		},
	})
}

func (c *ScheduleController) AddBlockedTime(ctx *gin.Context) {
	ownerID := ctx.Param("owner_id")
	if ownerID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Owner ID is required",
		})
		return
	}

	var request struct {
		StartTime string `json:"start_time" binding:"required"`
		EndTime   string `json:"end_time" binding:"required"`
		Reason    string `json:"reason"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// This would require an AddBlockedTimeUseCase
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error": "Add blocked time endpoint not implemented yet",
		"params": gin.H{
			"owner_id": ownerID,
			"request":  request,
		},
	})
}
