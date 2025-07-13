package valueobjects

import (
	"errors"
	"time"
)

type TimeRange struct {
	startTime time.Time
	endTime   time.Time
}

func NewTimeRange(startTime, endTime time.Time) (TimeRange, error) {
	if endTime.Before(startTime) {
		return TimeRange{}, errors.New("end time cannot be before start time")
	}
	
	if startTime.Equal(endTime) {
		return TimeRange{}, errors.New("start time and end time cannot be equal")
	}
	
	return TimeRange{
		startTime: startTime,
		endTime:   endTime,
	}, nil
}

func (tr TimeRange) StartTime() time.Time {
	return tr.startTime
}

func (tr TimeRange) EndTime() time.Time {
	return tr.endTime
}

func (tr TimeRange) Duration() time.Duration {
	return tr.endTime.Sub(tr.startTime)
}

func (tr TimeRange) OverlapsWith(other TimeRange) bool {
	return tr.startTime.Before(other.endTime) && tr.endTime.After(other.startTime)
}

func (tr TimeRange) Contains(other TimeRange) bool {
	return !tr.startTime.After(other.startTime) && !tr.endTime.Before(other.endTime)
}

func (tr TimeRange) IsWithin(other TimeRange) bool {
	return !other.startTime.After(tr.startTime) && !other.endTime.Before(tr.endTime)
}
