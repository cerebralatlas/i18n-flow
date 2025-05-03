# i18n Flow Documentation Site

[简体中文](README-zh.md) | English

This is the official documentation site for the i18n Flow project, built with VitePress and supporting both English and Chinese languages.

## Project Structure

- `docs/`: Documentation content
  - `en/`: English documentation
  - `zh/`: Chinese documentation
- `Dockerfile`: Docker image for the documentation site
- `nginx.conf`: Nginx configuration file

## Development Environment

### Local Development

```bash
# Install dependencies
pnpm install

# Start development server
pnpm dev

# Build static site
pnpm build

# Preview build result
pnpm preview
```

## Docker Deployment

### Integration into the main project docker-compose.yml

Add the following service configuration to the `docker-compose.yml` file in the project root directory:

```yaml
# Documentation site service
docs:
  build:
    context: ./i18n-flow-docs
    dockerfile: Dockerfile
  container_name: i18n_flow_docs
  restart: unless-stopped
  ports:
    - "8000:80" # Documentation site port mapped to host port 8000
  volumes:
    - docs_nginx_logs:/var/log/nginx
  networks:
    - i18n_flow_network
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:80"]
    interval: 30s
    timeout: 10s
    retries: 3
    start_period: 40s
```

Also, add to the `volumes` section:

```yaml
volumes:
  # Other existing volumes
  docs_nginx_logs:
```

After modification, use the following command to start the entire application:

```bash
# Execute in the project root directory
docker compose up -d
```

The documentation site will be accessible at `http://localhost:8000`.

## Accessing the Documentation Site

After deployment, you can access the documentation site through:

- **URL**: `http://localhost:8000`

### Available Content

- **English Documentation**: `http://localhost:8000/en/`
- **Chinese Documentation**: `http://localhost:8000/zh/`

### Documentation Sections

- Getting Started: Introduction to i18n Flow
- Installation Guide: Setup instructions
- Frontend Guide: Using the admin interface
- Backend Guide: API and configuration details
- CLI Guide: Command-line tool usage
- Usage Tutorial: Complete workflow examples

You can switch between English and Chinese versions using the language selector in the upper right corner of the site.

## Custom Configuration

### Modifying the Port

If you need to modify the port, edit the port mapping of the documentation service in the main project's `docker-compose.yml` file:

```yaml
ports:
  - "custom_port:80" # For example "9000:80"
```

### Using a Custom Domain

1. Modify the `server_name` directive in the `nginx.conf` file:

```nginx
server_name docs.your-domain.com;
```

2. Set up domain resolution to the server IP

3. Configure HTTPS (recommended for production environments):
   - Obtain a free certificate using Let's Encrypt
   - Configure Nginx SSL settings

---
