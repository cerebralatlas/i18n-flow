# Usage Tutorial

This tutorial will guide you through the complete process of using i18n Flow in real projects, from setup to translation management and integration into your application.

## Preparation

Before you start using i18n Flow, ensure you have completed the following preparation steps:

1. Successfully deployed the i18n Flow system (backend and frontend)
2. Installed the i18n Flow CLI tool
3. Prepared the application project that you want to internationalize

## Basic Workflow

The typical workflow in i18n Flow is as follows:

1. Create a project in the admin panel
2. Invite team members with appropriate roles (viewer, editor, or owner)
3. Define supported languages
4. Use the CLI to scan source code for translation keys
5. Push translation keys to the server
6. Add translation content in the admin panel
7. Sync translations back to your development project
8. Integrate translations into your application

## Step 1: Create a Project

First, you need to create a project in the i18n Flow admin interface:

1. Log in to the admin interface
2. Click on the "Projects" menu
3. Click the "Create Project" button
4. Fill in the project information:
   - **Project Name**: Your application name, e.g., "My Online Store"
   - **Project Description**: Brief description, e.g., "E-commerce platform internationalization content"
   - **Project Identifier (Slug)**: Unique identifier for API and CLI calls, e.g., "my-online-store"
5. Click the "Save" button

## Step 2: Manage Project Members

After creating a project, you can invite team members and assign appropriate roles:

### Understanding Roles and Permissions

i18n Flow implements a role-based access control system with three permission levels:

**Project Roles:**
- **Viewer**: Can view projects, translations, and export data
- **Editor**: Viewer permissions + can create, update, and delete translations
- **Owner**: Editor permissions + can manage project settings and members

**System Roles:**
- **Admin**: Full system access, can manage users, languages, and all projects
- **User**: Standard user with project-specific permissions

### Adding Project Members

1. Navigate to your project in the admin interface
2. Click on the "Members" tab
3. Click the "Add Member" button
4. Enter the user's email or username
5. Select the appropriate role for the user
6. Click "Save" to send the invitation

### Managing Member Permissions

- **Update Role**: Change a member's role as project needs evolve
- **Remove Member**: Remove users who no longer need access
- **Check Permissions**: Verify what actions a specific user can perform

### Best Practices for Team Collaboration

1. **Assign roles based on responsibilities**: Give viewers to stakeholders who only need to review, editors to translators, and owners to project managers
2. **Regular review**: Periodically review member access and adjust roles as needed
3. **Use the accessible projects view**: Team members can easily see all projects they have access to

## Step 3: Add Languages

Next, add the languages your application needs to support:

1. Click on "Languages" in the sidebar menu
2. Click the "Add Language" button
3. Add each language you need to support, for example:
   - English (en)
   - Simplified Chinese (zh-CN)
   - Japanese (ja)
   - Spanish (es)
4. For each language, fill in the following information:
   - **Language Name**: e.g., "Simplified Chinese"
   - **Language Code**: e.g., "zh-CN"
   - **Locale Name**: e.g., "Chinese"
   - **RTL Support**: Select Yes for right-to-left writing languages (like Arabic or Hebrew), otherwise No

## Step 3: Initialize CLI in Your Project

In your development project, initialize the i18n Flow CLI:

1. Open a terminal and navigate to your project directory
2. Run the initialization command:

```bash
i18n-flow init
```

3. Follow the prompts to configure the CLI:
   - **Server URL**: Address of your i18n Flow server, e.g., "http://your-i18n-flow-server.com"
   - **API Key**: CLI API key set in the i18n Flow backend configuration
   - **Project ID**: The slug of the project you just created, e.g., "my-online-store"
   - **Locales Directory**: Directory to store translation files, e.g., "./src/locales"
   - **Default Language**: Primary development language, e.g., "en"

## Step 4: Scan and Push Translation Keys

Now scan your source code to find all text that needs translation, and push these keys to the server:

```bash
i18n-flow push --scan
```

This will:

1. Scan your source code
2. Extract all keys matching the translation pattern
3. Push these keys to the i18n Flow server

::: tip
The CLI uses the regular expression `extractorPattern` to match translation functions. The default pattern is `(?:t|i18n\.t)\(['"]([\\w\\.\\-]+)['"]`, which matches calls like `t('common.hello')` or `i18n.t('user.greeting')`. If your project uses different translation functions, adjust this pattern accordingly.
:::

## Step 5: Add Translation Content

After keys are pushed to the server, you can add translations:

