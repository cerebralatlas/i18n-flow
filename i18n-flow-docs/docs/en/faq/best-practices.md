# Best Practices

This guide provides recommended best practices for working with i18n Flow to maximize efficiency and maintainability.

## Project Organization

### Translation Keys Structure

- **Use hierarchical namespaces**: Organize keys with namespaces (e.g., `common.buttons.save`, `user.profile.title`)
- **Be consistent**: Follow the same naming pattern throughout your project
- **Use meaningful keys**: Choose descriptive names rather than abstract identifiers

Example of good key structure:

```json
{
  "common": {
    "buttons": {
      "save": "Save",
      "cancel": "Cancel",
      "delete": "Delete"
    },
    "messages": {
      "success": "Operation successful",
      "error": "An error occurred"
    }
  },
  "user": {
    "profile": {
      "title": "User Profile",
      "email": "Email Address"
    }
  }
}
```

### Project Structure

- Consider splitting large applications into multiple i18n Flow projects:
  - One project per large feature area
  - Separate projects for different product lines
  - Shared projects for common components

## Translation Workflow

### Efficient Key Management

1. **Plan key structure in advance**: Design your translation key hierarchy before implementation
2. **Avoid changing keys**: Once in use, changing keys disrupts translations
3. **Add context comments**: Provide clear context for translators
4. **Use placeholders consistently**: Standardize placeholder format (e.g., `{{name}}` or `{name}`)

### Development Process

1. **Scan and push regularly**: Integrate key scanning into your development workflow
2. **Review before pushing**: Verify new keys before pushing to the server
3. **Use continuous integration**: Set up CI jobs to sync translations automatically
4. **Track completion rates**: Monitor translation completion for each language

## Technical Integration

### Code Integration

- **Use consistent translation functions**: Standardize how translations are accessed in code
- **Implement fallback strategy**: Have a plan for missing translations
- **Add type checking**: Use TypeScript or PropTypes for type safety with translation keys
- **Lazy load translations**: Load translations only when needed for better performance

Example of good React implementation:

```jsx
// Using react-i18next
function Button({ type = "save", onClick }) {
  const { t } = useTranslation();

  return (
    <button onClick={onClick} className={`btn btn-${type}`}>
      {t(`common.buttons.${type}`)}
    </button>
  );
}
```

### CLI Configuration

- **Fine-tune extractor pattern**: Customize the regex pattern to match your code style
- **Exclude test files**: Customize source patterns to exclude tests and mocks
- **Use nested structure for large projects**: Switch to nested structure when key count grows

Example of optimized CLI configuration:

```json
{
  "serverUrl": "https://i18n-flow.example.com",
  "apiKey": "your-api-key",
  "projectId": "my-project",
  "localesDir": "./src/i18n/locales",
  "defaultLocale": "en",
  "sourcePatterns": [
    "src/**/*.{js,jsx,ts,tsx}",
    "!src/**/*.{spec,test}.{js,jsx,ts,tsx}",
    "!src/mocks/**",
    "!**/node_modules/**"
  ],
  "extractorPattern": "(?:t|i18n\\.t|useTranslation\\(\\)[^}]*\\.)t)\\(['\"]([\\w\\.\\-]+)['\"]"
}
```

## Translation Content

### Writing for Translators

- **Keep text simple**: Use clear, concise language
- **Avoid idioms and slang**: They're difficult to translate accurately
- **Provide context**: Add descriptions for potentially ambiguous text
- **Include screenshots**: Visuals help translators understand the context

### Quality Control

- **Review translations**: Have native speakers review translations
- **Test in context**: Check translations in the actual UI
- **Watch string length**: Be aware that translations may be longer than English text
- **Verify variable placeholders**: Ensure all placeholders are preserved in translations

## Performance Optimization

### Frontend Performance

- **Split translations by route**: Load translations only for the current route
- **Cache translations**: Implement client-side caching
- **Optimize bundle size**: Remove unused translations from production builds

### Backend Performance

- **Implement caching**: Cache API responses to reduce database load
- **Use batch operations**: Prefer batch API calls over individual requests
- **Monitor API usage**: Track API performance to identify bottlenecks

## CI/CD Integration

### Automated Workflows

- **Automate translation sync**: Include translation sync in your CI pipeline
- **Verify translation completeness**: Fail builds if critical translations are missing
- **Generate translation reports**: Create completion reports during the build process

Example GitHub Actions workflow:

```yaml
name: i18n Workflow

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  i18n-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: "18"

      - name: Install dependencies
        run: npm ci

      - name: Install i18n Flow CLI
        run: npm install -g i18n-flow-cli

      - name: Push new keys
        run: i18n-flow push --scan

      - name: Check translation status
        run: i18n-flow status --min-completion=80

      - name: Sync translations
        run: i18n-flow sync
```

## Versioning Strategy

### Managing Versions

- **Version your translations**: Use version tags for major application releases
- **Archive old keys**: Don't delete old keys, mark them as deprecated
- **Plan for upgrades**: Have a strategy for migrating translations between major versions

By following these best practices, you can maintain an efficient, scalable, and high-quality translation workflow with i18n Flow.
