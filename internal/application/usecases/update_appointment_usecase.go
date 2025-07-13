package usecases

import (
	"errors"

	"github.com/visiab/appointment-calculator/internal/application/dto"
	"github.com/visiab/appointment-calculator/internal/domain/entities"
	"github.com/visiab/appointment-calculator/internal/domain/services"
	"github.com/visiab/appointment-calculator/internal/domain/valueobjects"
)

type UpdateAppointmentUseCase struct {
	appointmentRepo     AppointmentRepository
	scheduleRepo        ScheduleRepository
	notificationGateway NotificationGateway
	conflictDetector    *services.ConflictDetectionService
}

func NewUpdateAppointmentUseCase(
	appointmentRepo AppointmentRepository,
	scheduleRepo ScheduleRepository,
	notificationGateway NotificationGateway,
	conflictDetector *services.ConflictDetectionService,
) *UpdateAppointmentUseCase {
	return &UpdateAppointmentUseCase{
		appointmentRepo:     appointmentRepo,
		scheduleRepo:        scheduleRepo,
		notificationGateway: notificationGateway,
		conflictDetector:    conflictDetector,
	}
}

func (uc *UpdateAppointmentUseCase) Execute(appointmentID string, request dto.UpdateAppointmentRequest) (*dto.AppointmentResponse, error) {
	// Find existing appointment
	appointment, err := uc.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return nil, errors.New("appointment not found: " + err.Error())
	}

	// Check if appointment can be updated
	if appointment.Status() == entities.StatusCompleted {
		return nil, errors.New("cannot update completed appointment")
	}

	if appointment.Status() == entities.StatusCancelled {
		return nil, errors.New("cannot update cancelled appointment")
	}

	// Check if time is being updated
	if request.StartTime != nil || request.EndTime != nil {
		startTime := appointment.TimeRange().StartTime()
		endTime := appointment.TimeRange().EndTime()

		if request.StartTime != nil {
			startTime = *request.StartTime
		}
		if request.EndTime != nil {
			endTime = *request.EndTime
		}

		// Validate new time range
		newTimeRange, err := valueobjects.NewTimeRange(startTime, endTime)
		if err != nil {
			return nil, errors.New("invalid time range: " + err.Error())
		}

		// Check for conflicts with new time
		for _, attendeeID := range appointment.Attendees() {
			schedule, err := uc.scheduleRepo.FindByOwnerID(attendeeID)
			if err != nil {
				continue
			}

			// Temporarily remove this appointment from schedule to check conflicts
			schedule.RemoveAppointment(appointmentID)
			conflictResult := uc.conflictDetector.DetectConflicts(schedule, newTimeRange)
			// Add it back
			schedule.AddAppointment(appointment)

			if conflictResult.HasConflict {
				return nil, errors.New("updated time conflicts with existing schedule for participant " + attendeeID)
			}
		}

		// Update the appointment's time
		appointment.Reschedule(newTimeRange)
	}

	// Update other fields (this would require updating the appointment entity)
	// For now, we'll assume the appointment entity has update methods

	// Save updated appointment
	err = uc.appointmentRepo.Update(appointment)
	if err != nil {
		return nil, errors.New("failed to update appointment: " + err.Error())
	}

	// Send notification
	err = uc.notificationGateway.SendAppointmentUpdated(appointment)
	if err != nil {
		// Log error but don't fail the operation
	}

	return &dto.AppointmentResponse{
		ID:        appointment.ID(),
		Title:     appointment.Title(),
		StartTime: appointment.TimeRange().StartTime(),
		EndTime:   appointment.TimeRange().EndTime(),
		Duration:  appointment.TimeRange().Duration().String(),
		Attendees: appointment.Attendees(),
		Location:  appointment.Location(),
		Status:    appointment.Status(),
		CreatedAt: appointment.CreatedAt(),
		UpdatedAt: appointment.UpdatedAt(),
	}, nil
}

func (uc *UpdateAppointmentUseCase) Cancel(appointmentID string) error {
	// Find existing appointment
	appointment, err := uc.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return errors.New("appointment not found: " + err.Error())
	}

	// Check if appointment can be cancelled
	if appointment.Status() == entities.StatusCompleted {
		return errors.New("cannot cancel completed appointment")
	}

	if appointment.Status() == entities.StatusCancelled {
		return errors.New("appointment is already cancelled")
	}

	// Cancel the appointment
	appointment.Cancel()

	// Save updated appointment
	err = uc.appointmentRepo.Update(appointment)
	if err != nil {
		return errors.New("failed to cancel appointment: " + err.Error())
	}

	// Remove from participants' schedules
	for _, attendeeID := range appointment.Attendees() {
		schedule, err := uc.scheduleRepo.FindByOwnerID(attendeeID)
		if err != nil {
			continue
		}

		err = schedule.RemoveAppointment(appointmentID)
		if err != nil {
			continue
		}

		err = uc.scheduleRepo.Save(schedule)
		if err != nil {
			continue
		}
	}

	// Send notification
	err = uc.notificationGateway.SendAppointmentCancelled(appointment)
	if err != nil {
		// Log error but don't fail the operation
	}

	return nil
}
