package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type Schedule struct {
	id           string
	ownerID      string
	timezone     *time.Location
	workingHours valueobjects.TimeRange
	appointments []*Appointment
	blockedTimes []valueobjects.TimeRange
}

func NewSchedule(ownerID string, timezone *time.Location, workingHours valueobjects.TimeRange) (*Schedule, error) {
	if ownerID == "" {
		return nil, errors.New("owner ID cannot be empty")
	}
	
	if timezone == nil {
		timezone = time.UTC
	}
	
	return &Schedule{
		id:           uuid.New().String(),
		ownerID:      ownerID,
		timezone:     timezone,
		workingHours: workingHours,
		appointments: make([]*Appointment, 0),
		blockedTimes: make([]valueobjects.TimeRange, 0),
	}, nil
}

func (s *Schedule) ID() string {
	return s.id
}

func (s *Schedule) OwnerID() string {
	return s.ownerID
}

func (s *Schedule) Timezone() *time.Location {
	return s.timezone
}

func (s *Schedule) WorkingHours() valueobjects.TimeRange {
	return s.workingHours
}

func (s *Schedule) Appointments() []*Appointment {
	return s.appointments
}

func (s *Schedule) BlockedTimes() []valueobjects.TimeRange {
	return s.blockedTimes
}

func (s *Schedule) AddAppointment(appointment *Appointment) error {
	if !s.isWithinWorkingHours(appointment.TimeRange()) {
		return errors.New("appointment is outside working hours")
	}
	
	if s.hasConflict(appointment.TimeRange()) {
		return errors.New("appointment conflicts with existing schedule")
	}
	
	s.appointments = append(s.appointments, appointment)
	return nil
}

func (s *Schedule) RemoveAppointment(appointmentID string) error {
	for i, appointment := range s.appointments {
		if appointment.ID() == appointmentID {
			s.appointments = append(s.appointments[:i], s.appointments[i+1:]...)
			return nil
		}
	}
	return errors.New("appointment not found")
}

func (s *Schedule) AddBlockedTime(timeRange valueobjects.TimeRange) {
	s.blockedTimes = append(s.blockedTimes, timeRange)
}

func (s *Schedule) IsAvailable(timeRange valueobjects.TimeRange) bool {
	return s.isWithinWorkingHours(timeRange) && !s.hasConflict(timeRange)
}

func (s *Schedule) isWithinWorkingHours(timeRange valueobjects.TimeRange) bool {
	return s.workingHours.Contains(timeRange)
}

func (s *Schedule) hasConflict(timeRange valueobjects.TimeRange) bool {
	// Check appointments
	for _, appointment := range s.appointments {
		if appointment.Status() != StatusCancelled && appointment.TimeRange().OverlapsWith(timeRange) {
			return true
		}
	}
	
	// Check blocked times
	for _, blocked := range s.blockedTimes {
		if blocked.OverlapsWith(timeRange) {
			return true
		}
	}
	
	return false
}
