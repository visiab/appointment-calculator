package services

import (
	"errors"
	"time"
)

type TimezoneService struct {
	cache map[string]*time.Location
}

func NewTimezoneService() *TimezoneService {
	return &TimezoneService{
		cache: make(map[string]*time.Location),
	}
}

func (s *TimezoneService) GetLocation(timezone string) (*time.Location, error) {
	if timezone == "" {
		return time.UTC, nil
	}

	// Check cache first
	if location, exists := s.cache[timezone]; exists {
		return location, nil
	}

	// Load and cache the location
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, errors.New("invalid timezone: " + timezone)
	}

	s.cache[timezone] = location
	return location, nil
}

func (s *TimezoneService) ConvertTime(t time.Time, fromTz, toTz string) (time.Time, error) {
	fromLocation, err := s.GetLocation(fromTz)
	if err != nil {
		return time.Time{}, err
	}

	toLocation, err := s.GetLocation(toTz)
	if err != nil {
		return time.Time{}, err
	}

	// Convert to the source timezone first, then to target
	timeInFrom := t.In(fromLocation)
	return timeInFrom.In(toLocation), nil
}

func (s *TimezoneService) GetCommonTimezones() []string {
	return []string{
		"UTC",
		"America/New_York",
		"America/Los_Angeles",
		"America/Chicago",
		"America/Denver",
		"Europe/London",
		"Europe/Paris",
		"Europe/Berlin",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Asia/Mumbai",
		"Australia/Sydney",
		"Pacific/Auckland",
	}
}
