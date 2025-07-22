## Fundamentals of Go Tasks: Task 7 - Refactoring Task Management API using Clean Architecture Principles

## Natnael Eyuel - A2SV G62

## Overview
A refactored task management system following Clean Architecture principles to improve maintainability, testability, and scalability.

## Key Features
- JWT Authentication
- Role-based access control (Admin/User)
- CRUD operations for tasks
- MongoDB persistence

## Architecture
```
task-manager/
├── Delivery/        # HTTP handlers (Gin)
├── Domain/          # Entities and interfaces
├── Infrastructure/  # JWT, Hashing, Config
├── Repositories/    # MongoDB implementations
└── Usecases/        # Business logic
```

## Setup

### Prerequisites
- Go installed
- MongoDB installed

### Installation
```bash
go get .
```

### Configuration
1. Create `.env` file:
```env
JWT_SECRET=your_secret_key
```

### Running
```bash
go run Delivery/main.go
```

## API Documentation
See [API_DOCS.md](/docs/api_documentation.md) for endpoint specifications.

## Key Design Decisions
1. **Dependency Rule**:  
   Outer layers (Delivery/Infra) depend inward on interfaces defined in Domain.

2. **Validation**:  
   Business rules enforced in Use Case layer.



 

