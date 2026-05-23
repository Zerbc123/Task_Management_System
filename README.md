# Task Management System

REST API đơn giản để quản lý Task và Category, được xây dựng bằng Golang.

## Tech Stack

- Golang
- Gin Framework
- GORM
- PostgreSQL
- Docker Compose
- golang-migrate

## Features

### Task
- Create task
- Get all tasks
- Get task by ID
- Update task
- Delete task

### Category
- Create category
- Get all categories
- Get category by ID
- Update category
- Delete category

## Project Structure
cmd/
└── server/
    └── app.go

internal/
├── database/
├── middleware/
├── shared/
│   └── response/
├── task/
│   ├── dto/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   └── services/
└── category/
    ├── dto/
    ├── handler/
    ├── model/
    ├── repository/
    └── services/

migrations/
docker-compose.yml
.env
go.mod
README.md

## Architecture Flow

Client
  ↓
Middleware
  ↓
Handler
  ↓
Service
  ↓
Repository
  ↓
PostgreSQL

## Environment Variables

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=task_management
DB_SSLMODE=disable

## Run Database

docker compose up -d

- Kiểm tra container:

docker ps

## Run Migration

migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/task_management?sslmode=disable" up

## Rollback migration

migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/task_management?sslmode=disable" down

## Run Application

go run cmd/server/app.go

- Server chạy tại:

http://localhost:8080

## API Endpoints

### Task API
Method	        Endpoint	        Description
POST	        /tasks	            Create task
GET	            /tasks	            Get all tasks
GET	            /tasks/:id	        Get task by ID
PUT	            /tasks/:id	        Update task
DELETE	        /tasks/:id	        Delete task

### Category API
Method	        Endpoint	        Description
POST	        /categories	        Create category
GET	            /categories	        Get all categories
GET	            /categories/:id	    Get category by ID
PUT	            /categories/:id	    Update category
DELETE	        /categories/:id	    Delete category

## Example Create Task

- Request:

{
  "title": "Learn Go",
  "description": "Build CRUD API",
  "status": "TODO",
  "priority": 3,
  "assigned_to": "Duy"
}

- Response:

{
  "success": true,
  "message": "task created successfully",
  "data": {
    "id": "uuid",
    "title": "Learn Go",
    "description": "Build CRUD API",
    "status": "TODO",
    "priority": 3,
    "assigned_to": "Duy"
  }
}

## Response Format

- Success:

{
  "success": true,
  "message": "success message",
  "data": {}
}

- Error:

{
  "success": false,
  "message": "error message",
  "error": "error detail"
}

## Middleware

### Project có các middleware:

- Request ID
- Request Logger
- Recovery

## Repository Pattern

Project sử dụng Repository Pattern để tách business logic khỏi database logic.

Service
  ↓
Repository Interface
  ↓
GORM Repository
  ↓
Database

## Author
Duy
