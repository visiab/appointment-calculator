package valueobjects

import (
	"github.com/google/uuid"
)

type TimeSlot struct {
	id         string
	timeRange  TimeRange
	isAvailable bool
	resourceID string
}

func NewTimeSlot(timeRange TimeRange, isAvailable bool, resourceID string) TimeSlot {
	return TimeSlot{
		id:         uuid.New().String(),
		timeRange:  timeRange,
		isAvailable: isAvailable,
		resourceID: resourceID,
	}
}

func (ts TimeSlot) ID() string {
	return ts.id
}

func (ts TimeSlot) TimeRange() TimeRange {
	return ts.timeRange
}

func (ts TimeSlot) IsAvailable() bool {
	return ts.isAvailable
}

func (ts TimeSlot) ResourceID() string {
	return ts.resourceID
}

func (ts TimeSlot) Contains(timeRange TimeRange) bool {
	return ts.isAvailable && ts.timeRange.Contains(timeRange)
}

func (ts *TimeSlot) MarkAsUnavailable() {
	ts.isAvailable = false
}

func (ts *TimeSlot) MarkAsAvailable() {
	ts.isAvailable = true
}
