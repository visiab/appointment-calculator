package usecases

import (
	"errors"
	"time"

	"github.com/visiab/appointment-calculator/internal/application/dto"
	"github.com/visiab/appointment-calculator/internal/domain/services"
	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type ParticipantRepository interface {
	FindByID(id string) (*entities.Participant, error)
	FindByIDs(ids []string) ([]*entities.Participant, error)
	Save(participant *entities.Participant) error
	FindByEmail(email string) (*entities.Participant, error)
}

type FindAvailableTimeSlotsUseCase struct {
	participantRepo      ParticipantRepository
	optimalTimeFinder    *services.OptimalTimeFinderService
}

func NewFindAvailableTimeSlotsUseCase(
	participantRepo ParticipantRepository,
	optimalTimeFinder *services.OptimalTimeFinderService,
) *FindAvailableTimeSlotsUseCase {
	return &FindAvailableTimeSlotsUseCase{
		participantRepo:   participantRepo,
		optimalTimeFinder: optimalTimeFinder,
	}
}

func (uc *FindAvailableTimeSlotsUseCase) Execute(query dto.AvailabilityQuery) (*dto.AvailabilityResult, error) {
	// Validate input
	if query.StartDate.After(query.EndDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	if query.Duration <= 0 {
		return nil, errors.New("duration must be positive")
	}

	// Get participants
	participants, err := uc.participantRepo.FindByIDs(query.ParticipantIDs)
	if err != nil {
		return nil, errors.New("failed to find participants: " + err.Error())
	}

	if len(participants) == 0 {
		return nil, errors.New("no participants found")
	}

	// Create duration value object
	duration, err := valueobjects.NewDuration(time.Duration(query.Duration) * time.Minute)
	if err != nil {
		return nil, errors.New("invalid duration: " + err.Error())
	}

	// Parse timezone
	timezone := time.UTC
	if query.Timezone != "" {
		parsedTz, err := time.LoadLocation(query.Timezone)
		if err == nil {
			timezone = parsedTz
		}
	}

	// Adjust query times to specified timezone
	startTime := query.StartDate.In(timezone)
	endTime := query.EndDate.In(timezone)

	// Create request for optimal time finder
	request := services.FindOptimalTimeRequest{
		Participants:     participants,
		Duration:         duration,
		EarliestStart:    startTime,
		LatestEnd:        endTime,
		TimeSlotInterval: 15 * time.Minute, // 15-minute intervals
		MaxOptions:       50,               // Limit to 50 options
	}

	// Find optimal times
	timeOptions := uc.optimalTimeFinder.FindOptimalTimes(request)

	// Convert to response format
	availableSlots := make([]dto.TimeSlotResponse, len(timeOptions))
	for i, option := range timeOptions {
		availableSlots[i] = dto.TimeSlotResponse{
			StartTime:            option.TimeRange.StartTime(),
			EndTime:              option.TimeRange.EndTime(),
			Score:                option.Score,
			AvailableParticipants: len(option.Participants),
			TotalParticipants:    len(participants),
			Reason:               option.Reason,
			Conflicts:            option.Conflicts,
		}
	}

	// Prepare result
	result := &dto.AvailabilityResult{
		AvailableSlots: availableSlots,
	}

	// Set optimal slot (first one with highest score)
	if len(availableSlots) > 0 {
		result.OptimalSlot = &availableSlots[0]
	}

	// Set search period info
	result.SearchPeriod.StartDate = query.StartDate
	result.SearchPeriod.EndDate = query.EndDate
	result.SearchPeriod.Timezone = query.Timezone

	// Set summary
	result.Summary.TotalSlotsFound = len(availableSlots)
	result.Summary.Participants = len(participants)

	return result, nil
}
