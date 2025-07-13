package services

import (
	"fmt"
	"log"

	"github.com/visiab/appointment-calculator/internal/domain/entities"
)

type ConsoleNotificationService struct{}

func NewConsoleNotificationService() *ConsoleNotificationService {
	return &ConsoleNotificationService{}
}

func (s *ConsoleNotificationService) SendAppointmentCreated(appointment *entities.Appointment) error {
	message := fmt.Sprintf(
		"[NOTIFICATION] Appointment Created: %s (%s - %s) with attendees: %v",
		appointment.Title(),
		appointment.TimeRange().StartTime().Format("2006-01-02 15:04"),
		appointment.TimeRange().EndTime().Format("2006-01-02 15:04"),
		appointment.Attendees(),
	)
	log.Println(message)
	return nil
}

func (s *ConsoleNotificationService) SendAppointmentUpdated(appointment *entities.Appointment) error {
	message := fmt.Sprintf(
		"[NOTIFICATION] Appointment Updated: %s (%s - %s) with attendees: %v",
		appointment.Title(),
		appointment.TimeRange().StartTime().Format("2006-01-02 15:04"),
		appointment.TimeRange().EndTime().Format("2006-01-02 15:04"),
		appointment.Attendees(),
	)
	log.Println(message)
	return nil
}

func (s *ConsoleNotificationService) SendAppointmentCancelled(appointment *entities.Appointment) error {
	message := fmt.Sprintf(
		"[NOTIFICATION] Appointment Cancelled: %s (%s - %s) with attendees: %v",
		appointment.Title(),
		appointment.TimeRange().StartTime().Format("2006-01-02 15:04"),
		appointment.TimeRange().EndTime().Format("2006-01-02 15:04"),
		appointment.Attendees(),
	)
	log.Println(message)
	return nil
}
