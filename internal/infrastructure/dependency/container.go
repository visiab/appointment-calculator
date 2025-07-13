package dependency

import (
	"github.com/visiab/appointment-calculator/internal/application/usecases"
	"github.com/visiab/appointment-calculator/internal/domain/services"
	"github.com/visiab/appointment-calculator/internal/infrastructure/config"
	"github.com/visiab/appointment-calculator/internal/infrastructure/repositories"
	infraServices "github.com/visiab/appointment-calculator/internal/infrastructure/services"
	"github.com/visiab/appointment-calculator/internal/interfaces/http/controllers"
	"github.com/visiab/appointment-calculator/internal/interfaces/http/presenters"
)

type Container struct {
	Config *config.Config

	// Repositories
	AppointmentRepo  usecases.AppointmentRepository
	ScheduleRepo     usecases.ScheduleRepository
	ParticipantRepo  usecases.ParticipantRepository

	// Domain Services
	ConflictDetector    *services.ConflictDetectionService
	OptimalTimeFinder   *services.OptimalTimeFinderService
	RecurrenceCalculator *services.RecurrenceCalculatorService

	// Infrastructure Services
	NotificationGateway usecases.NotificationGateway
	TimezoneService     *infraServices.TimezoneService

	// Use Cases
	CreateAppointmentUseCase         *usecases.CreateAppointmentUseCase
	UpdateAppointmentUseCase         *usecases.UpdateAppointmentUseCase
	FindAvailableTimeSlotsUseCase    *usecases.FindAvailableTimeSlotsUseCase

	// Presenters
	AppointmentPresenter *presenters.AppointmentPresenter
	SchedulePresenter    *presenters.SchedulePresenter

	// Controllers
	AppointmentController *controllers.AppointmentController
	ScheduleController    *controllers.ScheduleController
	ParticipantController *controllers.ParticipantController
}

func NewContainer() *Container {
	cfg := config.Load()
	c := &Container{
		Config: cfg,
	}

	// Initialize repositories
	c.initRepositories()

	// Initialize domain services
	c.initDomainServices()

	// Initialize infrastructure services
	c.initInfrastructureServices()

	// Initialize use cases
	c.initUseCases()

	// Initialize presenters
	c.initPresenters()

	// Initialize controllers
	c.initControllers()

	return c
}

func (c *Container) initRepositories() {
	// For now, use in-memory repositories
	// In production, this would be configurable based on config.Database.Type
	c.AppointmentRepo = repositories.NewMemoryAppointmentRepository()
	c.ScheduleRepo = repositories.NewMemoryScheduleRepository()
	c.ParticipantRepo = repositories.NewMemoryParticipantRepository()
}

func (c *Container) initDomainServices() {
	c.ConflictDetector = services.NewConflictDetectionService()
	c.OptimalTimeFinder = services.NewOptimalTimeFinderService(c.ConflictDetector)
	c.RecurrenceCalculator = services.NewRecurrenceCalculatorService()
}

func (c *Container) initInfrastructureServices() {
	c.NotificationGateway = infraServices.NewConsoleNotificationService()
	c.TimezoneService = infraServices.NewTimezoneService()
}

func (c *Container) initUseCases() {
	c.CreateAppointmentUseCase = usecases.NewCreateAppointmentUseCase(
		c.AppointmentRepo,
		c.ScheduleRepo,
		c.NotificationGateway,
		c.ConflictDetector,
	)

	c.UpdateAppointmentUseCase = usecases.NewUpdateAppointmentUseCase(
		c.AppointmentRepo,
		c.ScheduleRepo,
		c.NotificationGateway,
		c.ConflictDetector,
	)

	c.FindAvailableTimeSlotsUseCase = usecases.NewFindAvailableTimeSlotsUseCase(
		c.ParticipantRepo,
		c.OptimalTimeFinder,
	)
}

func (c *Container) initPresenters() {
	c.AppointmentPresenter = presenters.NewAppointmentPresenter()
	c.SchedulePresenter = presenters.NewSchedulePresenter(c.AppointmentPresenter)
}

func (c *Container) initControllers() {
	c.AppointmentController = controllers.NewAppointmentController(
		c.CreateAppointmentUseCase,
		c.UpdateAppointmentUseCase,
	)

	c.ScheduleController = controllers.NewScheduleController(
		c.FindAvailableTimeSlotsUseCase,
	)

	c.ParticipantController = controllers.NewParticipantController()
}
