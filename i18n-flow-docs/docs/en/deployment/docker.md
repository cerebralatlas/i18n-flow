# Docker Deployment

This guide will help you deploy the i18n Flow system using Docker and Docker Compose, which is the recommended deployment method for both development and production environments.

## System Requirements

Before you begin, make sure your server meets the following requirements:

- **Operating System**: Linux, macOS, or Windows
- **Docker**: Version 20.10+
- **Docker Compose**: v2+
- **Available Memory**: Minimum 2GB RAM (4GB+ recommended)
- **CPU**: Minimum 2 cores
- **Storage**: Minimum 10GB available disk space
- **Network**: Access to external network for downloading Docker images

## Quick Deployment

### Step 1: Get the Code

First, clone the i18n Flow repository to your server:

```bash
git clone https://github.com/ilukemagic/i18n-flow.git
cd i18n-flow
```

### Step 2: Configure Environment Variables

Copy the environment variable template and configure it according to your needs:

```bash
cp .env.example .env
```

Edit the `.env` file with your preferred text editor:

```bash
nano .env
```

The main environment variables to configure are:

```properties
# Database Configuration
DB_ROOT_PASSWORD=secure_root_password  # MySQL root password
DB_USERNAME=i18nflow                   # Database username
DB_PASSWORD=secure_password            # Database password
DB_NAME=i18n_flow                      # Database name

# JWT Configuration
JWT_SECRET=your_secure_jwt_secret                 # JWT secret key
JWT_EXPIRATION_HOURS=24                           # JWT expiration time (hours)
JWT_REFRESH_SECRET=your_secure_refresh_secret     # JWT refresh secret key
JWT_REFRESH_EXPIRATION_HOURS=168                  # JWT refresh expiration time (hours)

# CLI API Configuration
CLI_API_KEY=your_secure_api_key                   # API key for CLI tool
```

::: warning Security Note
In production environments, be sure to use strong passwords and keys, and keep your `.env` file secure.
:::

### Step 3: Start the Services

Use Docker Compose to start all services:

```bash
docker compose up -d
```

This will start three services:

1. **MySQL Database**: For storing all i18n data
2. **Backend API Server**: Providing the RESTful API
3. **Frontend Admin Interface**: Providing the user interface

The first startup may take several minutes, as Docker needs to download images and build containers.

### Step 4: Verify Deployment

After startup, you can access the following addresses to verify that the deployment was successful:

- **Admin Interface**: `http://your-server-ip` (or `http://localhost` if deployed locally)
- **API Documentation**: `http://your-server-ip/swagger/index.html`

Initial login credentials:

- **Username**: `admin`
- **Password**: `admin123`

::: danger Important Note
After your first login, immediately change the default password to ensure system security!
:::

## Docker Compose Configuration Explained

The i18n Flow Docker Compose configuration file (`docker-compose.yml`) defines three services:

### Database Service

```yaml
db:
  image: mysql:8.0
  container_name: i18n_flow_db
  restart: unless-stopped
  env_file:
    - .env
  environment:
    MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD:-rootpassword}
    MYSQL_DATABASE: ${DB_NAME:-i18n_flow}
    MYSQL_USER: ${DB_USERNAME:-i18nflow}
    MYSQL_PASSWORD: ${DB_PASSWORD:-password}
  ports:
    - "3306:3306"
  volumes:
    - mysql_data:/var/lib/mysql
  networks:
    - i18n_flow_network
  healthcheck:
    test:
      [
        "CMD",
        "mysqladmin",
        "ping",
        "-h",
        "localhost",
        "-u",
        "root",
        "-p${DB_ROOT_PASSWORD:-rootpassword}",
      ]
    interval: 10s
    timeout: 5s
    retries: 5
```

### Backend Service

```yaml
backend:
  build:
    context: ./admin-backend
    dockerfile: Dockerfile
  container_name: i18n_flow_backend
  restart: unless-stopped
  depends_on:
    db:
      condition: service_healthy
  env_file:
    - .env
  environment:
    DB_USERNAME: ${DB_USERNAME:-i18nflow}
    DB_PASSWORD: ${DB_PASSWORD:-password}
    DB_HOST: db
    DB_PORT: 3306
    DB_NAME: ${DB_NAME:-i18n_flow}
    JWT_SECRET: ${JWT_SECRET:-your_secure_jwt_secret}
    JWT_EXPIRATION_HOURS: ${JWT_EXPIRATION_HOURS:-24}
    JWT_REFRESH_SECRET: ${JWT_REFRESH_SECRET:-your_secure_refresh_secret}
    JWT_REFRESH_EXPIRATION_HOURS: ${JWT_REFRESH_EXPIRATION_HOURS:-168}
    CLI_API_KEY: ${CLI_API_KEY:-your_secure_api_key}
  ports:
    - "8080:8080"
  networks:
    - i18n_flow_network
```

