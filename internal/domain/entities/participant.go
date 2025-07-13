package entities

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type Participant struct {
	id           string
	name         string
	email        string
	timezone     *time.Location
	availability []valueobjects.TimeSlot
}

func NewParticipant(name, email string, timezone *time.Location) (*Participant, error) {
	if name == "" {
		return nil, errors.New("participant name cannot be empty")
	}
	
	if !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}
	
	if timezone == nil {
		timezone = time.UTC
	}
	
	return &Participant{
		id:           uuid.New().String(),
		name:         name,
		email:        email,
		timezone:     timezone,
		availability: make([]valueobjects.TimeSlot, 0),
	}, nil
}

func (p *Participant) ID() string {
	return p.id
}

func (p *Participant) Name() string {
	return p.name
}

func (p *Participant) Email() string {
	return p.email
}

func (p *Participant) Timezone() *time.Location {
	return p.timezone
}

func (p *Participant) Availability() []valueobjects.TimeSlot {
	return p.availability
}

func (p *Participant) AddAvailability(slot valueobjects.TimeSlot) {
	p.availability = append(p.availability, slot)
}

func (p *Participant) IsAvailableAt(timeRange valueobjects.TimeRange) bool {
	for _, slot := range p.availability {
		if slot.Contains(timeRange) {
			return true
		}
	}
	return false
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
