package dto

import "time"

type CreateParticipantRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Timezone string `json:"timezone"`
}

type ParticipantResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Timezone string `json:"timezone"`
}

type UpdateParticipantRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
}

type AddAvailabilityRequest struct {
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Recurring bool      `json:"recurring"`
	Pattern   string    `json:"pattern,omitempty"`
	Interval  int       `json:"interval,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type AvailabilityResponse struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Recurring bool      `json:"recurring"`
	Pattern   string    `json:"pattern,omitempty"`
}
