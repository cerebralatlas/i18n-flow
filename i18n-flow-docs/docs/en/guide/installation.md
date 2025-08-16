# Installation Guide

This guide will help you install and set up all components of i18n Flow, including the backend server, frontend admin interface, and CLI tool.

## System Requirements

Before you begin, ensure your system meets the following requirements:

- **Operating System**: Linux, macOS, or Windows
- **Go**: 1.20 or higher
- **Node.js**: 18 or higher
- **MySQL**: 8.0 or higher
- **Package Manager**: pnpm (recommended) or npm
- **Docker** (optional, for containerized deployment): 20.10+ and Docker Compose v2+

## Basic Installation

### 1. Clone the Repository

First, clone the i18n Flow repository from GitHub:

```bash
git clone https://github.com/cerebralatlas/i18n-flow.git
cd i18n-flow
```

### 2. Backend Installation

Navigate to the backend directory and set up the environment:

```bash
cd admin-backend

# Create environment configuration file
cp .env.example .env

# Edit the .env file to configure database connection and other settings
# Modify the following according to your environment:
# - Database credentials
# - JWT keys
# - API keys
# - Other configuration options
```

Install dependencies and start the backend service:

```bash
go mod download
go run main.go
```

The backend API will be available at `http://localhost:8080`, and Swagger documentation at `http://localhost:8080/swagger/index.html`.

### 3. Frontend Installation

Navigate to the frontend directory and install dependencies:

```bash
cd ../admin-frontend
pnpm install
```

Start the development server:

```bash
pnpm dev
```

The admin interface will be available at `http://localhost:5173`.

### 4. CLI Tool Installation

The CLI tool can be installed globally for use in any project:

```bash
npm install -g i18n-flow-cli
```

## Docker Deployment

i18n Flow can be quickly deployed using Docker and Docker Compose:

```bash
# In the project root directory
cp .env.example .env

# Edit the .env file to configure necessary environment variables

# Start the services
docker compose up -d
```

This will start three services:

- MySQL database (port 3306)
- Backend API (port 8080)
- Frontend admin panel (port 80)

### Accessing the Application

- Admin Interface: http://localhost
- API and Swagger Documentation: http://localhost/swagger/index.html

### Initial Login

- Username: `admin`
- Password: `admin123`
- **Important Note**: Change the default password after your first login!

## Development Environment Setup

If you want to contribute to development, follow these steps to set up your development environment:

### Backend Development

The backend is developed using Go and the Gin framework:

```bash
cd admin-backend

# Set up development environment
go mod download

# Run tests
go test ./...

# Start development server
go run main.go
```

### Frontend Development

The frontend is developed using React, TypeScript, and Vite:

```bash
cd admin-frontend

# Install dependencies
pnpm install

# Start development server
pnpm dev

# Build production version
pnpm build
```

### CLI Tool Development

The CLI tool is developed using TypeScript/JavaScript and Node.js:

```bash
cd i18n-flow-cli

# Install dependencies
pnpm install

# Build the tool
pnpm build

# Link globally
npm link
```

## Next Steps

- [Frontend Guide](/en/guide/frontend-guide) - Learn how to use the admin interface
- [Backend Guide](/en/guide/backend-guide) - Understand the backend API and configuration
- [CLI Tool Guide](/en/guide/cli-guide) - Learn how to use the command-line tool
- [Docker Deployment](/en/deployment/docker) - Detailed Docker deployment instructions
