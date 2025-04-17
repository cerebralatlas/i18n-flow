# CLI Guide

## Installation

```bash
npm install -g i18n-flow-cli
```

## Commands

### Initialize Project

```bash
i18n-flow init
```

### Pull Translations

```bash
i18n-flow pull --project <project-id> --lang <language-code>
```

### Push Translations

```bash
i18n-flow push --project <project-id> --file <translation-file>
```

## Configuration

Create `.i18nflowrc.json` in your project root:

```json
{
  "apiUrl": "http://your-api-url",
  "projectId": "your-project-id",
  "token": "your-api-token"
}
```
