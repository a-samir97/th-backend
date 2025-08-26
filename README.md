# Thamaniyah Media Platform 🎥

A modern, scalable media platform built with Go, featuring microservices architecture, advanced search capabilities, and enterprise-grade technologies.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [API Documentation](#api-documentation)
- [Usage Examples](#usage-examples)
- [Database Schema](#database-schema)
- [Development Guide](#development-guide)
- [Testing](#testing)
- [Deployment](#deployment)
- [Performance](#performance)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

Thamaniyah is a production-ready media platform designed to handle video and audio content at scale. It provides a complete solution for media upload, storage, processing, and discovery through advanced search capabilities.

### Key Capabilities
- **Media Management**: Upload, store, and manage video/audio files
- **Advanced Search**: Elasticsearch-powered full-text search with relevance scoring
- **Microservices**: Clean separation between CMS and Discovery services
- **Scalability**: Horizontal scaling with containerized deployment
- **Production Ready**: Comprehensive error handling, logging, and monitoring

## 🏗️ Architecture

The platform follows a **microservices architecture** with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐
│   CMS Service   │    │Discovery Service│
│    (Port 8080)  │    │   (Port 8081)   │
│                 │    │                 │
│ • Media CRUD    │    │ • Search        │
│ • File Upload   │    │ • Suggestions   │
│ • Metadata      │    │ • Indexing      │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────┬───────────────┘
                 │
    ┌─────────────────────────────┐
    │     Infrastructure Layer     │
    │                             │
    │ • PostgreSQL (Metadata)     │
    │ • Elasticsearch (Search)    │
    │ • Redis (Caching)           │
    │ • Docker (Containers)       │
    └─────────────────────────────┘
```

### Architecture Principles
1. **Clean Architecture**: Domain-driven design with clear layer separation
2. **SOLID Principles**: Single responsibility, dependency inversion
3. **Repository Pattern**: Data access abstraction
4. **Service Layer**: Business logic encapsulation
5. **HTTP Communication**: RESTful APIs between services

## ✨ Features

### 🎥 Media Management (CMS Service)
- ✅ **File Upload**: Generate presigned URLs for direct S3-style uploads
- ✅ **CRUD Operations**: Create, read, update, delete media records
- ✅ **Metadata Extraction**: Automatic duration, format, and size detection
- ✅ **Status Tracking**: Upload, processing, ready, failed states
- ✅ **Pagination**: Efficient large dataset handling

### 🔍 Advanced Search (Discovery Service)
- ✅ **Full-Text Search**: Elasticsearch-powered search across title, description
- ✅ **Relevance Scoring**: Intelligent ranking with field boosting
- ✅ **Autocomplete**: Real-time search suggestions
- ✅ **Type Filtering**: Filter by video, podcast, or other media types
- ✅ **Bulk Indexing**: Efficient reindexing of large datasets

### 🛠️ Infrastructure Features
- ✅ **Health Checks**: Service availability monitoring
- ✅ **CORS Support**: Cross-origin request handling
- ✅ **Graceful Shutdown**: Clean service termination
- ✅ **Structured Logging**: Request/response logging with timestamps
- ✅ **Error Handling**: Comprehensive error responses

## 🚀 Technology Stack

### Backend Services
- **Language**: Go 1.21+
- **HTTP Framework**: Gin (high-performance HTTP web framework)
- **Database ORM**: GORM (The fantastic ORM library for Golang)

### Infrastructure
- **Primary Database**: PostgreSQL 15+ (metadata, relationships)
- **Search Engine**: Elasticsearch 8.11+ (full-text search, indexing)
- **Caching**: Redis 7+ (session management, caching)
- **Message Queue**: RabbitMQ 3+ (async processing) *[Future]*
- **Containerization**: Docker & Docker Compose

### Architecture Patterns
- **Clean Architecture**: Domain, service, repository, handler layers
- **Repository Pattern**: Data access abstraction
- **Dependency Injection**: Interface-based dependency management
- **Service Layer Pattern**: Business logic separation

## 📁 Project Structure

```
thamaniyah/
├── cmd/                        # Application entry points
│   ├── cms-service/           # CMS service main
│   ├── discovery-service/     # Discovery service main
│   ├── migrate/              # Database migration tool
│   └── utils/                # Utility commands
│
├── internal/                  # Private application code
│   ├── config/               # Configuration management
│   ├── domain/               # Business domain models
│   │   ├── media.go         # Media entity
│   │   ├── search.go        # Search models
│   │   ├── upload.go        # Upload models
│   │   ├── errors.go        # Business errors
│   │   └── constants.go     # Domain constants
│   │
│   ├── service/              # Business logic layer
│   │   ├── media_service.go     # Media business logic
│   │   ├── search_service.go    # Search business logic
│   │   └── media_event_handler.go # Event handling
│   │
│   ├── repository/           # Data access layer
│   │   ├── media_repository.go      # Media data access
│   │   ├── search_repository.go     # Search data access (PostgreSQL)
│   │   └── elasticsearch_search_repository.go # Elasticsearch implementation
│   │
│   ├── handler/              # HTTP request handlers
│   │   ├── media_handler.go     # Media API endpoints
│   │   └── search_handler.go    # Search API endpoints
│   │
│   └── middleware/           # HTTP middleware
│       └── middleware.go     # CORS, logging, recovery
│
├── pkg/                      # Public, reusable packages
│   ├── database/            # Database connections
│   ├── elasticsearch/       # Elasticsearch client
│   ├── httpclient/         # HTTP client utilities
│   └── messagequeue/       # Message queue interface
│
├── migrations/              # Database migration files
├── docker-compose.yml      # Container orchestration
├── go.mod                 # Go module dependencies
└── README.md             # This file
```

### Layer Responsibilities

#### 🎯 Domain Layer (`internal/domain/`)
- **Purpose**: Core business entities and rules
- **Contains**: Models, business errors, constants
- **Dependencies**: None (innermost layer)

#### 🏢 Service Layer (`internal/service/`)
- **Purpose**: Business logic and use cases
- **Contains**: Application services, business rules
- **Dependencies**: Domain layer only

#### 💾 Repository Layer (`internal/repository/`)
- **Purpose**: Data access and persistence
- **Contains**: Database operations, external service calls
- **Dependencies**: Domain layer

#### 🌐 Handler Layer (`internal/handler/`)
- **Purpose**: HTTP request/response handling
- **Contains**: API endpoints, input validation
- **Dependencies**: Service layer

## 📋 Prerequisites

Before running the project, ensure you have:

### Required Software
- **Go**: Version 1.21 or higher
  ```bash
  go version  # Should show go1.21+
  ```

- **Docker & Docker Compose**: For running infrastructure
  ```bash
  docker --version && docker-compose --version
  ```

- **PostgreSQL**: Version 15+ (or use Docker)
- **Elasticsearch**: Version 8.11+ (or use Docker)

### System Requirements
- **RAM**: Minimum 4GB (8GB recommended)
- **Disk Space**: 2GB for dependencies and data
- **OS**: macOS, Linux, or Windows with WSL2

## 🚀 Installation & Setup

### 1. Clone the Repository
```bash
git clone https://github.com/a-samir97/th-backend.git
cd thamaniyah
```

### 2. Install Go Dependencies
```bash
go mod download
go mod tidy
```

### 3. Start Infrastructure Services
```bash
# Start all infrastructure (PostgreSQL + Elasticsearch + Redis)
docker-compose up -d

# Or start individually
docker-compose up -d postgresql
docker-compose up -d elasticsearch
docker-compose up -d redis
```

### 4. Wait for Services to be Ready
```bash
# Check PostgreSQL
docker-compose exec postgresql pg_isready -U postgres

# Check Elasticsearch
curl -X GET "localhost:9200/_cluster/health?pretty"

# Should return "status": "green" or "yellow"
```

### 5. Run Database Migrations
```bash
# Create tables and indexes
go run cmd/migrate/main.go
```

### 6. Start Application Services

#### Option A: Run in Development Mode
```bash
# Terminal 1: Start CMS Service
go run cmd/cms-service/main.go

# Terminal 2: Start Discovery Service  
go run cmd/discovery-service/main.go
```

#### Option B: Run in Background
```bash
# Start both services in background
nohup go run cmd/cms-service/main.go > cms.log 2>&1 &
nohup go run cmd/discovery-service/main.go > discovery.log 2>&1 &
```

### 7. Verify Installation
```bash
# Check service health
curl -X GET http://localhost:8080/health  # CMS Service
curl -X GET http://localhost:8081/health  # Discovery Service

# Expected response:
# {"status":"ok","service":"cms-service","timestamp":"2025-08-27T..."}
```

### 8. Populate Search Index
```bash
# Index existing media into Elasticsearch
curl -X POST http://localhost:8081/api/v1/search/reindex
```

## 📚 API Documentation

### 🎥 CMS Service (Port 8080)

#### Media Upload Flow

**Step 1: Generate Upload URL**
```bash
POST /api/v1/media/upload-url
Content-Type: application/json

{
  "title": "My Video Tutorial",
  "description": "Learn Go programming basics",
  "filename": "go-tutorial.mp4",
  "file_size": 52428800,
  "type": "video"
}
```

**Response:**
```json
{
  "media_id": "550e8400-e29b-41d4-a716-446655440000",
  "upload_url": "http://localhost:8080/upload/uploads/550e8400-e29b-41d4-a716-446655440000.mp4",
  "expires_at": "2025-08-27T14:30:00Z"
}
```

**Step 2: Confirm Upload**
```bash
POST /api/v1/media/{media_id}/confirm
```

#### Media Management

**Get All Media (with pagination)**
```bash
GET /api/v1/media?limit=20&offset=0
```

**Get Single Media**
```bash
GET /api/v1/media/{media_id}
```

**Update Media Metadata**
```bash
PUT /api/v1/media/{media_id}
Content-Type: application/json

{
  "title": "Updated Video Title",
  "description": "Updated description"
}
```

**Delete Media**
```bash
DELETE /api/v1/media/{media_id}
```

### 🔍 Discovery Service (Port 8081)

#### Advanced Search

**Full-Text Search**
```bash
GET /api/v1/search?query=golang tutorial&limit=10&offset=0

# With type filter
GET /api/v1/search?query=machine learning&type=video

# Response includes relevance scores
{
  "results": [
    {
      "media": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "Advanced Golang Tutorial",
        "description": "Deep dive into Go programming...",
        "type": "video",
        "status": "ready"
      },
      "score": 2.3456789  # Elasticsearch relevance score
    }
  ],
  "total": 15,
  "query": "golang tutorial",
  "limit": 10,
  "offset": 0
}
```

**Search Suggestions**
```bash
GET /api/v1/search/suggest?query=gola&limit=5

{
  "suggestions": [
    {"text": "Golang Tutorial", "count": 3},
    {"text": "Golang Best Practices", "count": 1}
  ],
  "query": "gola"
}
```

**Rebuild Search Index**
```bash
POST /api/v1/search/reindex
# Use this when you need to refresh Elasticsearch with latest data
```

## 💾 Database Schema

### PostgreSQL Schema

#### `media_files` Table
```sql
CREATE TABLE media_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    duration INTEGER DEFAULT 0,        -- in seconds
    format VARCHAR(50),                -- mp4, mp3, avi, etc.
    type VARCHAR(20) NOT NULL,         -- video, podcast
    status VARCHAR(20) DEFAULT 'uploading', -- uploading, ready, failed
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL          -- soft delete
);

-- Indexes for performance
CREATE INDEX idx_media_created_at ON media_files(created_at DESC);
CREATE INDEX idx_media_type ON media_files(type);
CREATE INDEX idx_media_status ON media_files(status);
CREATE INDEX idx_media_files_deleted_at ON media_files(deleted_at);
```

#### `search_index` Table (Backup/Sync)
```sql
CREATE TABLE search_index (
    id UUID PRIMARY KEY,
    media_id UUID REFERENCES media_files(id),
    title VARCHAR(255),
    description TEXT,
    content TEXT,                      -- Combined searchable text
    type VARCHAR(20),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Full-text search index
CREATE INDEX idx_search_content ON search_index 
USING GIN(to_tsvector('english', content));
```

### Elasticsearch Mapping

```json
{
  "mappings": {
    "properties": {
      "id": {"type": "keyword"},
      "title": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "description": {"type": "text", "analyzer": "standard"},
      "content": {"type": "text", "analyzer": "standard"},
      "type": {"type": "keyword"},
      "status": {"type": "keyword"},
      "file_path": {"type": "keyword"},
      "file_size": {"type": "long"},
      "duration": {"type": "integer"},
      "format": {"type": "keyword"},
      "created_at": {"type": "date"},
      "updated_at": {"type": "date"}
    }
  },
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "standard": {"type": "standard"}
      }
    }
  }
}
```

## 🔧 Usage Examples

### Complete Media Upload Workflow

```bash
# 1. Generate upload URL
UPLOAD_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/media/upload-url 
  -H "Content-Type: application/json" 
  -d '{
    "title": "Golang Microservices Tutorial",
    "description": "Building scalable microservices with Go and Docker",
    "filename": "golang-microservices.mp4",
    "file_size": 104857600,
    "type": "video"
  }')

# Extract media ID
MEDIA_ID=$(echo $UPLOAD_RESPONSE | jq -r '.media_id')
echo "Created media: $MEDIA_ID"

# 2. Simulate file upload (in real app, frontend uploads to presigned URL)
echo "File would be uploaded to: $(echo $UPLOAD_RESPONSE | jq -r '.upload_url')"

# 3. Confirm upload completion
curl -X POST "http://localhost:8080/api/v1/media/$MEDIA_ID/confirm"

# 4. Verify media is ready
curl -X GET "http://localhost:8080/api/v1/media/$MEDIA_ID" | jq .

# 5. Index in search engine
curl -X POST http://localhost:8081/api/v1/search/reindex

# 6. Search for the content
curl -X GET "http://localhost:8081/api/v1/search?query=golang microservices" | jq .
```

### Advanced Search Examples

```bash
# Fuzzy search (handles typos)
curl -X GET "http://localhost:8081/api/v1/search?query=golong tutrial"  # finds "golang tutorial"

# Type-specific search
curl -X GET "http://localhost:8081/api/v1/search?query=machine learning&type=video"

# Pagination
curl -X GET "http://localhost:8081/api/v1/search?query=programming&limit=5&offset=10"

# Get suggestions for autocomplete
curl -X GET "http://localhost:8081/api/v1/search/suggest?query=prog&limit=3"
```

## 👩‍💻 Development Guide
### Code Standards

#### Naming Conventions
- **Packages**: lowercase, single words (`service`, `handler`)
- **Files**: snake_case (`media_service.go`)
- **Functions**: PascalCase for public, camelCase for private
- **Constants**: ALL_CAPS with underscores

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestMediaService_CreateUploadURL ./internal/service

# Run tests with race detection
go test -race ./...
```

## 🚢 Deployment

### Local Deployment

```bash
# Build binaries
go build -o bin/cms-service cmd/cms-service/main.go
go build -o bin/discovery-service cmd/discovery-service/main.go

# Run with Docker Compose
docker-compose up -d --build
```
#### Health Check Endpoints
```bash
# Service health
GET /health
GET /metrics      # Prometheus metrics (future)
GET /debug/pprof  # Go profiling (dev only)
