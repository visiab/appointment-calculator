package services

import (
	"sort"
	"time"

	"github.com/visiab/appointment-calculator/internal/domain/entities"
	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type OptimalTimeFinderService struct {
	conflictDetector *ConflictDetectionService
}

func NewOptimalTimeFinderService(conflictDetector *ConflictDetectionService) *OptimalTimeFinderService {
	return &OptimalTimeFinderService{
		conflictDetector: conflictDetector,
	}
}

type TimeOption struct {
	TimeRange     valueobjects.TimeRange
	Score         float64
	Reason        string
	Participants  []string
	Conflicts     int
}

type FindOptimalTimeRequest struct {
	Participants     []*entities.Participant
	Duration         valueobjects.Duration
	PreferredStart   time.Time
	PreferredEnd     time.Time
	EarliestStart    time.Time
	LatestEnd        time.Time
	TimeSlotInterval time.Duration
	MaxOptions       int
}

func (s *OptimalTimeFinderService) FindOptimalTimes(request FindOptimalTimeRequest) []TimeOption {
	options := make([]TimeOption, 0)

	// Generate time slots within the search range
	current := request.EarliestStart
	for current.Add(request.Duration.Value()).Before(request.LatestEnd) {
		endTime := current.Add(request.Duration.Value())
		timeRange, err := valueobjects.NewTimeRange(current, endTime)
		if err != nil {
			current = current.Add(request.TimeSlotInterval)
			continue
		}

		option := s.evaluateTimeOption(timeRange, request)
		if option.Score > 0 {
			options = append(options, option)
		}

		current = current.Add(request.TimeSlotInterval)
	}

	// Sort by score (highest first)
	sort.Slice(options, func(i, j int) bool {
		return options[i].Score > options[j].Score
	})

	// Limit results
	if len(options) > request.MaxOptions {
		options = options[:request.MaxOptions]
	}

	return options
}

func (s *OptimalTimeFinderService) evaluateTimeOption(timeRange valueobjects.TimeRange, request FindOptimalTimeRequest) TimeOption {
	option := TimeOption{
		TimeRange: timeRange,
		Score:     0,
		Conflicts: 0,
	}

	availableParticipants := make([]string, 0)
	conflictCount := 0

	// Check each participant's availability
	for _, participant := range request.Participants {
		if participant.IsAvailableAt(timeRange) {
			availableParticipants = append(availableParticipants, participant.ID())
		} else {
			conflictCount++
		}
	}

	option.Participants = availableParticipants
	option.Conflicts = conflictCount

	// Calculate base score based on participant availability
	participantScore := float64(len(availableParticipants)) / float64(len(request.Participants)) * 100

	// Bonus for preferred time range
	preferredScore := 0.0
	if !request.PreferredStart.IsZero() && !request.PreferredEnd.IsZero() {
		preferredRange, err := valueobjects.NewTimeRange(request.PreferredStart, request.PreferredEnd)
		if err == nil && preferredRange.Contains(timeRange) {
			preferredScore = 20.0
			option.Reason = "Within preferred time range"
		} else if preferredRange.OverlapsWith(timeRange) {
			preferredScore = 10.0
			option.Reason = "Partially overlaps preferred time"
		}
	}

	// Business hours bonus (9 AM - 5 PM)
	businessHoursScore := 0.0
	startHour := timeRange.StartTime().Hour()
	endHour := timeRange.EndTime().Hour()
	if startHour >= 9 && endHour <= 17 {
		businessHoursScore = 15.0
		if option.Reason == "" {
			option.Reason = "During business hours"
		}
	}

	// Penalty for conflicts
	conflictPenalty := float64(conflictCount) * 25.0

	option.Score = participantScore + preferredScore + businessHoursScore - conflictPenalty

	// Don't return options with too many conflicts
	if float64(conflictCount)/float64(len(request.Participants)) > 0.5 {
		option.Score = 0
		option.Reason = "Too many conflicts"
	}

	return option
}

func (s *OptimalTimeFinderService) FindNextAvailableSlot(schedule *entities.Schedule, duration valueobjects.Duration, after time.Time) (valueobjects.TimeRange, bool) {
	current := after
	interval := 15 * time.Minute // 15-minute intervals
	maxSearch := 30 * 24 * time.Hour // Search for up to 30 days

	for elapsed := time.Duration(0); elapsed < maxSearch; elapsed += interval {
		endTime := current.Add(duration.Value())
		timeRange, err := valueobjects.NewTimeRange(current, endTime)
		if err != nil {
			current = current.Add(interval)
			continue
		}

		conflictResult := s.conflictDetector.DetectConflicts(schedule, timeRange)
		if !conflictResult.HasConflict {
			return timeRange, true
		}

		current = current.Add(interval)
	}

	return valueobjects.TimeRange{}, false
}
