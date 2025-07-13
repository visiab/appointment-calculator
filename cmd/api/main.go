package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/visiab/appointment-calculator/internal/infrastructure/config"
	"github.com/visiab/appointment-calculator/internal/infrastructure/dependency"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize dependency container
	container := dependency.NewContainer()

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Setup routes with dependency injection
	setupRoutes(router, container)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRoutes(router *gin.Engine, container *dependency.Container) {
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
		// Appointment routes
		appointments := v1.Group("/appointments")
		{
			appointments.POST("", container.AppointmentController.CreateAppointment)
			appointments.GET("", container.AppointmentController.ListAppointments)
			appointments.GET("/:id", container.AppointmentController.GetAppointment)
			appointments.PUT("/:id", container.AppointmentController.UpdateAppointment)
			appointments.DELETE("/:id", container.AppointmentController.CancelAppointment)
		}

		// Schedule routes
		schedules := v1.Group("/schedules")
		{
			schedules.POST("/availability", container.ScheduleController.FindAvailableTimeSlots)
			schedules.GET("/:owner_id/overview", container.ScheduleController.GetScheduleOverview)
			schedules.GET("/:owner_id/detail", container.ScheduleController.GetScheduleDetail)
			schedules.POST("/:owner_id/blocked-times", container.ScheduleController.AddBlockedTime)
		}

		// Participant routes
		participants := v1.Group("/participants")
		{
			participants.POST("", container.ParticipantController.CreateParticipant)
			participants.GET("/:id", container.ParticipantController.GetParticipant)
			participants.PUT("/:id", container.ParticipantController.UpdateParticipant)
			participants.POST("/:id/availability", container.ParticipantController.AddAvailability)
			participants.GET("/:id/availability", container.ParticipantController.GetAvailability)
		}
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}