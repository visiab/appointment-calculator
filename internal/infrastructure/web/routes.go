package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/visiab/appointment-calculator/internal/interfaces/http/controllers"
)

func SetupRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "appointment-calculator",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// This is where we would set up dependency injection
		// For now, we'll create placeholder controllers
		appointmentController := controllers.NewAppointmentController(nil, nil)
		scheduleController := controllers.NewScheduleController(nil)
		participantController := controllers.NewParticipantController()

		// Appointment routes
		appointments := v1.Group("/appointments")
		{
			appointments.POST("", appointmentController.CreateAppointment)
			appointments.GET("", appointmentController.ListAppointments)
			appointments.GET("/:id", appointmentController.GetAppointment)
			appointments.PUT("/:id", appointmentController.UpdateAppointment)
			appointments.DELETE("/:id", appointmentController.CancelAppointment)
		}

		// Schedule routes
		schedules := v1.Group("/schedules")
		{
			schedules.POST("/availability", scheduleController.FindAvailableTimeSlots)
			schedules.GET("/:owner_id/overview", scheduleController.GetScheduleOverview)
			schedules.GET("/:owner_id/detail", scheduleController.GetScheduleDetail)
			schedules.POST("/:owner_id/blocked-times", scheduleController.AddBlockedTime)
		}

		// Participant routes
		participants := v1.Group("/participants")
		{
			participants.POST("", participantController.CreateParticipant)
			participants.GET("/:id", participantController.GetParticipant)
			participants.PUT("/:id", participantController.UpdateParticipant)
			participants.POST("/:id/availability", participantController.AddAvailability)
			participants.GET("/:id/availability", participantController.GetAvailability)
		}
	}
}
