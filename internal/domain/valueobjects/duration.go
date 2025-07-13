package valueobjects

import (
	"errors"
	"time"
)

type Duration struct {
	value time.Duration
}

func NewDuration(d time.Duration) (Duration, error) {
	if d <= 0 {
		return Duration{}, errors.New("duration must be positive")
	}
	
	return Duration{value: d}, nil
}

func (d Duration) Value() time.Duration {
	return d.value
}

func (d Duration) Minutes() float64 {
	return d.value.Minutes()
}

func (d Duration) Hours() float64 {
	return d.value.Hours()
}

func (d Duration) Add(other Duration) Duration {
	return Duration{value: d.value + other.value}
}

func (d Duration) Subtract(other Duration) (Duration, error) {
	result := d.value - other.value
	if result <= 0 {
		return Duration{}, errors.New("resulting duration would be negative or zero")
	}
	return Duration{value: result}, nil
}
