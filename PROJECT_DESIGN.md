# Appointment Calculator - Clean Architecture Design

## Overview

The Appointment Calculator is designed following Clean Architecture principles, ensuring separation of concerns, testability, and maintainability. The application calculates optimal appointment schedules, manages availability, and handles appointment conflicts.

## Clean Architecture Layers

### 1. Domain Layer (Core Business Logic)

The innermost layer containing the business rules and entities that are independent of any external concerns.

#### Entities
**Responsibility**: Core business objects that encapsulate enterprise-wide business rules.

- **Appointment**: Represents a scheduled appointment
  - Properties: id, title, startTime, endTime, attendees, location, status
  - Business rules: Validates time constraints, duration limits, overlap prevention

- **TimeSlot**: Represents available time periods
  - Properties: startTime, endTime, isAvailable, resourceId
  - Business rules: Ensures valid time ranges, prevents negative durations

- **Participant**: Represents appointment attendees
  - Properties: id, name, email, timezone, availability
  - Business rules: Manages participant constraints and preferences

- **Schedule**: Represents a collection of appointments and availability
  - Properties: ownerId, timeZone, workingHours, appointments, blockedTimes
  - Business rules: Enforces working hours, manages capacity limits

#### Value Objects
**Responsibility**: Immutable objects that represent descriptive aspects without identity.

- **Duration**: Represents time spans with validation
- **TimeRange**: Represents start/end time pairs with overlap detection
- **Availability**: Complex availability patterns with recurring rules

#### Domain Services
**Responsibility**: Encapsulate domain logic that doesn't belong to a specific entity.

- **ConflictDetectionService**: Identifies scheduling conflicts
- **OptimalTimeFinderService**: Calculates best appointment times
- **RecurrenceCalculatorService**: Handles recurring appointment patterns

### 2. Application Layer (Use Cases)

Contains application-specific business rules and orchestrates the flow of data between layers.

#### Use Cases
**Responsibility**: Implement specific application business rules and coordinate between domain entities.

- **CreateAppointmentUseCase**
  - Validates appointment data
  - Checks for conflicts
  - Persists appointment
  - Sends notifications

- **FindAvailableTimeSlotsUseCase**
  - Analyzes participant availability
  - Considers duration requirements
  - Returns optimal time options
  - Accounts for timezone differences

- **UpdateAppointmentUseCase**
  - Validates changes
  - Checks new conflicts
  - Updates related data
  - Triggers notifications

- **CalculateOptimalMeetingTimeUseCase**
  - Analyzes multiple participant schedules
  - Considers preferences and constraints
  - Suggests best meeting times
  - Handles timezone coordination

#### Application Services
**Responsibility**: Coordinate use cases and manage application state.

- **AppointmentApplicationService**: Orchestrates appointment-related operations
- **ScheduleApplicationService**: Manages schedule operations and calculations
- **NotificationApplicationService**: Handles appointment notifications and reminders

#### DTOs (Data Transfer Objects)
**Responsibility**: Transfer data between layers without exposing domain entities.

- **CreateAppointmentRequest/Response**
- **AvailabilityQuery/Result**
- **ScheduleOverview/Detail**

### 3. Interface Adapters Layer

Converts data between the format most convenient for use cases and entities, and the format most convenient for external agencies.

#### Controllers
**Responsibility**: Handle HTTP requests and convert them to use case calls.

- **AppointmentController**: REST endpoints for appointment management
- **ScheduleController**: Endpoints for schedule operations
- **AvailabilityController**: Endpoints for availability queries

#### Presenters
**Responsibility**: Format data for presentation layer consumption.

- **AppointmentPresenter**: Formats appointment data for different views
- **SchedulePresenter**: Prepares schedule visualizations
- **CalendarPresenter**: Formats data for calendar components

#### Gateways/Repository Interfaces
**Responsibility**: Define contracts for data access without implementation details.

- **AppointmentRepository**: CRUD operations for appointments
- **ParticipantRepository**: Participant data management
- **ScheduleRepository**: Schedule persistence operations
- **NotificationGateway**: External notification service interface

### 4. Infrastructure Layer (External Concerns)

Contains implementations of interfaces defined in inner layers and handles external system communication.

#### Repository Implementations
**Responsibility**: Concrete implementations of repository interfaces.

- **DatabaseAppointmentRepository**: SQL/NoSQL database implementation
- **CacheAppointmentRepository**: Redis/memory cache implementation
- **FileScheduleRepository**: File-based storage for testing

#### External Service Adapters
**Responsibility**: Integrate with external systems and services.

- **EmailNotificationService**: Email service integration
- **CalendarSyncService**: Integration with external calendars (Google, Outlook)
- **TimezoneService**: Timezone data and conversion service
- **ConferenceRoomBookingService**: Physical/virtual room booking

#### Framework Components
**Responsibility**: Handle framework-specific concerns and configurations.

- **WebConfiguration**: HTTP server and routing setup
- **DatabaseConfiguration**: Database connection and migration management
- **SecurityConfiguration**: Authentication and authorization setup

## Component Characteristics

### Domain Layer Characteristics
- **Independence**: No dependencies on external frameworks or libraries
- **Stability**: Changes least frequently
- **Testability**: Pure business logic, easily unit tested
- **Reusability**: Can be used across different applications

### Application Layer Characteristics
- **Orchestration**: Coordinates domain entities and services
- **Transaction Management**: Handles business transaction boundaries
- **Validation**: Enforces application-specific rules
- **Error Handling**: Manages application-level exceptions

### Interface Adapters Characteristics
- **Translation**: Converts between different data formats
- **Protocol Handling**: Manages communication protocols (HTTP, gRPC)
- **Presentation Logic**: Handles view-specific formatting
- **Input Validation**: Validates external input format and structure

### Infrastructure Characteristics
- **Implementation Details**: Contains concrete implementations
- **External Dependencies**: Manages third-party integrations
- **Configuration**: Handles environment-specific settings
- **Performance Optimization**: Implements caching, connection pooling

## Key Design Principles

### Dependency Inversion
- High-level modules don't depend on low-level modules
- Both depend on abstractions (interfaces)
- Infrastructure implements interfaces defined in inner layers

### Single Responsibility
- Each class/module has one reason to change
- Clear separation of concerns across layers
- Focused, cohesive components

### Open/Closed Principle
- Components are open for extension, closed for modification
- New features added through new implementations
- Existing code remains stable

### Interface Segregation
- Clients depend only on interfaces they use
- Small, focused interfaces over large, monolithic ones
- Reduces coupling between components

## Testing Strategy

### Unit Testing
- **Domain Layer**: Test business rules in isolation
- **Application Layer**: Test use case orchestration with mocks
- **Infrastructure**: Test concrete implementations with test doubles

### Integration Testing
- **API Layer**: Test controller endpoints with real use cases
- **Database Layer**: Test repository implementations with test database
- **External Services**: Test adapters with service mocks

### Acceptance Testing
- **End-to-End**: Test complete user scenarios
- **Contract Testing**: Verify external service contracts
- **Performance Testing**: Validate system performance characteristics

## Benefits of This Architecture

1. **Maintainability**: Clear separation makes changes predictable and localized
2. **Testability**: Each layer can be tested independently with appropriate test doubles
3. **Flexibility**: External concerns can be swapped without affecting business logic
4. **Scalability**: Components can be optimized or replaced independently
5. **Team Collaboration**: Different teams can work on different layers simultaneously
6. **Technology Independence**: Business logic is not tied to specific frameworks or databases