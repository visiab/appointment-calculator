package repositories

import (
	"errors"
	"sync"

	"github.com/visiab/appointment-calculator/internal/domain/entities"
)

type MemoryScheduleRepository struct {
	schedules map[string]*entities.Schedule
	mu        sync.RWMutex
}

func NewMemoryScheduleRepository() *MemoryScheduleRepository {
	return &MemoryScheduleRepository{
		schedules: make(map[string]*entities.Schedule),
	}
}

func (r *MemoryScheduleRepository) FindByOwnerID(ownerID string) (*entities.Schedule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, schedule := range r.schedules {
		if schedule.OwnerID() == ownerID {
			return schedule, nil
		}
	}
	return nil, errors.New("schedule not found")
}

func (r *MemoryScheduleRepository) Save(schedule *entities.Schedule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.schedules[schedule.ID()] = schedule
	return nil
}

func (r *MemoryScheduleRepository) FindByID(id string) (*entities.Schedule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	schedule, exists := r.schedules[id]
	if !exists {
		return nil, errors.New("schedule not found")
	}
	return schedule, nil
}

func (r *MemoryScheduleRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.schedules[id]; !exists {
		return errors.New("schedule not found")
	}
	
	delete(r.schedules, id)
	return nil
}

func (r *MemoryScheduleRepository) FindAll() ([]*entities.Schedule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*entities.Schedule, 0, len(r.schedules))
	for _, schedule := range r.schedules {
		result = append(result, schedule)
	}
	return result, nil
}
