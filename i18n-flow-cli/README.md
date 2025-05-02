# i18n-flow CLI

A powerful CLI tool for managing your application's translations with ease.

[![npm version](https://img.shields.io/npm/v/@i18n-flow/cli.svg)](https://www.npmjs.com/package/i18n-flow-cli)

[简体中文](README-zh.md) | English

## Features

- 🔄 Sync translations between your project and translation server
- 🔍 Scan your source code for translation keys
- 📤 Push new keys to translation server
- 📊 Check translation status for different locales
- 📂 Support for both flat and nested translation file structures

## Installation

```bash
# Install globally
npm install -g i18n-flow-cli

# Or use with npx
npx i18n-flow-cli <command>
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
