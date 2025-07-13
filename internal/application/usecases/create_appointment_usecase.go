package usecases

import (
	"errors"

	"github.com/visiab/appointment-calculator/internal/application/dto"
	"github.com/visiab/appointment-calculator/internal/domain/entities"
	"github.com/visiab/appointment-calculator/internal/domain/services"
	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type AppointmentRepository interface {
	Save(appointment *entities.Appointment) error
	FindByID(id string) (*entities.Appointment, error)
	FindByParticipant(participantID string) ([]*entities.Appointment, error)
	Update(appointment *entities.Appointment) error
	Delete(id string) error
}

type ScheduleRepository interface {
	FindByOwnerID(ownerID string) (*entities.Schedule, error)
	Save(schedule *entities.Schedule) error
}

type NotificationGateway interface {
	SendAppointmentCreated(appointment *entities.Appointment) error
	SendAppointmentUpdated(appointment *entities.Appointment) error
	SendAppointmentCancelled(appointment *entities.Appointment) error
}

type CreateAppointmentUseCase struct {
	appointmentRepo     AppointmentRepository
	scheduleRepo        ScheduleRepository
	notificationGateway NotificationGateway
	conflictDetector    *services.ConflictDetectionService
}

func NewCreateAppointmentUseCase(
	appointmentRepo AppointmentRepository,
	scheduleRepo ScheduleRepository,
	notificationGateway NotificationGateway,
	conflictDetector *services.ConflictDetectionService,
) *CreateAppointmentUseCase {
	return &CreateAppointmentUseCase{
		appointmentRepo:     appointmentRepo,
		scheduleRepo:        scheduleRepo,
		notificationGateway: notificationGateway,
		conflictDetector:    conflictDetector,
	}
}

func (uc *CreateAppointmentUseCase) Execute(request dto.CreateAppointmentRequest) (*dto.CreateAppointmentResponse, error) {
	// Validate time range
	timeRange, err := valueobjects.NewTimeRange(request.StartTime, request.EndTime)
	if err != nil {
		return nil, errors.New("invalid time range: " + err.Error())
	}

	// Create appointment entity
	appointment, err := entities.NewAppointment(request.Title, timeRange, request.Attendees, request.Location)
	if err != nil {
		return nil, errors.New("failed to create appointment: " + err.Error())
	}

	// Check for conflicts with each attendee's schedule
	for _, attendeeID := range request.Attendees {
		schedule, err := uc.scheduleRepo.FindByOwnerID(attendeeID)
		if err != nil {
			continue // Skip if schedule not found (participant might not have a schedule yet)
		}

		conflictResult := uc.conflictDetector.DetectConflicts(schedule, timeRange)
		if conflictResult.HasConflict {
			return nil, errors.New("appointment conflicts with existing schedule for participant " + attendeeID)
		}
	}

	// Save appointment
	err = uc.appointmentRepo.Save(appointment)
	if err != nil {
		return nil, errors.New("failed to save appointment: " + err.Error())
	}

	// Add appointment to participants' schedules
	for _, attendeeID := range request.Attendees {
		schedule, err := uc.scheduleRepo.FindByOwnerID(attendeeID)
		if err != nil {
			continue // Skip if schedule not found
		}

		err = schedule.AddAppointment(appointment)
		if err != nil {
			// Log error but don't fail the entire operation
			continue
		}

		err = uc.scheduleRepo.Save(schedule)
		if err != nil {
			// Log error but don't fail the entire operation
			continue
		}
	}

	// Send notification
	err = uc.notificationGateway.SendAppointmentCreated(appointment)
	if err != nil {
		// Log error but don't fail the operation
	}

	return &dto.CreateAppointmentResponse{
		ID:        appointment.ID(),
		Title:     appointment.Title(),
		StartTime: appointment.TimeRange().StartTime(),
		EndTime:   appointment.TimeRange().EndTime(),
		Attendees: appointment.Attendees(),
		Location:  appointment.Location(),
		Status:    appointment.Status(),
		CreatedAt: appointment.CreatedAt(),
	}, nil
}
