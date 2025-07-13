package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type AppointmentStatus string

const (
	StatusScheduled AppointmentStatus = "scheduled"
	StatusCancelled AppointmentStatus = "cancelled"
	StatusCompleted AppointmentStatus = "completed"
)

type Appointment struct {
	id          string
	title       string
	timeRange   valueobjects.TimeRange
	attendees   []string
	location    string
	status      AppointmentStatus
	createdAt   time.Time
	updatedAt   time.Time
}

func NewAppointment(title string, timeRange valueobjects.TimeRange, attendees []string, location string) (*Appointment, error) {
	if title == "" {
		return nil, errors.New("appointment title cannot be empty")
	}
	
	if len(attendees) == 0 {
		return nil, errors.New("appointment must have at least one attendee")
	}
	
	now := time.Now()
	return &Appointment{
		id:        uuid.New().String(),
		title:     title,
		timeRange: timeRange,
		attendees: attendees,
		location:  location,
		status:    StatusScheduled,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (a *Appointment) ID() string {
	return a.id
}

func (a *Appointment) Title() string {
	return a.title
}

func (a *Appointment) TimeRange() valueobjects.TimeRange {
	return a.timeRange
}

func (a *Appointment) Attendees() []string {
	return a.attendees
}

func (a *Appointment) Location() string {
	return a.location
}

func (a *Appointment) Status() AppointmentStatus {
	return a.status
}

func (a *Appointment) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Appointment) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Appointment) Cancel() {
	a.status = StatusCancelled
	a.updatedAt = time.Now()
}

func (a *Appointment) Complete() {
	a.status = StatusCompleted
	a.updatedAt = time.Now()
}

func (a *Appointment) Reschedule(newTimeRange valueobjects.TimeRange) {
	a.timeRange = newTimeRange
	a.updatedAt = time.Now()
}

func (a *Appointment) HasConflictWith(other *Appointment) bool {
	return a.timeRange.OverlapsWith(other.timeRange)
}
