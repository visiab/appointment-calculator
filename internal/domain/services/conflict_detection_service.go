package services

import (
	"github.com/visiab/appointment-calculator/internal/domain/entities"
	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type ConflictDetectionService struct{}

func NewConflictDetectionService() *ConflictDetectionService {
	return &ConflictDetectionService{}
}

type ConflictResult struct {
	HasConflict        bool
	ConflictingSlots   []ConflictingSlot
	ConflictType       ConflictType
	Severity          ConflictSeverity
}

type ConflictingSlot struct {
	AppointmentID string
	TimeRange     valueobjects.TimeRange
	OverlapRange  valueobjects.TimeRange
}

type ConflictType string

const (
	ConflictTypeAppointment ConflictType = "appointment"
	ConflictTypeBlocked     ConflictType = "blocked_time"
	ConflictTypeWorkingHours ConflictType = "working_hours"
)

type ConflictSeverity string

const (
	SeverityMinor    ConflictSeverity = "minor"     // Small overlap
	SeverityMajor    ConflictSeverity = "major"     // Significant overlap
	SeverityCritical ConflictSeverity = "critical"  // Complete overlap
)

func (s *ConflictDetectionService) DetectConflicts(schedule *entities.Schedule, proposedTimeRange valueobjects.TimeRange) ConflictResult {
	result := ConflictResult{
		HasConflict:      false,
		ConflictingSlots: make([]ConflictingSlot, 0),
	}

	// Check working hours conflict
	if !schedule.WorkingHours().Contains(proposedTimeRange) {
		result.HasConflict = true
		result.ConflictType = ConflictTypeWorkingHours
		result.Severity = SeverityCritical
		return result
	}

	// Check appointment conflicts
	for _, appointment := range schedule.Appointments() {
		if appointment.Status() == entities.StatusCancelled {
			continue
		}

		if appointment.TimeRange().OverlapsWith(proposedTimeRange) {
			result.HasConflict = true
			result.ConflictType = ConflictTypeAppointment
			
			overlapRange := s.calculateOverlap(appointment.TimeRange(), proposedTimeRange)
			conflictingSlot := ConflictingSlot{
				AppointmentID: appointment.ID(),
				TimeRange:     appointment.TimeRange(),
				OverlapRange:  overlapRange,
			}
			result.ConflictingSlots = append(result.ConflictingSlots, conflictingSlot)
		}
	}

	// Check blocked time conflicts
	for _, blockedTime := range schedule.BlockedTimes() {
		if blockedTime.OverlapsWith(proposedTimeRange) {
			result.HasConflict = true
			result.ConflictType = ConflictTypeBlocked
			
			overlapRange := s.calculateOverlap(blockedTime, proposedTimeRange)
			conflictingSlot := ConflictingSlot{
				AppointmentID: "blocked",
				TimeRange:     blockedTime,
				OverlapRange:  overlapRange,
			}
			result.ConflictingSlots = append(result.ConflictingSlots, conflictingSlot)
		}
	}

	// Determine severity based on overlap
	if result.HasConflict {
		result.Severity = s.calculateSeverity(proposedTimeRange, result.ConflictingSlots)
	}

	return result
}

func (s *ConflictDetectionService) DetectMultiParticipantConflicts(participants []*entities.Participant, proposedTimeRange valueobjects.TimeRange) map[string]ConflictResult {
	conflicts := make(map[string]ConflictResult)

	for _, participant := range participants {
		if !participant.IsAvailableAt(proposedTimeRange) {
			conflicts[participant.ID()] = ConflictResult{
				HasConflict:  true,
				ConflictType: ConflictTypeAppointment,
				Severity:     SeverityCritical,
			}
		}
	}

	return conflicts
}

func (s *ConflictDetectionService) calculateOverlap(timeRange1, timeRange2 valueobjects.TimeRange) valueobjects.TimeRange {
	start := timeRange1.StartTime()
	if timeRange2.StartTime().After(start) {
		start = timeRange2.StartTime()
	}

	end := timeRange1.EndTime()
	if timeRange2.EndTime().Before(end) {
		end = timeRange2.EndTime()
	}

	overlapRange, _ := valueobjects.NewTimeRange(start, end)
	return overlapRange
}

func (s *ConflictDetectionService) calculateSeverity(proposedTimeRange valueobjects.TimeRange, conflicts []ConflictingSlot) ConflictSeverity {
	proposedDuration := proposedTimeRange.Duration()
	maxOverlapDuration := proposedDuration * 0

	for _, conflict := range conflicts {
		overlapDuration := conflict.OverlapRange.Duration()
		if overlapDuration > maxOverlapDuration {
			maxOverlapDuration = overlapDuration
		}
	}

	overlapPercentage := float64(maxOverlapDuration) / float64(proposedDuration)

	if overlapPercentage >= 0.8 {
		return SeverityCritical
	} else if overlapPercentage >= 0.3 {
		return SeverityMajor
	}
	return SeverityMinor
}
