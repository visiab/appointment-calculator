package services

import (
	"errors"
	"time"

	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type RecurrenceCalculatorService struct{}

func NewRecurrenceCalculatorService() *RecurrenceCalculatorService {
	return &RecurrenceCalculatorService{}
}

type RecurrencePattern string

const (
	PatternDaily   RecurrencePattern = "daily"
	PatternWeekly  RecurrencePattern = "weekly"
	PatternMonthly RecurrencePattern = "monthly"
	PatternYearly  RecurrencePattern = "yearly"
	PatternCustom  RecurrencePattern = "custom"
)

type RecurrenceRule struct {
	Pattern    RecurrencePattern
	Interval   int               // Every N days/weeks/months
	DaysOfWeek []time.Weekday    // For weekly patterns
	DayOfMonth int               // For monthly patterns
	EndDate    *time.Time        // When to stop recurring
	MaxCount   int               // Maximum number of occurrences
}

type RecurrenceResult struct {
	TimeRanges []valueobjects.TimeRange
	Count      int
	NextDate   *time.Time
}

func (s *RecurrenceCalculatorService) CalculateRecurrences(baseTimeRange valueobjects.TimeRange, rule RecurrenceRule) (RecurrenceResult, error) {
	if rule.Interval <= 0 {
		return RecurrenceResult{}, errors.New("recurrence interval must be positive")
	}

	result := RecurrenceResult{
		TimeRanges: make([]valueobjects.TimeRange, 0),
	}

	// Add the base time range
	result.TimeRanges = append(result.TimeRanges, baseTimeRange)
	result.Count = 1

	current := baseTimeRange.StartTime()
	duration := baseTimeRange.Duration()

	for {
		// Calculate next occurrence
		next := s.calculateNextOccurrence(current, rule)
		if next.IsZero() {
			break
		}

		// Check end conditions
		if rule.EndDate != nil && next.After(*rule.EndDate) {
			break
		}

		if rule.MaxCount > 0 && result.Count >= rule.MaxCount {
			result.NextDate = &next
			break
		}

		// Create time range for this occurrence
		endTime := next.Add(duration)
		timeRange, err := valueobjects.NewTimeRange(next, endTime)
		if err != nil {
			break
		}

		result.TimeRanges = append(result.TimeRanges, timeRange)
		result.Count++
		current = next
	}

	return result, nil
}

func (s *RecurrenceCalculatorService) calculateNextOccurrence(current time.Time, rule RecurrenceRule) time.Time {
	switch rule.Pattern {
	case PatternDaily:
		return current.AddDate(0, 0, rule.Interval)

	case PatternWeekly:
		if len(rule.DaysOfWeek) == 0 {
			return current.AddDate(0, 0, 7*rule.Interval)
		}
		return s.findNextWeekday(current, rule.DaysOfWeek, rule.Interval)

	case PatternMonthly:
		if rule.DayOfMonth > 0 {
			return s.findNextMonthlyDate(current, rule.DayOfMonth, rule.Interval)
		}
		return current.AddDate(0, rule.Interval, 0)

	case PatternYearly:
		return current.AddDate(rule.Interval, 0, 0)

	default:
		return time.Time{}
	}
}

func (s *RecurrenceCalculatorService) findNextWeekday(current time.Time, daysOfWeek []time.Weekday, interval int) time.Time {
	currentWeekday := current.Weekday()

	// Find next valid weekday in current week
	for _, day := range daysOfWeek {
		if day > currentWeekday {
			daysUntil := int(day - currentWeekday)
			return current.AddDate(0, 0, daysUntil)
		}
	}

	// No valid day in current week, go to next week(s)
	weeksToAdd := interval
	firstDayOfWeek := daysOfWeek[0]
	daysUntil := int(firstDayOfWeek) + 7*weeksToAdd - int(currentWeekday)
	if daysUntil <= 0 {
		daysUntil += 7
	}

	return current.AddDate(0, 0, daysUntil)
}

func (s *RecurrenceCalculatorService) findNextMonthlyDate(current time.Time, dayOfMonth, interval int) time.Time {
	nextMonth := current.AddDate(0, interval, 0)
	year, month, _ := nextMonth.Date()

	// Handle end-of-month cases
	lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, current.Location()).Day()
	if dayOfMonth > lastDayOfMonth {
		dayOfMonth = lastDayOfMonth
	}

	return time.Date(year, month, dayOfMonth, current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), current.Location())
}

func (s *RecurrenceCalculatorService) ValidateRecurrenceRule(rule RecurrenceRule) error {
	if rule.Interval <= 0 {
		return errors.New("interval must be positive")
	}

	if rule.Pattern == PatternWeekly && len(rule.DaysOfWeek) == 0 {
		return errors.New("weekly pattern requires at least one day of week")
	}

	if rule.Pattern == PatternMonthly && rule.DayOfMonth < 1 {
		return errors.New("monthly pattern requires valid day of month")
	}

	if rule.EndDate != nil && rule.MaxCount > 0 {
		return errors.New("cannot specify both end date and max count")
	}

	return nil
}
