package presenters

import (
	"time"

	"github.com/visiab/appointment-calculator/internal/application/dto"
	"github.com/visiab/appointment-calculator/internal/domain/entities"
)

type SchedulePresenter struct {
	appointmentPresenter *AppointmentPresenter
}

func NewSchedulePresenter(appointmentPresenter *AppointmentPresenter) *SchedulePresenter {
	return &SchedulePresenter{
		appointmentPresenter: appointmentPresenter,
	}
}

func (p *SchedulePresenter) PresentScheduleOverview(schedule *entities.Schedule) dto.ScheduleOverview {
	appointments := schedule.Appointments()
	todayAppointments := p.countTodayAppointments(appointments)
	nextAppointment := p.findNextAppointment(appointments)

	var nextAppointmentResponse *dto.AppointmentResponse
	if nextAppointment != nil {
		response := p.appointmentPresenter.PresentAppointment(nextAppointment)
		nextAppointmentResponse = &response
	}

	return dto.ScheduleOverview{
		OwnerID:           schedule.OwnerID(),
		Timezone:          schedule.Timezone().String(),
		WorkingHoursStart: p.formatTime(schedule.WorkingHours().StartTime()),
		WorkingHoursEnd:   p.formatTime(schedule.WorkingHours().EndTime()),
		AppointmentsToday: todayAppointments,
		NextAppointment:   nextAppointmentResponse,
		TotalAppointments: len(appointments),
		LastUpdated:       time.Now(),
	}
}

func (p *SchedulePresenter) PresentScheduleDetail(schedule *entities.Schedule) dto.ScheduleDetail {
	overview := p.PresentScheduleOverview(schedule)
	appointmentResponses := make([]dto.AppointmentResponse, len(schedule.Appointments()))
	
	for i, appointment := range schedule.Appointments() {
		appointmentResponses[i] = p.appointmentPresenter.PresentAppointment(appointment)
	}

	blockedTimes := make([]struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		Reason    string    `json:"reason"`
	}, len(schedule.BlockedTimes()))

	for i, blockedTime := range schedule.BlockedTimes() {
		blockedTimes[i] = struct {
			StartTime time.Time `json:"start_time"`
			EndTime   time.Time `json:"end_time"`
			Reason    string    `json:"reason"`
		}{
			StartTime: blockedTime.StartTime(),
			EndTime:   blockedTime.EndTime(),
			Reason:    "Blocked time", // This could be enhanced with actual reasons
		}
	}

	return dto.ScheduleDetail{
		ScheduleOverview: overview,
		Appointments:     appointmentResponses,
		BlockedTimes:     blockedTimes,
	}
}

func (p *SchedulePresenter) countTodayAppointments(appointments []*entities.Appointment) int {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	count := 0

	for _, appointment := range appointments {
		if appointment.Status() == entities.StatusCancelled {
			continue
		}
		appointmentDate := appointment.TimeRange().StartTime().Truncate(24 * time.Hour)
		if appointmentDate.Equal(today) && appointmentDate.Before(tomorrow) {
			count++
		}
	}

	return count
}

func (p *SchedulePresenter) findNextAppointment(appointments []*entities.Appointment) *entities.Appointment {
	now := time.Now()
	var nextAppointment *entities.Appointment

	for _, appointment := range appointments {
		if appointment.Status() == entities.StatusCancelled {
			continue
		}
		if appointment.TimeRange().StartTime().After(now) {
			if nextAppointment == nil || appointment.TimeRange().StartTime().Before(nextAppointment.TimeRange().StartTime()) {
				nextAppointment = appointment
			}
		}
	}

	return nextAppointment
}

func (p *SchedulePresenter) formatTime(t time.Time) string {
	return t.Format("15:04")
}
