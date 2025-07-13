package repositories

import (
	"errors"
	"sync"

	"github.com/visiab/appointment-calculator/internal/domain/entities"
)

type MemoryParticipantRepository struct {
	participants map[string]*entities.Participant
	mu           sync.RWMutex
}

func NewMemoryParticipantRepository() *MemoryParticipantRepository {
	return &MemoryParticipantRepository{
		participants: make(map[string]*entities.Participant),
	}
}

func (r *MemoryParticipantRepository) FindByID(id string) (*entities.Participant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	participant, exists := r.participants[id]
	if !exists {
		return nil, errors.New("participant not found")
	}
	return participant, nil
}

func (r *MemoryParticipantRepository) FindByIDs(ids []string) ([]*entities.Participant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*entities.Participant, 0, len(ids))
	for _, id := range ids {
		if participant, exists := r.participants[id]; exists {
			result = append(result, participant)
		}
	}
	return result, nil
}

func (r *MemoryParticipantRepository) Save(participant *entities.Participant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.participants[participant.ID()] = participant
	return nil
}

func (r *MemoryParticipantRepository) FindByEmail(email string) (*entities.Participant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, participant := range r.participants {
		if participant.Email() == email {
			return participant, nil
		}
	}
	return nil, errors.New("participant not found")
}

func (r *MemoryParticipantRepository) Update(participant *entities.Participant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.participants[participant.ID()]; !exists {
		return errors.New("participant not found")
	}
	
	r.participants[participant.ID()] = participant
	return nil
}

func (r *MemoryParticipantRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.participants[id]; !exists {
		return errors.New("participant not found")
	}
	
	delete(r.participants, id)
	return nil
}

func (r *MemoryParticipantRepository) FindAll() ([]*entities.Participant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*entities.Participant, 0, len(r.participants))
	for _, participant := range r.participants {
		result = append(result, participant)
	}
	return result, nil
}
