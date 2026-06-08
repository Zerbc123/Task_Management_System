# Task Management System

A Task Management REST API built with Go, Gin, PostgreSQL, Redis, JWT Authentication, WebSocket, Background Worker, and Docker.

---

## Features

* JWT Authentication (Register/Login)
* Project Management
* Task Management
* Task Comment System
* Redis Cache
* Background Notification Worker
* WebSocket Realtime Updates
* Health Check Endpoint
* Docker & Docker Compose
* Unit Tests & Integration Tests

---

## Tech Stack

* Go
* Gin
* PostgreSQL
* Redis
* GORM
* JWT
* Gorilla WebSocket
* Docker
* Docker Compose

---

## Project Structure

```text
cmd/
├── server/
│   └── app.go
└── worker/
    └── main.go

internal/
├── auth/
├── cache/
├── comment/
├── database/
├── health/
├── middleware/
├── notification/
├── project/
├── task/
├── user/
└── websocket/
```

---

## Running with Docker

Build and start all services:

```bash
docker compose up --build -d
```

Check containers:

```bash
docker ps
```

Stop all services:

```bash
docker compose down
```

---

## Running Locally

### Start PostgreSQL and Redis

```bash
docker compose up postgres redis -d
```

### Run API Server

```bash
go run ./cmd/server
```

### Run Background Worker

```bash
go run ./cmd/worker
```

---

## Environment Variables

Create a `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=task_management
DB_SSLMODE=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=your-secret-key
```

---

## Health Check

### Request

```http
GET /health
```

### Response

```json
{
  "status": "UP",
  "database": "UP",
  "redis": "UP"
}
```

---

## Authentication APIs

### Register

```http
POST /auth/register
```

Body:

```json
{
  "email": "user@example.com",
  "password": "123456",
  "full_name": "John Doe"
}
```

---

### Login

```http
POST /auth/login
```

Body:

```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

Use returned token:

```http
Authorization: Bearer <token>
```

---

## Project APIs

| Method | Endpoint       |
| ------ | -------------- |
| POST   | /projects      |
| GET    | /projects      |
| GET    | /projects/{id} |
| PUT    | /projects/{id} |
| DELETE | /projects/{id} |

---

## Task APIs

| Method | Endpoint    |
| ------ | ----------- |
| POST   | /tasks      |
| GET    | /tasks      |
| GET    | /tasks/{id} |
| PUT    | /tasks/{id} |
| DELETE | /tasks/{id} |

---

## Comment APIs

| Method | Endpoint             |
| ------ | -------------------- |
| POST   | /tasks/{id}/comments |
| GET    | /tasks/{id}/comments |
| PUT    | /comments/{id}       |
| DELETE | /comments/{id}       |

---

## WebSocket

Connect:

```text
ws://localhost:8080/ws
```

Realtime Events:

### Task Updated

```json
{
  "type": "task.updated",
  "task_id": "uuid",
  "title": "Update Docker",
  "status": "DONE"
}
```

### Comment Created

```json
{
  "type": "comment.created",
  "task_id": "uuid",
  "content": "Redis queue completed"
}
```

---

## Testing

Run all tests:

```bash
go test ./...
```

Run integration tests:

```bash
go test ./internal/integration -v
```

Run race detector:

```bash
go test -race ./...
```

Run vet:

```bash
go vet ./...
```

---

## Docker Services

| Service    | Port |
| ---------- | ---- |
| API        | 8080 |
| PostgreSQL | 5432 |
| Redis      | 6379 |

---

## CI Pipeline

GitHub Actions automatically runs:

* go test ./...
* go vet ./...
* golangci-lint

for every push and pull request.

---

## Author

**Đào Lê Hoàng Duy**

Task Management System Project - Go Backend Learning Journey
