# CLI Tool Guide

The i18n Flow CLI is a powerful command-line tool for managing your application translations, enabling seamless integration between your development workflow and the translation management system.

## Features

- 🔄 Sync translations between your project and the translation server
- 🔍 Scan source code for translation keys
- 📤 Push new keys to the translation server
- 📊 Check translation status across different locales
- 📂 Support for both flat and nested translation file structures

## Installation

The i18n Flow CLI can be installed via npm or yarn:

```bash
# Global installation with npm
npm install -g i18n-flow-cli

# Or with yarn
yarn global add i18n-flow-cli

# Or with pnpm
pnpm add -g i18n-flow-cli

# Or use directly with npx
npx i18n-flow-cli <command>
```

## Quick Start

### 1. Initialize Your Project

First, you need to initialize i18n Flow in your project:

```bash
i18n-flow init
```

This will start an interactive setup wizard that guides you through the configuration process, including:

- Server URL
- API key
- Project ID
- Locales directory
- Default language

After completion, the configuration will be saved in a `.i18nflowrc.json` file in your project's root directory.

### 2. Sync Translations from Server

After initialization, you can sync translations from the server:

```bash
i18n-flow sync
```

This will download all translations from the server and save them to your configured locales directory.

### 3. Push New Translation Keys

You can scan your source code and push new translation keys to the server:

```bash
i18n-flow push --scan
```

## Command Reference

### `init`

Initialize i18n Flow configuration.

```bash
i18n-flow init
```

This interactive setup will guide you through configuring your project, including:

- Server URL
- API key
- Project ID
- Locales directory
- Default language

### `sync`

Sync translations from the server to your local project.

```bash
i18n-flow sync [options]
```

Options:

- `-l, --locale <locales>` - Locales to sync (comma-separated)
- `-f, --force` - Force overwrite local translations
- `-n, --nested` - Use nested directory structure (one folder per language, split by first part of keys)

Examples:

```bash
# Sync all languages
i18n-flow sync

# Sync specific languages
i18n-flow sync -l en,zh,ja

# Force overwrite local translations
i18n-flow sync -f

# Sync with nested structure
i18n-flow sync -n
```

### `push`

Push translation keys to the server.

```bash
i18n-flow push [options]
```

Options:

- `-s, --scan` - Scan source files for translation keys
- `-p, --path <patterns>` - Glob patterns for source files (comma-separated)
- `-d, --dry-run` - Show keys but don't push to server
- `-n, --nested` - Use nested directory structure (read from each language folder)

Examples:

```bash
# Push existing translation keys
i18n-flow push

# Scan source code and push new keys
i18n-flow push --scan

# Scan specific files
i18n-flow push -s -p "src/components/**/*.tsx,src/pages/**/*.tsx"

# Push with nested structure
i18n-flow push -n

# Show scan results without pushing
i18n-flow push -s -d
```

### `status`

Check the translation status of your project.

```bash
i18n-flow status
```

Displays status information including:

- Number of keys per language
- Keys missing on the server
- Keys missing locally
- Completion rate for each language

## File Structure

The i18n Flow CLI supports two translation file structures: flat and nested.

### Flat Structure (Default)

By default, i18n Flow organizes translations in a flat structure:

```
locales/
  en.json
  zh.json
  fr.json
  ...
```

Each file contains all keys for that language:

```json
{
  "common.buttons.save": "Save",
  "common.buttons.cancel": "Cancel",
  "user.profile.title": "User Profile"
}
```

### Nested Structure

When using the `-n, --nested` option, i18n Flow organizes translations in a nested structure:

```
locales/
  en/
    common.json
    user.json
  zh/
    common.json
    user.json
  ...
```

Each namespace file contains only relevant keys:

```json
// locales/en/common.json
{
  "buttons": {
    "save": "Save",
    "cancel": "Cancel"
  }
}

// locales/en/user.json
{
  "profile": {
    "title": "User Profile"
  }
}
```

## Configuration

Configuration is stored in a `.i18nflowrc.json` file in your project's root directory:

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
  "extractorPattern": "(?:t|i18next\\.t|useTranslation\\(\\)[^}]*\\.)t)\\(['\"]([\\w\\.\\-]+)['\"]"
}
```

### Configuration Options

- `serverUrl`: URL of the i18n Flow server
- `apiKey`: API key for CLI authentication
- `projectId`: Project ID or slug
- `localesDir`: Local directory for translation files
- `defaultLocale`: Default language code
- `sourcePatterns`: Source file patterns to scan (supports glob)
- `extractorPattern`: Regular expression to extract translation keys

## Practical Examples

### Example 1: Integrating with a React Project

Suppose you have a React project using react-i18next for internationalization:

1. Initialize i18n Flow:

```bash
i18n-flow init
```

2. Configure the extractor pattern to match react-i18next usage:

```json
{
  "extractorPattern": "(?:t|i18next\\.t|useTranslation\\(\\)[^}]*\\.)t)\\(['\"]([\\w\\.\\-]+)['\"]"
}
```

3. Sync the latest translations from the server:

```bash
i18n-flow sync
```

4. After adding new translation keys in development:

```bash
i18n-flow push --scan
```

### Example 2: Continuous Integration Workflow

In your CI/CD pipeline, you can automate translation synchronization:

```yaml
# Example GitHub Actions workflow
name: Sync Translations

on:
  schedule:
    - cron: "0 0 * * *" # Run daily

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: "18"
      - run: npm install -g i18n-flow-cli
      - run: i18n-flow sync
      - name: Commit changes
        uses: stefanzweifel/git-auto-commit-action@v4
        with:
          commit_message: "chore: sync translations"
```

## Troubleshooting

### Authentication Failure

If you encounter authentication errors:

1. Verify that your API key is correct
2. Check that your server URL is correct
3. Confirm that the server is running

### Translation Keys Not Detected

If some translation keys are not being detected during scanning:

1. Verify that the extractor pattern (`extractorPattern`) matches the pattern used in your code
2. Check that the source file patterns (`sourcePatterns`) include all relevant files

### File Structure Issues

If the translation file structure is not as expected:

1. Confirm whether you're correctly using the `-n` option
2. Check that the i18n library configuration matches the i18n Flow configuration
