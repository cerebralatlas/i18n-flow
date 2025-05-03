# i18n-flow

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[简体中文](README-zh.md) | English

A comprehensive internationalization (i18n) management platform designed to streamline the translation workflow for development teams, content creators, and localization specialists.

## 🌟 Features

- **Project Management**: Organize translations by project for better workflow organization
- **Multi-language Support**: Manage translations for unlimited languages with customizable language settings
- **Key-based Translation System**: Maintain consistent translations across your application with a key-based approach
- **Batch Operations**: Import, export, and update translations in batch to save time
- **Context Support**: Add contextual information to improve translation accuracy
- **Excel Import/Export**: Support for standard formats to integrate with existing workflows
- **RESTful API**: Well-documented API for integration with your systems
- **User Authentication**: Secure JWT-based authentication system
- **Responsive UI**: Modern admin dashboard built with React and Ant Design
- **CLI Integration**: Command-line tool for seamless integration with your development workflow

## 📦 Components

i18n-flow consists of three main components:

1. **Admin Backend**: Go-based API server that manages all i18n data and provides RESTful endpoints
2. **Admin Frontend**: React-based dashboard for visual management of translation projects
3. **CLI Tool**: Command-line interface for developers to sync translations with their codebase

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or higher
- Node.js 18 or higher
- MySQL 8.0
- pnpm (preferred) or npm

### Backend Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/ilukemagic/i18n-flow.git
   cd i18n-flow/admin-backend
   ```

2. Configure your environment:

   ```bash
   cp .env.example .env
   # Edit .env with your database credentials and other settings
   ```

3. Run the backend:

   ```bash
   go mod download
   go run main.go
   ```

4. The API will be available at `http://localhost:8080` with Swagger documentation at `http://localhost:8080/swagger/index.html`

### Frontend Setup

1. Navigate to the frontend directory:

   ```bash
   cd ../admin-frontend
   ```

2. Install dependencies:

   ```bash
   pnpm install
   ```

3. Start the development server:

   ```bash
   pnpm dev
   ```

4. The admin interface will be available at `http://localhost:5173`

### CLI Tool Setup

1. Install the CLI tool globally:

   ```bash
   npm install -g i18n-flow-cli
   ```

2. Initialize i18n-flow in your project:

   ```bash
   i18n-flow init
   ```

3. Follow the interactive setup to configure your project with i18n-flow.

## 🐳 Docker Deployment

You can easily deploy i18n-flow using Docker and Docker Compose.

### Prerequisites

- Docker 20.10+
- Docker Compose v2+

### Quick Deployment

1. Clone the repository:

   ```bash
   git clone https://github.com/ilukemagic/i18n-flow.git
   cd i18n-flow
   ```

2. Configure environment variables:

   ```bash
   cp .env.example .env
   ```

   Edit the `.env` file to set:

   - Database credentials
   - JWT secrets
   - API keys
   - Other configuration options

3. Start the services:

   ```bash
   docker compose up -d
   ```

   This will start three services:

   - **MySQL** database on port 3306
   - **Backend API** on port 8080
   - **Frontend admin panel** on port 80

4. Access the application:

   - Admin interface: <http://localhost>
   - API & Swagger docs: <http://localhost/swagger/index.html>

5. Initial login:
   - Username: `admin`
   - Password: `admin123`
   - **Important**: Change the default password after first login!

### Docker Compose Configuration

The default `docker-compose.yml` includes:

- **Database**: MySQL 8.0 with data persistence through Docker volumes
- **Backend**: Go 1.23 API server with MySQL connection
- **Frontend**: React SPA served through Nginx with API proxying

### Custom Configuration

You can adjust the deployment by:

1. Modifying environment variables in the `.env` file
2. Changing port mappings in `docker-compose.yml`
3. Updating Docker build contexts or volumes as needed

### Updating the Application

To update to a new version:

```bash
git pull
docker compose down
docker compose up -d --build
```

### Troubleshooting

- **Database connection issues**: Check your DB credentials in the `.env` file
- **Frontend not loading**: Verify the nginx proxy configuration is correctly pointing to the backend
- **Backend not starting**: Check the backend logs with `docker compose logs backend`

## 📚 Documentation

### Documentation Site

The i18n Flow documentation is available as a dedicated website:

- **URL**: When deployed via Docker, access at `http://localhost:8000/en/`
- **Languages**: Available in both English and Chinese
- **Content**:
  - Comprehensive guides for all components
  - Installation and setup instructions
  - Usage tutorials with practical examples
  - API reference documentation
  - Deployment guides

Access specific language versions:

- English: `http://localhost:8000/en/`
- Chinese: `http://localhost:8000/zh/`

### API Documentation

The API documentation is available via Swagger UI when running the backend server:

- Open `http://localhost:8080/swagger/index.html` in your browser

### Admin Dashboard

The admin dashboard provides an intuitive interface for:

- Managing projects and their translation keys
- Adding and updating language definitions
- Entering and editing translations
- Importing and exporting translations in various formats
- User management and access control
- Monitoring translation status and progress

### CLI Tool

The CLI tool provides the following commands:

- `init`: Initialize i18n-flow in your project
- `sync`: Sync translations from the server to your local project
- `push`: Push translation keys to the server
- `status`: Check translation status for your project

Full CLI documentation is available by running `i18n-flow --help`.

## 🏗️ Architecture

i18n-flow is built with a modern tech stack and follows a clean architecture approach:

### Backend

- **Language**: Go
- **Framework**: Gin
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT & API keys
- **Documentation**: Swagger/OpenAPI

### Frontend

- **Language**: TypeScript
- **Framework**: React 19
- **UI Library**: Ant Design 5
- **State Management**: React Context API
- **Routing**: React Router 7
- **Styling**: Tailwind CSS 4
- **Build Tool**: Vite

### CLI Tool

- **Language**: TypeScript/JavaScript
- **Runtime**: Node.js
- **Distribution**: npm package

## 🔄 Workflow

The typical workflow in i18n-flow:

1. Create a project through the admin dashboard
2. Define languages you need to support
3. Use the CLI to scan your source code for translation keys
4. Push new keys to the translation server with `i18n-flow push --scan`
5. Translate keys through the admin dashboard
6. Sync translations to your local project with `i18n-flow sync`
7. Update translations as needed

## 📜 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework in Go
- [GORM](https://gorm.io/) - ORM library for Go
- [React](https://reactjs.org/) - A JavaScript library for building user interfaces
- [Ant Design](https://ant.design/) - A design system for enterprise-level products
- [Tailwind CSS](https://tailwindcss.com/) - A utility-first CSS framework
- [Vite](https://vitejs.dev/) - Next generation frontend build tool

## 📧 Contact

Project Link: [https://github.com/ilukemagic/i18n-flow](https://github.com/ilukemagic/i18n-flow)

---

Made with ❤️ for better i18n workflows
