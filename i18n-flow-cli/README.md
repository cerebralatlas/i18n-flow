# i18n-flow CLI

A powerful CLI tool for managing your application's translations with ease.

[![npm version](https://img.shields.io/npm/v/@i18n-flow/cli.svg)](https://www.npmjs.com/package/@i18n-flow/cli)
[![License](https://img.shields.io/npm/l/@i18n-flow/cli.svg)](https://github.com/i18n-flow/cli/blob/main/LICENSE)

## Features

- 🔄 Sync translations between your project and translation server
- 🔍 Scan your source code for translation keys
- 📤 Push new keys to translation server
- 📊 Check translation status for different locales
- 📂 Support for both flat and nested translation file structures

## Installation

```bash
# Install globally
npm install -g @i18n-flow/cli

# Or use with npx
npx @i18n-flow/cli <command>
```

## Quick Start

1. Initialize i18n-flow in your project:

```bash
i18n-flow init
```

2. Sync translations from server:

```bash
i18n-flow sync
```

3. Scan for translation keys and push to server:

```bash
i18n-flow push --scan
```

## Commands

### `init`

Initialize i18n-flow in your project.

```bash
i18n-flow init
```

This interactive setup will guide you through configuring your project with i18n-flow, including:

- Server URL
- API key
- Project ID
- Locales directory
- Default locale

### `sync`

Sync translations from the server to your local project.

```bash
i18n-flow sync [options]
```

Options:

- `-l, --locale <locales>` - Comma-separated list of locales to sync
- `-f, --force` - Force overwrite local translations
- `-n, --nested` - Use nested directory structure (one folder per locale, split by first key part)

### `push`

Push translation keys to the server.

```bash
i18n-flow push [options]
```

Options:

- `-s, --scan` - Scan source files for translation keys
- `-p, --path <patterns>` - Glob patterns for source files (comma-separated)
- `-d, --dry-run` - Show keys without pushing to server
- `-n, --nested` - Use nested directory structure (read from folders per locale)

### `status`

Check translation status for your project.

```bash
i18n-flow status
```

Shows the status of your translations, including:

- Number of keys per locale
- Missing keys from server
- Extra keys not on server

## File Structure

### Traditional (Flat) Structure

By default, i18n-flow organizes translations in a flat structure:

```
locales/
  en.json
  fr.json
  zh.json
  ...
```

Each file contains all keys for that locale:

```json
{
  "account.info": "Account Information",
  "account.settings": "Account Settings",
  "dashboard.title": "Dashboard"
}
```

### Nested Structure

With the `-n, --nested` option, i18n-flow organizes translations in a nested structure:

```
locales/
  en/
    account.json
    dashboard.json
  fr/
    account.json
    dashboard.json
  ...
```

Each namespace file contains only relevant keys:

```json
// locales/en/account.json
{
  "info": "Account Information",
  "settings": "Account Settings"
}

// locales/en/dashboard.json
{
  "title": "Dashboard"
}
```

## Configuration

Configuration is stored in `.i18nflowrc.json` at the root of your project:

```json
{
  "serverUrl": "https://your-i18n-server.com",
  "apiKey": "your-api-key",
  "projectId": "your-project-id",
  "localesDir": "./src/locales",
  "defaultLocale": "en",
  "sourcePatterns": [
    "src/**/*.{js,jsx,ts,tsx}",
    "!src/**/*.{spec,test}.{js,jsx,ts,tsx}",
    "!**/node_modules/**"
  ],
  "extractorPattern": "(?:t|i18n\\.t)\\(['\"]([\\w\\.\\-]+)['\"]"
}
```

## Usage Examples

### Initialize a new project

```bash
i18n-flow init
```

### Sync translations with nested structure

```bash
i18n-flow sync --nested
```

### Push keys from nested structure

```bash
i18n-flow push --nested
```

### Scan specific files for keys

```bash
i18n-flow push --scan --path "src/components/**/*.tsx,src/pages/**/*.tsx"
```

### Check status of translations

```bash
i18n-flow status
```

## License

MIT
