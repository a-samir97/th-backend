# Thamaniyah Media Platform

A scalable backend system for managing and discovering media content (videos and podcasts) built with Go.

## Architecture

The system follows a microservices architecture with clean separation of concerns:

- **CMS Service**: Handles media CRUD operations and upload coordination
- **Discovery Service**: Provides search functionality across media content  
- **Media Processor**: Background service for processing media metadata
- **Message Queue**: Asynchronous communication between services

## Features

### CMS Service
- Media upload workflow with presigned URLs
- CRUD operations for media metadata
- File storage abstraction (local/S3)
- Async processing pipeline

### Discovery Service  
- Full-text search across media metadata
- Elasticsearch integration
- Advanced filtering and sorting
- Caching for performance

### Media Processor
- Background metadata extraction
- Elasticsearch indexing
- Event-driven processing
- Error handling and retries

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 15
- **Search**: Elasticsearch 8.11
- **Cache**: Redis 7
- **Message Queue**: RabbitMQ 3
- **HTTP Framework**: Gin
- **ORM**: GORM

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose

### Quick Start

1. Clone the repository
2. Start dependencies:
   ```bash
   docker-compose up -d
   ```

3. Run database migrations:
   ```bash
   go run cmd/migrations/main.go
   ```

4. Start services:
   ```bash
   # CMS Service
   go run cmd/cms-service/main.go
   
   # Discovery Service  
   go run cmd/discovery-service/main.go
   
   # Media Processor
   go run cmd/media-processor/main.go
   ```

## API Documentation

### CMS Service (Port 8080)
- `POST /api/v1/media/upload-url` - Generate upload URL
- `POST /api/v1/media/{id}/confirm` - Confirm upload completion
- `GET /api/v1/media` - List media
- `GET /api/v1/media/{id}` - Get media details
- `PUT /api/v1/media/{id}` - Update media metadata
- `DELETE /api/v1/media/{id}` - Delete media

### Discovery Service (Port 8081)
- `GET /api/v1/search` - Search media content
- `GET /api/v1/search/suggest` - Get search suggestions

## Configuration

Environment variables can be used to configure the services. See `internal/config/config.go` for available options.

## Project Structure

```
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── domain/            # Business entities
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic
│   ├── handler/           # HTTP handlers
│   ├── config/            # Configuration
│   └── middleware/        # HTTP middleware
├── pkg/                   # Shared packages
├── migrations/            # Database migrations
└── docker-compose.yml     # Development environment
```

## Development

- Follow Go coding standards
- Write tests for business logic
- Use dependency injection for clean architecture
- Keep services loosely coupled via message queues
