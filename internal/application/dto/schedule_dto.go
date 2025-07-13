package dto

import (
	"time"
)

type AvailabilityQuery struct {
	ParticipantIDs []string  `json:"participant_ids" binding:"required,min=1"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
	Duration       int       `json:"duration_minutes" binding:"required,min=1"`
	Timezone       string    `json:"timezone"`
}

type TimeSlotResponse struct {
	StartTime            time.Time `json:"start_time"`
	EndTime              time.Time `json:"end_time"`
	Score                float64   `json:"score"`
	AvailableParticipants int      `json:"available_participants"`
	TotalParticipants    int       `json:"total_participants"`
	Reason               string    `json:"reason"`
	Conflicts            int       `json:"conflicts"`
}

type AvailabilityResult struct {
	AvailableSlots []TimeSlotResponse `json:"available_slots"`
	OptimalSlot    *TimeSlotResponse  `json:"optimal_slot,omitempty"`
	SearchPeriod   struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		Timezone  string    `json:"timezone"`
	} `json:"search_period"`
	Summary struct {
		TotalSlotsFound int `json:"total_slots_found"`
		Participants    int `json:"participants"`
	} `json:"summary"`
}

type ScheduleOverview struct {
	OwnerID           string    `json:"owner_id"`
	Timezone          string    `json:"timezone"`
	WorkingHoursStart string    `json:"working_hours_start"`
	WorkingHoursEnd   string    `json:"working_hours_end"`
	AppointmentsToday int       `json:"appointments_today"`
	NextAppointment   *AppointmentResponse `json:"next_appointment,omitempty"`
	TotalAppointments int       `json:"total_appointments"`
	LastUpdated       time.Time `json:"last_updated"`
}

type ScheduleDetail struct {
	ScheduleOverview
	Appointments []AppointmentResponse `json:"appointments"`
	BlockedTimes []struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		Reason    string    `json:"reason"`
	} `json:"blocked_times"`
}
