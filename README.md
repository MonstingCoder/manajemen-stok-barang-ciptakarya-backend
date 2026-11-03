# Cipta Karya Inventory Management Backend

A Go backend for managing Cipta Karya's inventory. The project includes application layers for users, roles, inventory items, and tags, with PostgreSQL as its database.

## Current status

This project is under development. The API handlers and services are being built, but the `cmd/api` entry point and HTTP routing are not implemented yet, so the API server cannot currently be started.

## Features

- User registration, login, and user management
- Inventory item creation, listing, search, update, and deletion
- Tag creation, listing, search, and deletion
- Role deletion

## Tech stack

- Go (see `go.mod` for the required Go version)
- PostgreSQL
- GORM
- JWT-based authentication

## Requirements

- Go
- Docker Compose

## Configuration

Create a `.env` file in the project root with the following values:

```env
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=your_database_name
JWT_SECRET=your_jwt_secret

# Required by the optional pgAdmin service in compose.yml
PGADMIN_EMAIL=admin@example.com
PGADMIN_PASSWORD=your_pgadmin_password
```

Do not commit real credentials or secrets.

## Start the database

Start PostgreSQL with Docker Compose:

```bash
docker compose up -d postgres
```

The Compose configuration reads the database credentials from `.env`. To start pgAdmin as well, run:

```bash
docker compose up -d
```

## Build and test

Run the Go checks from the project root:

```bash
go test ./...
go build ./...
```

## Project layout

```text
cmd/api/          API entry point (in progress)
internal/config/  Environment configuration
internal/database/Database connection
internal/handler/ HTTP handlers
internal/model/   Data models and DTOs
internal/repository/Database access
internal/service/ Application logic
```
