# Go Hexagonal Architecture Framework

A production-ready Go microservice framework following **Hexagonal Architecture** (also known as Ports & Adapters) with **Clean Architecture** principles. This framework provides a scalable, testable, and maintainable structure for building Go applications.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Core Concepts](#core-concepts)
  - [Edges (Inbound Adapters)](#edges-inbound-adapters)
  - [Chain (Application Layer)](#chain-application-layer)
  - [Service (Business Logic)](#service-business-logic)
  - [Data (Domain Layer)](#data-domain-layer)
  - [Infrastructure (Outbound Adapters)](#infrastructure-outbound-adapters)
- [AhuM CLI Tool](#ahum-cli-tool)
- [Configuration](#configuration)
- [Security](#security)
- [Getting Started](#getting-started)
- [Creating a New Feature](#creating-a-new-feature)
- [Testing](#testing)
- [Deployment](#deployment)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              EDGES (Inbound)                                │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Connect (gRPC/HTTP)  │  REST API  │  GraphQL  │  CLI  │  Events   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────┬────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          CHAIN (Application Layer)                          │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Orchestration  │  Validation  │  Error Mapping  │  DI Container   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────┬────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         SERVICE (Business Logic)                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Use Cases  │  Business Rules  │  Domain Operations  │  JWT/Crypto │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────┬────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           DATA (Domain Layer)                               │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Entities  │  Repository Interfaces  │  Domain Errors  │  Contracts│   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────┬────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      INFRASTRUCTURE (Outbound Adapters)                     │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  PostgreSQL  │  Redis  │  SMS Gateway  │  Email  │  External APIs  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Dependency Flow

```
Edge → Chain → Service → Data ← Infrastructure
                           │
                           └── Interfaces defined here, implemented in Infrastructure
```

**Key Principle**: Dependencies point inward. The domain layer (`data`) has no knowledge of infrastructure or delivery mechanisms.

---

## Project Structure

```
.
├── main.go                     # Application entry point
├── config/
│   ├── config.go               # Configuration loading & infrastructure setup
│   └── config.json             # Application configuration
├── edge/
│   ├── edge.go                 # Edge interface & edge starter
│   └── connect/                # Connect (gRPC-Web) edge implementation
│       ├── connect.go          # HTTP/2 server setup
│       ├── gen/                # Generated protobuf code
│       └── {domain}/           # Domain-specific edge handlers
│           ├── registrar.go    # Service registration
│           └── v1/             # API versioning
│               ├── edge.go     # Edge handler implementation
│               └── *.proto     # Protocol buffer definitions
├── chain/
│   └── {domain}/
│       ├── chain.go            # Dependency injection & service composition
│       ├── entity.go           # Chain-level DTOs (request/response)
│       └── handler.go          # Application orchestration layer
├── service/
│   └── {domain}/
│       ├── service.go          # Service constructor with DI
│       ├── entity.go           # Service-level DTOs
│       ├── handler.go          # Business logic implementation
│       └── error.go            # Domain-specific errors
├── data/
│   └── {domain}/
│       ├── entity.go           # Domain entities (core models)
│       ├── repo.go             # Repository interface definition
│       └── error.go            # Domain errors
├── infrastructure/
│   ├── postgres/
│   │   ├── config.go           # Database configuration
│   │   ├── connection.go       # Connection management
│   │   ├── migrations/         # SQL migration files
│   │   └── {domain}/
│   │       ├── model.go        # Database models & mappers
│   │       └── repo.go         # Repository implementation
│   ├── sms/                    # SMS provider implementation
│   └── smstest/                # Mock SMS for testing
├── crypto/
│   ├── config.go               # Encryption configuration
│   ├── encrypter.go            # Encrypter interface
│   └── mock/                   # Mock implementations for testing
├── security/
│   └── string.go               # SecureString type with auto-encryption
├── jwthelper/
│   ├── config.go               # JWT configuration
│   └── helper.go               # JWT token utilities
├── log/
│   └── log.go                  # Structured logging (slog)
└── Dockerfile                  # Multi-stage Docker build
```

---

## Core Concepts

### Edges (Inbound Adapters)

**Location**: `edge/`

Edges are the entry points to your application. They handle protocol-specific concerns and translate external requests into domain operations.

```go
// edge/edge.go
type Edge interface {
    Configure()
    Start(context.Context, *sync.WaitGroup)
}
```

**Supported Edge Types**:
- **Connect (gRPC-Web)**: HTTP/2 with Connect protocol (default)
- REST API, GraphQL, CLI, Event consumers (extensible)

**Edge Structure**:
```go
// edge/connect/{domain}/v1/edge.go
type Edge struct {
    chain domain.Chain  // Injected chain
}

func NewEdge() *Edge {
    chain := domain.NewChain()
    return &Edge{chain: *chain}
}

func (e *Edge) Handler(ctx context.Context, req *Request) (*Response, error) {
    result, err := e.chain.Operation(ctx, req.Param)
    // Map errors, build response
    return response, nil
}
```

**API Versioning**: Each domain edge supports multiple versions (`v1/`, `v2/`, etc.) for backward compatibility.

---

### Chain (Application Layer)

**Location**: `chain/`

The chain layer acts as the orchestration layer between edges and services. It handles:
- **Dependency Injection**: Composing services with their dependencies
- **Validation**: Input validation before reaching business logic
- **Error Mapping**: Converting service errors to edge-appropriate errors
- **Cross-cutting Concerns**: Logging, metrics, tracing

```go
// chain/{domain}/chain.go
type Chain struct {
    svc *service.Service
}

func NewChain() *Chain {
    // Infrastructure wiring
    repo := postgres.NewRepo(postgres.Db)
    externalSvc := external.New()
    
    // Service composition
    svc := service.NewService(repo, externalSvc)
    return &Chain{svc}
}
```

```go
// chain/{domain}/handler.go
func (c *Chain) Operation(ctx context.Context, input string) (*Response, error) {
    // Validate input
    if !validator.Validate(input) {
        return nil, ErrInvalidInput
    }
    
    // Delegate to service
    result, err := c.svc.DoOperation(ctx, input)
    if err != nil {
        log.Logger.Error("chain: operation failed", "error", err)
        return nil, mapError(err)
    }
    
    return &Response{Data: result.Data}, nil
}
```

---

### Service (Business Logic)

**Location**: `service/`

Services contain pure business logic. They are protocol-agnostic and depend only on repository interfaces defined in the `data` layer.

```go
// service/{domain}/service.go
type Service struct {
    repo    data.Repo          // Repository interface
    sms     sms.Service        // External service interface
}

func NewService(repo data.Repo, sms sms.Service) *Service {
    return &Service{repo: repo, sms: sms}
}
```

```go
// service/{domain}/handler.go
func (svc *Service) DoOperation(ctx context.Context, input string) (*Result, error) {
    // Business logic
    entity := data.NewEntity(input)
    
    // Persist through repository interface
    saved, err := svc.repo.Create(entity)
    if err != nil {
        return nil, ErrPersistFailed
    }
    
    // External service call
    err = svc.sms.Send(ctx, entity.Phone, "message")
    
    return &Result{ID: saved.ID}, nil
}
```

---

### Data (Domain Layer)

**Location**: `data/`

The data layer defines the core domain entities and repository interfaces. This is the innermost layer and has no external dependencies.

```go
// data/{domain}/entity.go
type Entity struct {
    ID        *string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
    Status    Status
}

func New(name string) *Entity {
    return &Entity{
        Name:   name,
        Status: Pending,
    }
}
```

```go
// data/{domain}/repo.go
type Repo interface {
    Create(Entity) (Entity, error)
    Find(Entity) (Entity, error)
    Update(Entity) (Entity, error)
    Delete(id string) error
}
```

**Key Principle**: Repository interfaces are defined here but implemented in the `infrastructure` layer.

---

### Infrastructure (Outbound Adapters)

**Location**: `infrastructure/`

Infrastructure implements the interfaces defined in the `data` layer. This includes database repositories, external service clients, and third-party integrations.

```go
// infrastructure/postgres/{domain}/repo.go
type Repo struct {
    db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
    return &Repo{db: db}
}

func (r *Repo) Create(e data.Entity) (data.Entity, error) {
    model := toModel(e)
    model.ID = cuid2.Generate()
    err := r.db.Create(&model).Error
    if err != nil {
        return data.Entity{}, err
    }
    return toEntity(model), nil
}
```

```go
// infrastructure/postgres/{domain}/model.go
type Model struct {
    ID        string    `gorm:"primaryKey,type:varchar"`
    Name      string    `gorm:"type:varchar(255)"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func toModel(e data.Entity) Model { /* ... */ }
func toEntity(m Model) data.Entity { /* ... */ }
```

**Test Implementations**: Each infrastructure has a corresponding test implementation (e.g., `smstest/`) for unit testing without external dependencies.

---

## AhuM CLI Tool

**Repository**: [https://github.com/Ahu-Tools/ahum](https://github.com/Ahu-Tools/ahum)

AhuM is a powerful command-line interface (CLI) tool designed to streamline the development and management of projects following this hexagonal architecture. It automates project scaffolding, code generation, and maintains consistency across the codebase.

### Installation

```bash
go install github.com/Ahu-Tools/ahum@latest
```

### Core Commands

| Command | Description |
|---------|-------------|
| `ahum init` | Initialize a new project with interactive TUI |
| `ahum service create` | Create a new service (chain, service, data layers) |
| `ahum edge generate` | Generate a new edge (inbound adapter) |
| `ahum infra generate` | Add new infrastructure components |
| `ahum connect service add` | Add a new Connect/gRPC service |
| `ahum connect service version add` | Add a new version to a Connect service |
| `ahum connect service method add` | Add a new method to a Connect service |
| `ahum connect gen` | Generate protobuf code |

### Project Initialization

Create a new project with the interactive wizard:

```bash
ahum init
```

This launches a TUI (Terminal User Interface) that guides you through:
- Project name and module path
- Infrastructure selection (PostgreSQL, Redis, etc.)
- Edge selection (Connect, Gin, Asynq)
- Initial configuration

### Managing Services

Create a new service with all three layers (chain, service, data):

```bash
ahum service create
```

This generates:
```
chain/{service}/
  ├── chain.go      # Dependency injection
  ├── entity.go     # Chain DTOs
  └── handler.go    # Application orchestration

service/{service}/
  ├── service.go    # Service constructor
  ├── entity.go     # Service DTOs
  └── handler.go    # Business logic

data/{service}/
  ├── entity.go     # Domain entities
  └── repo.go       # Repository interface
```

### Managing Connect (gRPC) Services

#### Add a New Service

```bash
ahum connect service add [service_name]
```

Example:
```bash
ahum connect service add product
```

#### Add a New Version

```bash
ahum connect service version add [version_name] [service_name]
```

Example:
```bash
ahum connect service version add v2 product
```

#### Add a New Method

```bash
ahum connect service method add [method_name] [service_name] [version_name]
```

Example:
```bash
ahum connect service method add GetProduct product v1
```

This automatically:
1. Adds the method definition to the `.proto` file
2. Runs `buf generate` to generate Go code
3. Adds the method stub to the edge handler

#### Regenerate Protobuf Code

```bash
ahum connect gen
```

### Managing Gin Routes

#### Add a Version

```bash
ahum gin route version add [version]
```

#### Add an Entity

```bash
ahum gin route entity add [version_name] [entity_name]
```

#### Add a Handler

```bash
ahum gin route handle add [version_name] [entity_name] [method_name]
```

### Managing Asynq (Background Tasks)

#### Add a Module

```bash
ahum asynq edge add module [version] [module_name]
```

#### Add a Task Handler

```bash
ahum asynq edge add task [version] [module_name] [task_name]
```

### Infrastructure Management

Generate new infrastructure components:

```bash
ahum infra generate
```

This launches an interactive form to configure:
- PostgreSQL connection
- Redis connection
- Other supported infrastructures

### Edge Management

Generate a new edge:

```bash
ahum edge generate
```

Supported edge types:
- **Connect** (gRPC-Web compatible)
- **Gin** (REST API)
- **Asynq** (Background task processing)

### Code Markers

AhuM uses special comment markers (e.g., `// @ahum: imports`, `// @ahum: services`) to identify injection points in generated code. These markers allow AhuM to safely add new code without overwriting existing customizations.

Example markers in the codebase:
```go
// edge/connect/connect.go
func RegisterServices(mux *http.ServeMux) {
    hello.RegisterService(mux)
    otp.RegisterService(mux)
    // @ahum: services    <-- New services are injected here
}
```

```go
// config/config.go
func ConfigInfras() error {
    // @ahum:infras.group
    err := postgres.Configure()
    // @ahum:end.infras.group
    
    //@ahum: loads        <-- New infrastructure loads are injected here
    return nil
}
```

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--path` | `-p` | Project root path (default: current directory) |
| `--config` | | Custom config file path |

Example:
```bash
ahum -p /path/to/project connect service add user
```

### Workflow Example

Complete workflow for adding a new "product" feature:

```bash
# 1. Create the service (chain, service, data layers)
ahum service create
# Follow the TUI prompts to create "product" service

# 2. Add a Connect service for the product
ahum connect service add product

# 3. Add version v1
ahum connect service version add v1 product

# 4. Add methods
ahum connect service method add CreateProduct product v1
ahum connect service method add GetProduct product v1
ahum connect service method add ListProducts product v1

# 5. Implement the infrastructure (PostgreSQL repo)
ahum infra generate
# Select PostgreSQL and configure for "product"

# 6. Wire everything in chain.go and implement business logic
```

---

## Configuration

**Location**: `config/`

Configuration is managed via Viper with support for JSON files and environment variables.

```json
// config/config.json
{
  "app": {
    "secret_key": "your-secret-key",
    "env": "dev"
  },
  "jwt": {
    "private_key_file": "config/private.pem",
    "public_key_file": "config/public.pem",
    "signing_algorithm": "RS256"
  },
  "edges": {
    "connect": {
      "server": {
        "host": "0.0.0.0",
        "port": "8080"
      }
    }
  },
  "infras": {
    "postgres": {
      "user": "postgres",
      "password": "postgres",
      "db_name": "mydb",
      "host": "localhost",
      "port": "5432",
      "sslmode": "disable"
    }
  }
}
```

**Configuration Flow**:
1. `CheckConfigs()`: Validates required configuration
2. `ConfigInfras()`: Initializes infrastructure connections

---

## Security

### SecureString Type

**Location**: `security/string.go`

`SecureString` provides automatic encryption/decryption when storing sensitive data in the database.

```go
type User struct {
    ID    string
    Phone security.SecureString  // Auto-encrypted in DB
    Email security.SecureString  // Auto-encrypted in DB
}
```

Features:
- **Automatic encryption** via GORM's `Value()` hook
- **Automatic decryption** via GORM's `Scan()` hook
- **Redacted JSON output** to prevent accidental logging

### Encryption Interface

```go
// crypto/encrypter.go
type Encrypter interface {
    Encrypt(plaintext []byte) ([]byte, error)
    Decrypt(ciphertext []byte) ([]byte, error)
    ComputeBlindIndex(plaintext string) string  // For searchable encryption
}
```

### JWT Authentication

**Location**: `jwthelper/`

RSA-based asymmetric JWT tokens for secure authentication.

```go
// Generate token
claims := jwt.MapClaims{"user_id": "123", "exp": time.Now().Add(24*time.Hour).Unix()}
token, err := jwthelper.GenerateToken(claims)

// Verify token
parsedToken, err := jwthelper.ParseToken(tokenString)
```

---

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 14+
- Protocol Buffers compiler (for Connect/gRPC)
- AhuM CLI (recommended)

### Installation

#### Option 1: Using AhuM CLI (Recommended)

Install AhuM and create a new project:

```bash
# Install AhuM CLI
go install github.com/Ahu-Tools/ahum@latest

# Initialize a new project
ahum init
```

The interactive TUI will guide you through:
- Project name and Go module path
- Infrastructure selection (PostgreSQL, Redis, etc.)
- Edge selection (Connect, Gin, Asynq)
- Initial configuration

#### Option 2: Clone Existing Project

```bash
# Clone the repository
git clone <repository-url>
cd example

# Install dependencies
go mod download

# Generate protobuf code (if using Connect)
ahum connect gen
# Or manually: buf generate

# Setup database
# Run migrations from infrastructure/postgres/migrations/

# Start the application
go run main.go
```

### Running with Docker

```bash
# Build the image
docker build -t myapp .

# Run the container
docker run -p 8080:8080 myapp
```

---

## Creating a New Feature

Follow these steps to add a new domain feature (e.g., "product") using the **AhuM CLI**:

### 1. Create Service Layers (Chain, Service, Data)

Use AhuM to generate all three layers at once:

```bash
ahum service create
```

Follow the interactive TUI prompts:
- Enter service name: `product`
- Configure options as needed

This automatically generates:
```
chain/product/
  ├── chain.go      # Dependency injection
  ├── entity.go     # Chain DTOs
  └── handler.go    # Application orchestration

service/product/
  ├── service.go    # Service constructor
  ├── entity.go     # Service DTOs
  └── handler.go    # Business logic

data/product/
  ├── entity.go     # Domain entities
  └── repo.go       # Repository interface
```

### 2. Add Connect (gRPC) Service

```bash
ahum connect service add product
```

### 3. Add API Version

```bash
ahum connect service version add v1 product
```

### 4. Add Methods

```bash
ahum connect service method add CreateProduct product v1
ahum connect service method add GetProduct product v1
ahum connect service method add ListProducts product v1
```

Each method command automatically:
- Adds the RPC definition to the `.proto` file
- Runs `buf generate` to generate Go code
- Adds the method stub to the edge handler

### 5. Implement Infrastructure (Manual Step)

> **Note**: Infrastructure repository generation is planned for a future AhuM release (see [Roadmap](#ahum-roadmap)).

For now, manually create the PostgreSQL repository:

```bash
mkdir -p infrastructure/postgres/product
```

```go
// infrastructure/postgres/product/model.go
package product

type Product struct {
    ID        string    `gorm:"primaryKey,type:varchar"`
    Name      string    `gorm:"type:varchar(255)"`
    Price     float64   `gorm:"type:decimal(10,2)"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func toModel(e data.Product) Product { /* ... */ }
func toEntity(m Product) data.Product { /* ... */ }
```

```go
// infrastructure/postgres/product/repo.go
package product

func NewRepo(db *gorm.DB) *Repo {
    return &Repo{db: db}
}

func (r *Repo) Create(e data.Product) (data.Product, error) {
    model := toModel(e)
    model.ID = cuid2.Generate()
    err := r.db.Create(&model).Error
    if err != nil {
        return data.Product{}, err
    }
    return toEntity(model), nil
}
```

### 6. Wire Dependencies in Chain

Update the generated `chain/product/chain.go` to inject the repository:

```go
// chain/product/chain.go
package product

import (
    "gitlab.com/yourproject/infrastructure/postgres"
    postgresProduct "gitlab.com/yourproject/infrastructure/postgres/product"
    productSvc "gitlab.com/yourproject/service/product"
)

func NewChain() *Chain {
    repo := postgresProduct.NewRepo(postgres.Db)
    svc := productSvc.NewService(repo)
    return &Chain{svc: svc}
}
```

### 7. Implement Business Logic

Add your business logic in the generated handler files:
- `chain/product/handler.go` - Validation and orchestration
- `service/product/handler.go` - Core business rules

### Quick Reference

| Step | Command |
|------|---------|
| Create service layers | `ahum service create` |
| Add Connect service | `ahum connect service add [name]` |
| Add version | `ahum connect service version add [version] [service]` |
| Add method | `ahum connect service method add [method] [service] [version]` |
| Regenerate protos | `ahum connect gen` |

---

## Testing

### Unit Testing

Each layer can be tested independently using mock implementations.

```go
// Use mock repository
mockRepo := &mock.ProductRepo{}
svc := service.NewService(mockRepo)

// Test service logic
result, err := svc.CreateProduct(ctx, input)
assert.NoError(t, err)
assert.NotNil(t, result)
```

### Integration Testing

Use test implementations from `infrastructure/*test/` packages:

```go
// Use test SMS service
smsSvc := smstest.New()
chain := otp.NewChainWithSMS(repo, smsSvc)
```

---

## Deployment

### Docker

```bash
docker build -t myapp .
docker run -e APP_ENV=prod -p 8080:8080 myapp
```

### Environment Variables

Override configuration via environment variables:
- `APP_SECRET_KEY`
- `INFRAS_POSTGRES_HOST`
- `INFRAS_POSTGRES_PASSWORD`

---

## Key Technologies

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| API Protocol | Connect (gRPC-Web compatible) |
| Database | PostgreSQL with GORM |
| Configuration | Viper |
| Logging | slog (structured logging) |
| Authentication | JWT (RSA256) |
| Encryption | AES-GCM with blind indexing |
| ID Generation | CUID2 |
| Protobuf | Protocol Buffers v3 |

---

## Design Principles

1. **Dependency Inversion**: High-level modules don't depend on low-level modules
2. **Single Responsibility**: Each layer has a specific purpose
3. **Interface Segregation**: Small, focused interfaces
4. **Open/Closed**: Extend behavior without modifying existing code
5. **Testability**: Every layer is independently testable

---

## License

MIT

---
