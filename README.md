# Appointment Calculator

A Go-based appointment scheduling and availability calculation service built with Clean Architecture principles.

## Overview

The Appointment Calculator provides REST API endpoints to:
- Create and manage appointments
- Calculate optimal meeting times for multiple participants
- Handle participant availability and scheduling conflicts
- Support recurring appointments and timezone management

## Architecture

This application follows Clean Architecture principles with clear separation of concerns:

### Domain Layer
- **Entities**: Core business objects (Appointment, Participant, Schedule)
- **Value Objects**: Immutable objects (TimeRange, Duration, TimeSlot)
- **Domain Services**: Business logic services (ConflictDetection, OptimalTimeFinder, RecurrenceCalculator)

### Application Layer
- **Use Cases**: Application-specific business rules
- **DTOs**: Data transfer objects for API communication

### Interface Adapters
- **Controllers**: HTTP request handlers
- **Presenters**: Data formatting for responses
- **Repository Interfaces**: Data access contracts

### Infrastructure Layer
- **Repositories**: Concrete data access implementations
- **External Services**: Notification, timezone services
- **Configuration**: Environment-based configuration

## Getting Started

### Prerequisites
- Go 1.21 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/visiab/appointment-calculator.git
cd appointment-calculator
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

### Configuration

The application can be configured using environment variables:

```bash
# Server Configuration
PORT=8080
HOST=0.0.0.0
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
IDLE_TIMEOUT=120s

# Database Configuration
DB_TYPE=memory  # Currently only memory is supported

# Logging Configuration
LOG_LEVEL=info  # debug, info, warn, error
LOG_FORMAT=text # text, json
```

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Appointments
- `POST /api/v1/appointments` - Create a new appointment
- `GET /api/v1/appointments` - List appointments
- `GET /api/v1/appointments/{id}` - Get appointment details
- `PUT /api/v1/appointments/{id}` - Update appointment
- `DELETE /api/v1/appointments/{id}` - Cancel appointment

### Schedules
- `POST /api/v1/schedules/availability` - Find available time slots
- `GET /api/v1/schedules/{owner_id}/overview` - Get schedule overview
- `GET /api/v1/schedules/{owner_id}/detail` - Get detailed schedule
- `POST /api/v1/schedules/{owner_id}/blocked-times` - Add blocked time

### Participants
- `POST /api/v1/participants` - Create participant
- `GET /api/v1/participants/{id}` - Get participant details
- `PUT /api/v1/participants/{id}` - Update participant
- `POST /api/v1/participants/{id}/availability` - Add availability
- `GET /api/v1/participants/{id}/availability` - Get availability

## API Examples

### Create an Appointment

```bash
curl -X POST http://localhost:8080/api/v1/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "start_time": "2024-01-15T14:00:00Z",
    "end_time": "2024-01-15T15:00:00Z",
    "attendees": ["user1", "user2"],
    "location": "Conference Room A"
  }'
```

### Find Available Time Slots

```bash
curl -X POST http://localhost:8080/api/v1/schedules/availability \
  -H "Content-Type: application/json" \
  -d '{
    "participant_ids": ["user1", "user2"],
    "start_date": "2024-01-15T09:00:00Z",
    "end_date": "2024-01-15T17:00:00Z",
    "duration_minutes": 60,
    "timezone": "America/New_York"
  }'
```

## Development

### Project Structure

```
.
├── cmd/api/                    # Application entry point
├── internal/
│   ├── domain/                 # Domain layer
│   │   ├── entities/          # Business entities
│   │   ├── valueobjects/      # Value objects
│   │   └── services/          # Domain services
│   ├── application/           # Application layer
│   │   ├── dto/              # Data transfer objects
│   │   └── usecases/         # Use cases
│   ├── interfaces/           # Interface adapters
│   │   └── http/
│   │       ├── controllers/  # HTTP controllers
│   │       └── presenters/   # Response presenters
│   └── infrastructure/       # Infrastructure layer
│       ├── repositories/     # Data access implementations
│       ├── services/        # External service implementations
│       ├── config/          # Configuration
│       ├── dependency/      # Dependency injection
│       └── web/            # Web framework setup
├── go.mod
├── go.sum
├── README.md
└── PROJECT_DESIGN.md          # Detailed architecture documentation
```

### Building

```bash
# Build the application
go build -o bin/appointment-calculator cmd/api/main.go

# Run tests
go test ./...

# Run with race detection
go run -race cmd/api/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## Features

### Current Features
- ✅ Clean Architecture implementation
- ✅ RESTful API endpoints
- ✅ In-memory data storage
- ✅ Appointment management
- ✅ Conflict detection
- ✅ Optimal time finding
- ✅ Timezone support
- ✅ Graceful shutdown
- ✅ CORS support
- ✅ Configuration management

### Planned Features
- 🔄 Database persistence (PostgreSQL, MySQL)
- 🔄 Email notifications
- 🔄 Calendar integrations (Google Calendar, Outlook)
- 🔄 Recurring appointments
- 🔄 User authentication
- 🔄 Rate limiting
- 🔄 Metrics and monitoring
- 🔄 Docker support
- 🔄 API documentation (Swagger)

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Add tests for your changes
5. Ensure tests pass: `go test ./...`
6. Commit your changes: `git commit -am 'Add feature'`
7. Push to the branch: `git push origin feature-name`
8. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Design Documentation

For detailed architecture and design information, see [PROJECT_DESIGN.md](PROJECT_DESIGN.md).
