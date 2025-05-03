# Getting Started

## What is i18n Flow?

i18n Flow is a modern internationalization management system designed to streamline the translation process for web applications. It provides a centralized platform for managing translations across multiple projects and languages, helping development teams, content creators, and localization experts collaborate more efficiently.

## Key Features

- **Project Management**: Organize translations by projects for better workflow management
- **Multiple Language Support**: Manage unlimited languages with custom language settings
- **Key-Value Translation System**: Maintain consistency in your application translations through key-value pairs
- **Batch Operations**: Import, export, and update translations in bulk to save time
- **Context Support**: Add contextual information to improve translation accuracy
- **Excel Import/Export**: Support for standard formats for easy integration with existing workflows
- **RESTful API**: Well-documented API for integration with your systems
- **User Authentication**: Secure JWT-based authentication system
- **Responsive UI**: Modern admin dashboard built with React and Ant Design
- **CLI Integration**: Command-line tools for seamless integration with development workflows

## System Requirements

- Go 1.20 or higher
- Node.js 18 or higher
- MySQL 8.0
- pnpm (recommended) or npm

## Architecture

i18n Flow consists of three main components:

1. **Admin Backend**: A Go-based API server that manages all i18n data and provides RESTful interfaces
2. **Admin Frontend**: A React-based dashboard for visually managing translation projects
3. **CLI Tool**: Command-line interface for developers to sync translations with codebase

## Workflow

The typical workflow in i18n Flow is as follows:

1. Create a project through the admin dashboard
2. Define the languages you need to support
3. Use the CLI to scan source code for translation keys
4. Push new keys to the translation server using `i18n-flow push --scan`
5. Add translations through the admin interface
6. Sync the latest translations to your project using `i18n-flow sync`
7. Integrate translations in your application using translation functions

## Next Steps

- [Installation Guide](/en/guide/installation) - Learn how to install and configure i18n Flow
- [Frontend Guide](/en/guide/frontend-guide) - Understand how to use the admin interface
- [Backend Guide](/en/guide/backend-guide) - Dive deeper into the API and configuration
- [CLI Tool Guide](/en/guide/cli-guide) - Learn how to use the command-line tool
