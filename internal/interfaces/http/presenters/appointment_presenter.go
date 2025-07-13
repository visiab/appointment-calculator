package presenters

import (
	"github.com/visiab/appointment-calculator/internal/application/dto"
	"github.com/visiab/appointment-calculator/internal/domain/entities"
)

type AppointmentPresenter struct{}

func NewAppointmentPresenter() *AppointmentPresenter {
	return &AppointmentPresenter{}
}

func (p *AppointmentPresenter) PresentAppointment(appointment *entities.Appointment) dto.AppointmentResponse {
	return dto.AppointmentResponse{
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
	}
}

func (p *AppointmentPresenter) PresentAppointmentList(appointments []*entities.Appointment, total, page, limit int) dto.AppointmentListResponse {
	responses := make([]dto.AppointmentResponse, len(appointments))
	for i, appointment := range appointments {
		responses[i] = p.PresentAppointment(appointment)
	}

	return dto.AppointmentListResponse{
		Appointments: responses,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}
}

func (p *AppointmentPresenter) PresentCreateResponse(appointment *entities.Appointment) dto.CreateAppointmentResponse {
	return dto.CreateAppointmentResponse{
		ID:        appointment.ID(),
		Title:     appointment.Title(),
		StartTime: appointment.TimeRange().StartTime(),
		EndTime:   appointment.TimeRange().EndTime(),
		Attendees: appointment.Attendees(),
		Location:  appointment.Location(),
		Status:    appointment.Status(),
		CreatedAt: appointment.CreatedAt(),
	}
}
