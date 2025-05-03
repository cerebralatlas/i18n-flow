# Troubleshooting

This guide provides solutions to common issues you might encounter when using i18n Flow.

## Authentication Issues

### Unable to Login

**Problem**: Cannot login to the admin interface.

**Solutions**:

- Verify your username and password are correct
- Check if the backend service is running
- Clear browser cookies and cache
- If using Docker, check if all containers are running: `docker compose ps`

### API Key Not Working

**Problem**: CLI commands fail with authentication errors.

**Solutions**:

- Verify the API key in your `.i18nflowrc.json` matches the one in the server's `.env` file
- Confirm the server URL is correct and accessible
- Check network connectivity between your machine and the server

## Backend Issues

### Database Connection Failures

**Problem**: Backend service fails to connect to the database.

**Solutions**:

- Verify database credentials in `.env` are correct
- Check if the database service is running and accessible
- If using Docker, ensure the database container is healthy
- Check database logs for any errors

### Server Won't Start

**Problem**: The backend server won't start.

**Solutions**:

- Check logs for error messages: `docker compose logs backend`
- Verify all environment variables are properly set
- Ensure required ports are not in use by other applications
- Check for sufficient system resources (memory, disk space)

## Frontend Issues

### Blank or Error Screen

**Problem**: Admin interface shows a blank screen or error message.

**Solutions**:

- Check browser console for JavaScript errors
- Verify the backend API is accessible from the browser
- Clear browser cache
- Try a different browser

### UI Rendering Issues

**Problem**: Elements of the UI are not displaying properly.

**Solutions**:

- Update your browser to the latest version
- Disable browser extensions that might interfere
- Try in incognito/private browsing mode

## CLI Tool Issues

### CLI Not Detecting Keys

**Problem**: The CLI tool doesn't detect translation keys when scanning.

**Solutions**:

- Verify the `extractorPattern` in `.i18nflowrc.json` matches your code's translation function calls
- Check that source files are included in the scan patterns
- Run with verbose logging: `i18n-flow push --scan --debug`

### Sync Failures

**Problem**: Cannot sync translations with the server.

**Solutions**:

- Check network connectivity
- Verify project ID and server URL are correct
- Ensure the API key has the necessary permissions
- Check server logs for any API errors

## Translation Management Issues

### Missing Translations

**Problem**: Translations are missing after syncing.

**Solutions**:

- Verify the language code matches between CLI and server
- Check if the translations exist on the server
- Try forcing a sync: `i18n-flow sync --force`
- Check for file permission issues in your locales directory

### Import/Export Failures

**Problem**: Unable to import or export translations.

**Solutions**:

- Verify the file format is supported
- Check that the file structure matches the expected format
- Try with a smaller file to rule out size limitations
- Check server logs for detailed error messages

## Docker Deployment Issues

### Container Startup Failures

**Problem**: Docker containers fail to start.

**Solutions**:

- Check logs: `docker compose logs`
- Verify port availability: `netstat -tuln`
- Ensure Docker has sufficient resources
- Check for errors in your `docker-compose.yml` file

### Volume Permission Issues

**Problem**: Container cannot write to volumes.

**Solutions**:

- Check file ownership and permissions
- Recreate the volume: `docker volume rm [volume_name]`
- Use explicit user mapping in Docker Compose file

## Performance Issues

### Slow API Responses

**Problem**: API requests are taking too long.

**Solutions**:

- Check server load and resources
- Optimize database queries and indexes
- Consider scaling up resources or services
- Enable caching where possible

### High Memory Usage

**Problem**: Services consuming excessive memory.

**Solutions**:

- Check for memory leaks
- Optimize container resource limits
- Consider database query optimizations
- Implement pagination for large datasets

## Getting Help

If you continue to experience issues after trying these troubleshooting steps:

1. Check the [GitHub Issues](https://github.com/ilukemagic/i18n-flow/issues) for similar problems
2. Gather relevant logs and error messages
3. Create a new issue with detailed information about your setup and the problem
4. Join the community discussion on [Discord/Slack] for real-time help