### Frontend Service

```yaml
frontend:
  build:
    context: ./admin-frontend
    dockerfile: Dockerfile
  container_name: i18n_flow_frontend
  restart: unless-stopped
  depends_on:
    - backend
  env_file:
    - .env
  ports:
    - "80:80"
  networks:
    - i18n_flow_network
```

## Custom Configuration

### Changing Port Mappings

If you need to change the default port mappings, you can edit the `ports` section in the `docker-compose.yml` file:

```yaml
# Change frontend service port (from 80 to 8000)
frontend:
  ports:
    - "8000:80"

# Change backend API port (from 8080 to 9000)
backend:
  ports:
    - "9000:8080"
```

### Using an External Database

If you want to use an external database instead of the containerized one, you can:

1. Modify the database configuration in the `.env` file:

```properties
DB_HOST=your-external-db-host
DB_PORT=3306
DB_USERNAME=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=your-db-name
```

2. Remove the database service and dependency in `docker-compose.yml`:

```yaml
# Remove db service definition
# Modify backend service's depends_on section
backend:
  depends_on: [] # Remove db dependency
```

### Setting Up HTTPS

In production environments, it's recommended to use HTTPS. You can configure it in one of two ways:

1. **Using a Reverse Proxy** (recommended): Configure Nginx or Traefik as a reverse proxy in front, using Let's Encrypt certificates

2. **Updating Frontend Nginx Configuration**: Modify the `admin-frontend/nginx.conf` file, add SSL configuration, and update the frontend service configuration in Docker Compose to mount certificate files

## Data Persistence

The Docker Compose configuration uses a named volume `mysql_data` to persist database data. This ensures that data remains even after containers are removed.

If you need to backup data, you can:

1. Use Docker commands to backup the database:

```bash
docker exec i18n_flow_db mysqldump -u root -p<root-password> i18n_flow > i18n_flow_backup.sql
```

2. Backup the Docker volume:

```bash
docker run --rm -v i18n-flow_mysql_data:/source -v $(pwd):/backup alpine tar -czvf /backup/mysql_data_backup.tar.gz /source
```

## Updating the System

To update the i18n Flow system to a new version:

```bash
# Pull the latest code
git pull

# Rebuild and start services
docker compose down
docker compose up -d --build
```

## Monitoring and Logging

### View Service Status

```bash
docker compose ps
```

### View Logs

```bash
# View logs for all services
docker compose logs

# View logs for a specific service (e.g. backend)
docker compose logs backend

# View logs in real-time
docker compose logs -f backend
```

### Monitor Container Resource Usage

```bash
docker stats
```

## Troubleshooting

### Service Startup Failure

If a service fails to start, check the logs:

```bash
docker compose logs <service-name>
```

### Database Connection Issues

If the backend cannot connect to the database:

1. Confirm that the database service is healthy: `docker compose ps db`
2. Check database connection settings: `docker compose logs backend | grep "database"`
3. Try restarting the service: `docker compose restart backend`

### Network Issues

If services cannot communicate with each other:

1. Check if the network is correctly created: `docker network ls`
2. Verify that containers are on the same network: `docker network inspect i18n_flow_network`

### Permission Issues

If you encounter permission errors:

```bash
# Ensure data volumes have correct permissions
docker compose down
docker volume rm i18n-flow_mysql_data
docker compose up -d
```

## Performance Optimization

For production environments, consider the following optimizations:

1. **Increase Database Resources**: Configure more resources for the database service in `docker-compose.yml`

```yaml
db:
  deploy:
    resources:
      limits:
        cpus: "2"
        memory: 2G
```

2. **Enable Database Caching**: Add MySQL cache configuration in the `.env` file

3. **Use Redis Caching**: Add a Redis service for caching

4. **Configure Load Balancing**: Add a load balancer in front and deploy multiple backend and frontend service instances

## Multi-Environment Deployment

For different environments (development, testing, production), you can create different environment files and Docker Compose configurations:

```bash
# Development environment
cp .env.example .env.dev
cp docker-compose.yml docker-compose.dev.yml

# Testing environment
cp .env.example .env.test
cp docker-compose.yml docker-compose.test.yml

# Production environment
cp .env.example .env.prod
cp docker-compose.yml docker-compose.prod.yml
```

Then use the `-f` option to specify which configuration file to use:

```bash
# Development environment
docker compose -f docker-compose.dev.yml --env-file .env.dev up -d

# Testing environment
docker compose -f docker-compose.test.yml --env-file .env.test up -d

# Production environment
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
```
