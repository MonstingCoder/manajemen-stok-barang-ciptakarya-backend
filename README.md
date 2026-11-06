# Cipta Karya Inventory Management Backend

A Go-based backend for managing stock and inventory data for Cipta Karya. The service includes user authentication, role-based authorization, inventory items, and tag-based organization for item listings.

## Overview

This project provides a REST API built with Go, PostgreSQL, and GORM. It is designed for managing:

- Users and roles
- Inventory items
- Tags and tag-to-item relationships
- Authenticated access with JWT
- CORS-friendly API behavior for frontend clients

## Features

- User registration and login
- JWT-based authentication
- Profile lookup for the authenticated user
- User listing and management
- Role-based access control
- Inventory item creation, retrieval, update, and deletion
- Search by item name
- Tag creation, retrieval, update, and deletion
- Inventory filtering by tag
- Health check endpoint
- Seeded sample data for local development

## Tech stack

- Go
- PostgreSQL
- GORM
- JWT
- Docker Compose

## Project structure

```text
cmd/
  api/        API server entry point
  db/         Database migration and seed utilities
internal/
  config/     Environment and app configuration
  database/   PostgreSQL connection setup
  handler/    HTTP handlers and request validation
  middleware/ Auth, CORS, and shared middleware
  model/      Database models and DTOs
  repository/ Database access layer
  service/    Business logic layer
  utility/    Common response helpers
```

## Requirements

- Go 1.27.1 or newer (see `go.mod`)
- Docker and Docker Compose
- PostgreSQL (provided by the local Compose setup)

## Configuration

Create a `.env` file in the project root before starting the application. Example:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ciptakarya
DB_TIMEZONE=Asia/Jakarta
JWT_SECRET=change_this_to_a_secure_secret
JWT_EXPIRE_TIME=30m
APP_PORT=:8080
ALLOWED_ORIGINS=http://localhost:5173
ALLOWED_METHODS=GET, POST, PUT, PATCH, DELETE, OPTIONS
ALLOWED_HEADERS=Content-Type, Authorization

# Optional pgAdmin configuration for docker compose
PGADMIN_EMAIL=admin@example.com
PGADMIN_PASSWORD=your_pgadmin_password
```

Important:

- Do not commit production credentials or real secrets.
- `APP_PORT` should include the leading colon, for example `:8080`.
- If you override `ALLOWED_METHODS`, make sure `PATCH` is included for account and inventory updates.

## Local database setup

Start PostgreSQL with Docker Compose:

```bash
docker compose up -d postgres
```

To also run pgAdmin locally:

```bash
docker compose up -d
```

## Database migration and seed

Before running the API, initialize the tables:

```bash
go run ./cmd/db migrate
```

Optional sample data:

```bash
go run ./cmd/db seed
```

The seed script creates default roles and sample users, including:

- `admin@example.com` / `Admin123!`
- `crew@example.com` / `Crew123!`

## Run the API

```bash
go run ./cmd/api
```

The server listens on the value of `APP_PORT` (default is `:8080`).

A basic health check is available at:

```bash
curl http://localhost:8080/health
```

## API routes

| Method | Route | Access | Description |
| --- | --- | --- | --- |
| `GET` | `/health` | Public | Server health check |
| `POST` | `/reg` | Public | Register a new user |
| `POST` | `/login` | Public | Authenticate and receive a JWT |
| `GET` | `/users` | Authenticated + admin | List all users |
| `GET` | `/users/me` | Authenticated | Get the current user's profile |
| `PATCH` | `/users/{id}` | Authenticated | Update a user record |
| `DELETE` | `/users/{id}` | Authenticated + admin | Delete a user |
| `DELETE` | `/roles/{id}` | Authenticated + admin | Delete a role |
| `POST` | `/items` | Authenticated + admin | Create an item |
| `GET` | `/items` | Authenticated | List items |
| `GET` | `/items/{name}` | Authenticated | Find an item by name |
| `PATCH` | `/items/{id}` | Authenticated + admin | Update an item |
| `DELETE` | `/items/{id}` | Authenticated + admin | Delete an item |
| `POST` | `/tags` | Authenticated + admin | Create a tag |
| `GET` | `/tags` | Authenticated | List tags |
| `GET` | `/tags/{name}` | Authenticated | Find a tag by name |
| `GET` | `/tags/{name}/items` | Authenticated | List items associated with a tag |
| `DELETE` | `/tags/{id}` | Authenticated + admin | Delete a tag |

## Example usage

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin123!"
  }'
```

Use the returned JWT in the `Authorization` header:

```bash
curl http://localhost:8080/users/me \
  -H "Authorization: Bearer <your_jwt>"
```

## Build and test

Run the project's Go checks from the repository root:

```bash
go test ./...
go build ./...
```

## Notes

- The API exposes JWT-protected routes and expects a valid bearer token in the `Authorization` header.
- CORS is enabled through the environment variables defined in `.env`.
- This backend is intended to work with a frontend app that consumes the API endpoints listed above.
