package dto

import (
	"time"

	"github.com/visiab/appointment-calculator/internal/domain/entities"
)

type CreateAppointmentRequest struct {
	Title     string    `json:"title" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Attendees []string  `json:"attendees" binding:"required,min=1"`
	Location  string    `json:"location"`
}

type CreateAppointmentResponse struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
	Attendees []string               `json:"attendees"`
	Location  string                 `json:"location"`
	Status    entities.AppointmentStatus `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
}

type UpdateAppointmentRequest struct {
	Title     *string    `json:"title,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Attendees []string   `json:"attendees,omitempty"`
	Location  *string    `json:"location,omitempty"`
}

type AppointmentResponse struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
	Duration  string                 `json:"duration"`
	Attendees []string               `json:"attendees"`
	Location  string                 `json:"location"`
	Status    entities.AppointmentStatus `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type AppointmentListResponse struct {
	Appointments []AppointmentResponse `json:"appointments"`
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	Limit        int                   `json:"limit"`
}
