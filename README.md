# i18n-flow

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

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

## 📚 Documentation

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

## 🏗️ Architecture

i18n-flow is built with a modern tech stack and follows a clean architecture approach:

### Backend

- **Language**: Go
- **Framework**: Gin
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI

### Frontend

- **Language**: TypeScript
- **Framework**: React
- **UI Library**: Ant Design
- **State Management**: React Context API
- **Routing**: React Router
- **Styling**: Tailwind CSS

## 🔄 Workflow

The typical workflow in i18n-flow:

1. Create a project
2. Define languages you need to support
3. Import existing translation keys or create new ones
4. Translate keys into different languages
5. Export translations to use in your application
6. Update translations as needed

## 🤝 Contributing

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

## 📧 Contact

Project Link: [https://github.com/ilukemagic/i18n-flow](https://github.com/ilukemagic/i18n-flow)

---

Made with ❤️ for better i18n workflows
