# i18n-flow Admin Frontend

A modern React-based admin panel for managing internationalization (i18n) projects and translations.

## Overview

The i18n-flow Admin Frontend provides a user-friendly interface for managing multilingual content across projects. It allows users to:

- Manage i18n projects
- Add and edit translation keys
- Manage translations across multiple languages
- Import/export translation data
- Monitor translation status and progress

## Tech Stack

- **Framework**: React 19
- **UI Library**: Ant Design 5
- **Styling**: Tailwind CSS 4
- **Routing**: React Router 7
- **HTTP Client**: Axios
- **Build Tool**: Vite
- **Language**: TypeScript
- **Data Export**: XLSX

## Features

### Authentication

- Secure login system
- Protected routes for authenticated users

### Project Management

- Create, view, edit, and delete i18n projects
- Project details including name, description, and slug

### Translation Management

- Manage translation keys and their values across multiple languages
- Filter and search for specific translations
- Batch operations for translation updates
- Context support for translations

### Dashboard

- Overview of translation progress
- Project statistics and status

## Getting Started

### Prerequisites

- Node.js (latest LTS version recommended)
- pnpm

### Installation

```bash
# Clone the repository
git clone [repository-url]

# Navigate to the project directory
cd i18n-flow/admin-frontend

# Install dependencies
pnpm install
```

### Development

```bash
# Start the development server
pnpm dev
```

The application will be available at <http://localhost:5173> by default.

### Building for Production

```bash
# Build the application
pnpm build

# Preview the production build
pnpm preview
```

## Project Structure

```
src/
├── assets/       # Static assets like images
├── components/   # Reusable UI components
├── contexts/     # React context providers
├── hooks/        # Custom React hooks
├── pages/        # Application pages
├── services/     # API service integrations
├── types/        # TypeScript type definitions
└── utils/        # Utility functions
```

## API Integration

The frontend connects to a backend API for data operations. The API endpoints are abstracted through service modules:

- `projectService`: Manages project CRUD operations
- `translationService`: Handles translation-related operations
- `dashboardService`: Provides dashboard statistics and data

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request