1. In the admin interface, navigate to your project
2. Switch to the "Translations" tab
3. You'll see all keys extracted from your code
4. Add translations for each key in different languages
5. You can edit them individually or use the bulk import feature

### Bulk Import of Translations

For large projects, you can use the bulk import feature:

1. Click the "Import" button
2. Select the file format (JSON, Excel, CSV)
3. Upload your translation file or paste content
4. Select the languages to update
5. Click "Import" to start processing

### Export Translations

You can also export translations to share or backup:

1. Click the "Export" button
2. Select the languages to export
3. Select the file format
4. Click "Export" to download the file

## Step 6: Sync Translations to Your Project

Once translations are ready, sync them to your development project:

```bash
i18n-flow sync
```

This will download all translations from the server and save them to your configured locales directory as JSON files.

::: tip File Structure
The i18n Flow CLI supports two file structures:

- **Flat Structure** (default): One file per language containing all keys
- **Nested Structure**: Use the `-n` option to organize keys by namespace
  :::

## Step 7: Integrate into Your Application

Finally, integrate translations into your application. This depends on your frontend framework and i18n library.

### React Project Example (using react-i18next)

1. Install i18next and react-i18next:

```bash
npm install i18next react-i18next i18next-http-backend i18next-browser-languagedetector
```

2. Create an i18n configuration file (src/i18n.js):

```javascript
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import Backend from "i18next-http-backend";
import LanguageDetector from "i18next-browser-languagedetector";

// Import translation files
import enTranslation from "./locales/en.json";
import zhTranslation from "./locales/zh-CN.json";

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en: { translation: enTranslation },
      "zh-CN": { translation: zhTranslation },
    },
    fallbackLng: "en",
    debug: process.env.NODE_ENV === "development",
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
```

3. Initialize i18n in App.js:

```jsx
import React from "react";
import { useTranslation } from "react-i18next";
import "./i18n";

function App() {
  const { t } = useTranslation();

  return (
    <div className="App">
      <h1>{t("app.title")}</h1>
      <p>{t("app.greeting")}</p>
      <button>{t("common.buttons.save")}</button>
    </div>
  );
}

export default App;
```

## Continuous Integration Process

You can integrate i18n Flow into your CI/CD pipeline:

### Development Phase

1. Developers add new translation keys in the code
2. Developers run `i18n-flow push --scan` to push new keys to the server
3. Translators add translations via the admin interface

### Build Phase

1. CI system runs `i18n-flow sync` before building to get the latest translations
2. Application is built with the latest translations

### Example GitHub Actions Workflow

```yaml
name: Build with translations

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v2

      - name: Setup Node.js
        uses: actions/setup-node@v2
        with:
          node-version: "18"

      - name: Install dependencies
        run: npm ci

      - name: Install i18n Flow CLI
        run: npm install -g i18n-flow-cli

      - name: Sync translations
        run: i18n-flow sync
        env:
          I18N_FLOW_API_KEY: ${{ secrets.I18N_FLOW_API_KEY }}

      - name: Build application
        run: npm run build
```

## Advanced Usage

### Namespace-based Organization of Translations

For large projects, you might want to organize translations by namespace:

1. Configure nested structure during initialization:

```bash
i18n-flow init
# Choose nested structure in the setup wizard
```

2. Or use the nested option in an already configured project:

```bash
i18n-flow sync -n
i18n-flow push -n
```

### Handling Key Changes During Code Refactoring

When you refactor code and change translation keys:

1. In the admin interface, use the "Batch Edit" feature to rename keys
2. Update key references in your code
3. Run `i18n-flow push --scan` to update keys on the server
4. Run `i18n-flow sync` to sync changes back to your project

## Troubleshooting

### Keys Not Being Scanned

If some keys are not detected:

1. Check if the key format matches the extraction pattern
2. Verify that file paths are included in the scan patterns
3. Adjust the `extractorPattern` in `.i18nflowrc.json` to match your translation function call style

### Sync Errors

If you encounter sync errors:

1. Verify that the API key and server URL are correct
2. Check network connectivity
3. Verify that the project ID exists on the server

## Best Practices

1. **Use descriptive keys**: Use meaningful key names like `user.profile.title` instead of `title1`
2. **Add context**: Add descriptions for each key in the admin interface to help translators understand the context
3. **Sync regularly**: Regularly sync translations to ensure your development environment has the latest translations
4. **Automate the process**: Integrate translation synchronization into your CI/CD pipeline
5. **Provide default text**: Include default text in your code to display reasonable content when translations are missing
