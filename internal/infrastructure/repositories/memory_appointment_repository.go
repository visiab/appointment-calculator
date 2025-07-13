package repositories

import (
	"errors"
	"sync"

	"github.com/visiab/appointment-calculator/internal/domain/entities"
)

type MemoryAppointmentRepository struct {
	appointments map[string]*entities.Appointment
	mu           sync.RWMutex
}

func NewMemoryAppointmentRepository() *MemoryAppointmentRepository {
	return &MemoryAppointmentRepository{
		appointments: make(map[string]*entities.Appointment),
	}
}

func (r *MemoryAppointmentRepository) Save(appointment *entities.Appointment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.appointments[appointment.ID()] = appointment
	return nil
}

func (r *MemoryAppointmentRepository) FindByID(id string) (*entities.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	appointment, exists := r.appointments[id]
	if !exists {
		return nil, errors.New("appointment not found")
	}
	return appointment, nil
}

func (r *MemoryAppointmentRepository) FindByParticipant(participantID string) ([]*entities.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*entities.Appointment
	for _, appointment := range r.appointments {
		for _, attendee := range appointment.Attendees() {
			if attendee == participantID {
				result = append(result, appointment)
				break
			}
		}
	}
	return result, nil
}

func (r *MemoryAppointmentRepository) Update(appointment *entities.Appointment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.appointments[appointment.ID()]; !exists {
		return errors.New("appointment not found")
	}
	
	r.appointments[appointment.ID()] = appointment
	return nil
}

func (r *MemoryAppointmentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.appointments[id]; !exists {
		return errors.New("appointment not found")
	}
	
	delete(r.appointments, id)
	return nil
}

func (r *MemoryAppointmentRepository) FindAll() ([]*entities.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*entities.Appointment, 0, len(r.appointments))
	for _, appointment := range r.appointments {
		result = append(result, appointment)
	}
	return result, nil
}
