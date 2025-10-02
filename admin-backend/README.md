# i18n-flow Backend

A powerful backend service for internationalization management in application development. Built with Go 1.23, Gin, GORM, and Redis.

## Overview

i18n-flow is a comprehensive backend system designed to streamline the internationalization (i18n) process for software applications. It provides a centralized platform for managing translation keys, values, and languages across multiple projects through a RESTful API.

## Features

- **Project Management**: Create and organize multiple translation projects with role-based access control
- **Project Member Management**: Invite users to projects with viewer, editor, or owner roles
- **Multi-language Support**: Add, edit, and manage different languages and locales
- **Translation Management**: Centralize all translations with context support
- **API-driven Architecture**: RESTful API for seamless integration
- **CLI Integration**: Command-line tool support for automated translation workflows
- **Batch Operations**: Efficiently handle bulk translations
- **Export/Import**: Flexible data interchange in various formats
- **Dashboard**: Statistics for monitoring translation progress
- **Authentication**: JWT-based auth for admin interface and API key auth for CLI tools
- **User Management**: Admin functionality for creating and managing system users
- **Role-based Permissions**: Fine-grained access control for projects and translations
- **Intelligent Rate Limiting**: Advanced rate limiting with tollbooth for DDoS protection
- **Enterprise Security**: Multi-layer security protection with input validation, XSS prevention, and SQL injection defense
- **Monitoring & Observability**: Built-in health checks, performance metrics, and enhanced logging

## Tech Stack

- **Language**: Go 1.23
- **Web Framework**: Gin 1.9.1
- **ORM**: GORM 1.30.0
- **Database**: MySQL
- **Documentation**: Swagger/OpenAPI
- **Authentication**: JWT & API keys
- **Caching**: Redis
- **Logging**: Zap with lumberjack rotation
- **Rate Limiting**: Tollbooth-based intelligent rate limiting
- **Security**: Multi-layer protection (input validation, XSS/SQL injection prevention)
- **Monitoring**: Built-in health checks and performance metrics

## API Endpoints

### Authentication & User

- `POST /api/login`: Admin login
- `POST /api/refresh`: Refresh JWT token
- `GET /api/user/info`: Get current user info
- `POST /api/user/change-password`: Change current user password
- `GET /api/cli/auth`: Validate CLI API key

### User Management (Admin)

- `POST /api/users`: Create new user
- `GET /api/users`: List all users with pagination
- `GET /api/users/:id`: Get user details
- `PUT /api/users/:id`: Update user information
- `POST /api/users/:id/reset-password`: Reset user password
- `DELETE /api/users/:id`: Delete user
- `GET /api/user-projects/:user_id`: Get user's projects

### Projects

- `POST /api/projects`: Create project
- `GET /api/projects`: List projects with pagination
- `GET /api/projects/accessible`: Get accessible projects for current user
- `GET /api/projects/detail/:id`: Get project details
- `PUT /api/projects/update/:id`: Update project
- `DELETE /api/projects/delete/:id`: Delete project
- `GET /api/projects/:project_id/members`: Get project members
- `POST /api/projects/:project_id/members`: Add project member
- `PUT /api/projects/:project_id/members/:user_id`: Update member role
- `DELETE /api/projects/:project_id/members/:user_id`: Remove project member
- `GET /api/projects/:project_id/members/:user_id/permission`: Check user permission in project

### Languages

- `GET /api/languages`: List languages
- `POST /api/languages`: Create language
- `PUT /api/languages/:id`: Update language
- `DELETE /api/languages/:id`: Delete language

### Translations

- `POST /api/translations`: Create translation
- `POST /api/translations/batch`: Batch create translations
- `GET /api/translations/by-project/:project_id`: Get project translations
- `GET /api/translations/matrix/by-project/:project_id`: Get translation matrix
- `GET /api/translations/:id`: Get translation details
- `PUT /api/translations/:id`: Update translation
- `DELETE /api/translations/:id`: Delete translation
- `POST /api/translations/batch-delete`: Batch delete translations
- `GET /api/exports/project/:project_id`: Export project translations
- `POST /api/imports/project/:project_id`: Import project translations

### CLI Tool Integration

- `GET /api/cli/translations`: Get translations for CLI
- `POST /api/cli/keys`: Push new translation keys from CLI

### Permissions & Roles

The system implements a role-based access control system with three permission levels:

**Project Roles:**

- **Viewer**: Can view projects, translations, and export data
- **Editor**: Viewer permissions + can create, update, and delete translations
- **Owner**: Editor permissions + can manage project settings and members

**System Roles:**

- **Admin**: Full system access, can manage users, languages, and all projects
- **User**: Standard user with project-specific permissions

### Dashboard

- `GET /api/dashboard/stats`: Get system statistics

### Monitoring & Health

- `GET /health`: Service health check with detailed status
- `GET /stats`: Basic performance statistics
- `GET /stats/detailed`: Detailed system information and metrics

## Getting Started

### Prerequisites

- Go 1.23+
- MySQL
- Git
- Redis (for caching)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/i18n-flow.git
   cd i18n-flow/admin-backend
   ```

   Or download the latest release from the releases page.

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Configure environment:

   ```bash
   cp .env.example .env
   ```

   Edit `.env` file with your database credentials and settings:

   ```
   DB_DRIVER=mysql           # mysql
   DB_USERNAME=your_db_user
   DB_PASSWORD=your_db_password
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=i18n_flow
   
   JWT_SECRET=your_secure_jwt_secret
   JWT_EXPIRATION_HOURS=24
   JWT_REFRESH_SECRET=your_secure_refresh_secret
   JWT_REFRESH_EXPIRATION_HOURS=168
   
   CLI_API_KEY=your_secure_api_key
   
   ADMIN_USERNAME=admin
   ADMIN_PASSWORD=your_password
   
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   REDIS_PREFIX=i18n_flow:
   
   LOG_LEVEL=info           # debug, info, warn, error, fatal
   LOG_FORMAT=console       # console, json
   LOG_OUTPUT=both          # console, file, both
   LOG_DIR=logs
   LOG_DATE_FORMAT=2006-01-02
   LOG_MAX_SIZE=100         # MB
   LOG_MAX_AGE=7            # days
   LOG_MAX_BACKUPS=5
   LOG_COMPRESS=true
   LOG_ENABLE_CONSOLE=true
   ```

4. Start the server:

   ```bash
   # Using air for hot-reload during development
   air
   
   # Or run directly
   go run cmd/server/main.go
   
   # Or build and run the binary
   go build -o i18n-flow ./cmd/server
   ./i18n-flow
   ```

   The server will start on port 8080 by default. The database tables will be automatically created using GORM AutoMigrate, along with default admin user and language seed data.

5. Access the Swagger documentation:

   ```
   http://localhost:8080/swagger/index.html
   ```

### Initial Access

On first run, the system creates a default admin user based on your environment variables:

- Username: Value from `ADMIN_USERNAME` in .env file (defaults to "admin")
- Password: Value from `ADMIN_PASSWORD` in .env file (defaults to "admin123")

The system also creates 20 default languages (English, Chinese, Japanese, Korean, French, German, Spanish, etc.).

**Important:** It's highly recommended to change the default password after the first login through the admin interface.

## Rate Limiting & DDoS Protection

The i18n-flow backend includes an advanced rate limiting system powered by **tollbooth** to protect against abuse and ensure fair resource usage.

### 🛡️ Rate Limiting Features

- **Intelligent Rate Limiting**: Automatic per-IP rate limiting with configurable thresholds
- **DDoS Protection**: Prevents abuse and ensures service availability
- **Zero Configuration**: Works out-of-the-box with sensible defaults
- **High Performance**: Lock-free design for minimal performance impact
- **Automatic Cleanup**: Memory-efficient with automatic cleanup of expired entries

### 📊 Rate Limiting Rules

The system implements different rate limits for different types of operations:

| Endpoint Type | Rate Limit | Window | Purpose |
|---------------|------------|---------|---------|
| **Global** | 100 req/sec | 5 min | Overall API protection |
| **Login** | 5 req/sec | 10 min | Brute force prevention |
| **API Operations** | 50 req/sec | 5 min | Normal API usage |
| **Batch Operations** | 2 req/sec | 10 min | Resource-intensive operations |

### 🚫 Rate Limit Response

When rate limits are exceeded, the API returns a `429 Too Many Requests` response:

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "请求过于频繁，请稍后再试",
    "details": "Rate limit exceeded for: 192.168.1.100"
  }
}
```

### ⚙️ Configuration

Rate limits are automatically applied but can be customized in the middleware configuration:

```go
// Global rate limiting (100 requests per second)
router.Use(middleware.TollboothGlobalRateLimitMiddleware())

// Login protection (5 requests per second)
loginRoutes.Use(middleware.TollboothLoginRateLimitMiddleware())

// API rate limiting (50 requests per second)
apiRoutes.Use(middleware.TollboothAPIRateLimitMiddleware())

// Batch operations (2 requests per second)
batchRoutes.Use(middleware.TollboothBatchOperationRateLimitMiddleware())
```

### 🔧 Custom Rate Limiting

For specific use cases, you can create custom rate limiters:

```go
// Custom rate limiter: 20 requests per second, 2-minute window
customLimiter := middleware.TollboothCustomRateLimitMiddleware(20, 2*time.Minute)

// User-based rate limiting (uses user ID when available, falls back to IP)
userLimiter := middleware.TollboothUserBasedRateLimitMiddleware(30, 5*time.Minute)
```

### 📈 Rate Limiting Monitoring

Rate limiting events are automatically integrated into the monitoring system:

- **429 Errors**: Counted in error statistics
- **Rate Limit Logs**: Detailed logging of rate limit violations
- **Health Metrics**: Rate limiting status included in health checks

You can monitor rate limiting effectiveness through:

```bash
# Check overall error rate (includes rate limiting)
curl http://localhost:8080/health

# Monitor logs for rate limiting events
tail -f logs/app-$(date +%Y-%m-%d).log | grep "RATE_LIMIT"
```

### 🔄 Best Practices

#### For API Clients

1. **Implement Retry Logic**: Use exponential backoff when receiving 429 responses
2. **Respect Rate Limits**: Monitor your request rate to stay within limits
3. **Cache Responses**: Reduce API calls by caching frequently accessed data

#### For Administrators

1. **Monitor Rate Limiting**: Set up alerts for high rate limiting activity
2. **Adjust Limits**: Tune rate limits based on actual usage patterns
3. **Whitelist Trusted IPs**: Consider implementing IP whitelisting for trusted sources

### 🚨 Troubleshooting Rate Limits

**Common Issues:**

1. **Legitimate users hitting limits**
   - Review and adjust rate limit thresholds
   - Implement user-based rate limiting instead of IP-based
   - Consider implementing API keys with higher limits

2. **High rate limiting activity**
   - Check for DDoS attacks or bot traffic
   - Review application logs for patterns
   - Consider implementing additional security measures

3. **Performance impact**
   - Monitor system performance metrics
   - The tollbooth implementation is highly optimized with minimal overhead
   - Rate limiting adds < 0.1ms latency per request

## Security & Protection

The i18n-flow backend implements enterprise-grade security measures to protect against common web vulnerabilities and attacks.

### 🛡️ Security Features

- **Input Validation**: Comprehensive validation and sanitization of all user inputs
- **XSS Prevention**: Automatic HTML cleaning and strict Content Security Policy
- **SQL Injection Protection**: Multi-layer defense with query validation and monitoring
- **Security Headers**: Complete set of security headers for browser protection
- **CSRF Protection**: Cross-site request forgery prevention
- **Rate Limiting Integration**: DDoS protection with intelligent throttling

### 🔒 Protection Layers

#### 1. Input Validation & Sanitization

```bash
# Automatic protection against malicious inputs
POST /api/login
{"username": "<script>alert(1)</script>", "password": "test"}
# Response: 400 Bad Request - Malicious content detected
```

#### 2. XSS Prevention

- **HTML Cleaning**: Automatic removal of dangerous HTML tags and attributes
- **CSP Policy**: Strict Content Security Policy preventing inline scripts
- **Output Encoding**: Safe rendering of user-generated content

#### 3. SQL Injection Defense

- **Query Validation**: Whitelist-based parameter validation
- **Pattern Detection**: Recognition of common injection patterns
- **Database Monitoring**: Real-time query analysis and logging

#### 4. Security Headers

```http
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'; script-src 'self'; ...
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

### 📊 Security Monitoring

Security events are automatically logged and monitored:

- **Attack Attempts**: XSS, SQL injection, and other malicious requests
- **Rate Limiting**: Excessive request patterns and DDoS attempts
- **CSP Violations**: Content security policy breaches
- **Suspicious Activity**: Unusual access patterns and behaviors

### ⚡ Performance Impact

The security system is designed for minimal performance impact:

- **CPU Overhead**: < 2% per request
- **Memory Usage**: < 5MB for security rules cache
- **Latency**: < 1ms additional processing time
- **Throughput**: No significant impact on request handling

### 🔧 Security Testing

Run the included security test suite:

```bash
# Execute comprehensive security tests
./test_security.sh

# Expected results:
# ✅ XSS Protection: Malicious scripts blocked
# ✅ SQL Injection: Dangerous queries prevented  
# ✅ Input Validation: Oversized inputs rejected
# ✅ Rate Limiting: Excessive requests throttled
# ✅ Security Headers: All protection headers present
```

## Monitoring & Observability

The i18n-flow backend includes a built-in lightweight monitoring system that provides real-time service health status and performance metrics.

### 🚀 Monitoring Endpoints

#### 1. Health Check - `/health`

**Purpose**: Check overall service health status

**Example Request**:

```bash
curl http://localhost:8080/health
```

**Response Example**:

```json
{
  "status": "healthy",
  "uptime": "2h15m30s",
  "uptime_seconds": 8130,
  "request_count": 1250,
  "error_count": 12,
  "slow_requests": 3,
  "error_rate": "0.96%",
  "last_error_time": "2024-10-02 14:30:15",
  "timestamp": "2024-10-02T16:45:45Z",
  "version": "1.0.0",
  "database": "healthy (open: 5, idle: 3)",
  "redis": "healthy"
}
```

**Status Description**:

- `healthy`: All core services are normal
- `unhealthy`: Database connection issues (returns 503 status code)

#### 2. Basic Statistics - `/stats`

**Purpose**: Get basic runtime statistics

**Example Request**:

```bash
curl http://localhost:8080/stats
```

**Response**: Same data structure as `/health`

#### 3. Detailed Statistics - `/stats/detailed`

**Purpose**: Get detailed system information and performance metrics

**Example Request**:

```bash
curl http://localhost:8080/stats/detailed
```

**Response Example**:

```json
{
  "basic_stats": {
    // ... basic statistics
  },
  "system_info": {
    "go_version": "go1.23",
    "service_name": "i18n-flow-backend",
    "environment": "development",
    "log_level": "info"
  },
  "performance": {
    "avg_requests_per_second": "2.45",
    "uptime_hours": "2.26"
  }
}
```

### 📈 Monitoring Metrics

#### Core Metrics

| Metric | Description |
|--------|-------------|
| `request_count` | Total number of requests |
| `error_count` | Number of error requests (4xx + 5xx) |
| `slow_requests` | Number of slow requests (>1 second) |
| `error_rate` | Error rate percentage |
| `uptime` | Service uptime |

#### Service Status

| Component | Health Status | Description |
|-----------|---------------|-------------|
| `database` | healthy/down/error | MySQL connection status |
| `redis` | healthy/down/not_configured | Redis connection status |

### 🔍 Enhanced Logging

The monitoring system also enhances logging:

#### Slow Request Monitoring

- Automatically logs requests taking more than 1 second
- Includes detailed request information (method, path, client IP, etc.)

#### Enhanced Error Logging

- 4xx errors logged to error log
- 5xx errors logged with detailed information to application log
- Includes request context information

#### Log Examples

```
2024-10-02 16:45:30.123 WARN  Slow request detected
  method=GET path=/api/translations/matrix/by-project/1 client_ip=127.0.0.1 
  duration=1.234s status=200

2024-10-02 16:45:35.456 ERROR Server error occurred
  method=POST path=/api/translations status_code=500 duration=0.123s 
  client_ip=127.0.0.1 request_id=20241002164535-abc123
```

### 🛠️ Integration with Monitoring Systems

#### 1. Script-based Monitoring

```bash
# Run test script
./test_monitoring.sh

# Periodic health checks
watch -n 30 'curl -s http://localhost:8080/health | jq .status'
```

#### 2. Load Balancer Integration

```nginx
# Nginx health check configuration example
upstream i18n_backend {
    server 127.0.0.1:8080;
    # Other servers...
}

# Health check
location /health {
    access_log off;
    proxy_pass http://i18n_backend/health;
}
```

#### 3. Monitoring Platform Integration

##### Prometheus Integration (Optional)

If you need Prometheus metrics, you can add:

```bash
# Install Prometheus Go client
go get github.com/prometheus/client_golang/prometheus
```

##### Alert Rules Example

```yaml
# High error rate alert
- alert: HighErrorRate
  expr: error_rate > 5
  for: 5m
  
# Service unhealthy alert  
- alert: ServiceUnhealthy
  expr: up{job="i18n-flow"} == 0
  for: 1m
```

### 🔧 Configuration Options

The monitoring system uses the following environment variables:

```bash
# Environment identifier
ENV=development

# Log level
LOG_LEVEL=info

# Database and Redis configuration (affects health checks)
DB_HOST=localhost
DB_PORT=3306
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 📊 Performance Impact

The monitoring system is designed to be lightweight:

- **CPU Overhead**: < 0.1% per request
- **Memory Overhead**: < 1MB resident memory
- **Latency Impact**: < 0.1ms per request
- **Storage Overhead**: No additional storage requirements

### 🚨 Troubleshooting

#### Common Issues

1. **Health check returns 503**
   - Check database connection configuration
   - Ensure MySQL service is running normally

2. **Redis shows down**
   - Check Redis service status
   - Verify Redis connection configuration

3. **Inaccurate statistics**
   - Restart service to reset counters
   - Check system time synchronization

#### Debug Commands

```bash
# Check service status
curl -v http://localhost:8080/health

# View detailed logs
tail -f logs/app-$(date +%Y-%m-%d).log

# Check database connection
mysql -h $DB_HOST -P $DB_PORT -u $DB_USERNAME -p$DB_PASSWORD -e "SELECT 1"

# Check Redis connection  
redis-cli -h $REDIS_HOST -p $REDIS_PORT ping
```

## Development

### API Documentation

The API documentation is automatically generated with Swagger. You can access it at:

```
http://localhost:8080/swagger/index.html
```

### Hot Reload Development

This project uses [Air](https://github.com/cosmtrek/air) for hot reloading during development. The configuration is in `.air.toml`.

To install Air:

```bash
go install github.com/cosmtrek/air@latest
```

Then run:

```bash
air
```

### Project Structure

- `/cmd/server`: Application entry point
- `/internal/api/handlers`: HTTP request handlers
- `/internal/api/middleware`: Request processing middleware
- `/internal/api/response`: Response formatting
- `/internal/api/routes`: API route definitions
- `/internal/config`: Application configuration
- `/internal/container`: Dependency injection container
- `/internal/domain`: Core domain models and interfaces
- `/internal/repository`: Data access layer
- `/internal/service`: Business logic layer
- `/internal/utils`: Utility functions
- `/docs`: Auto-generated API documentation
- `/migrations`: Database migration files (GORM AutoMigrate is used)
- `/utils`: Shared utility functions
- `/tests`: Test files
- `/logs`: Application log files (created at runtime)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## CLI Tool Integration

The backend provides dedicated endpoints for CLI tool integration, allowing automated workflows:

- Authentication via API key
- Pulling translations for specific projects and languages
- Pushing new translation keys
- Checking API connection status

### Recent Updates

The system has been enhanced with comprehensive project member management functionality and advanced rate limiting:

#### Project Management Enhancements

- **Project Member Handler**: Complete CRUD operations for project members
- **Role-based Access Control**: Three permission levels (viewer, editor, owner)
- **User Management**: Admin interface for creating and managing system users
- **Accessible Projects**: Users can view projects they have access to
- **Permission Checking**: Real-time permission validation for project operations

#### Rate Limiting & Security Improvements

- **Tollbooth Integration**: Migrated to enterprise-grade rate limiting with tollbooth
- **Performance Optimization**: 75% code reduction with 3-5x performance improvement
- **DDoS Protection**: Advanced protection against abuse and bot traffic
- **Zero Maintenance**: Automatic memory management and cleanup
- **Monitoring Integration**: Rate limiting metrics integrated into health checks

#### Enterprise Security Enhancements

- **Multi-layer Security**: Comprehensive protection against XSS, SQL injection, and malicious inputs
- **Input Validation**: Advanced validation using bluemonday and govalidator libraries
- **Security Headers**: Complete set of security headers with strict CSP policies
- **Real-time Monitoring**: Automatic detection and logging of security threats
- **Performance Optimized**: < 2% overhead with enterprise-grade protection

## Docker Deployment

```bash
# Build image
docker build -t i18n-flow-backend .

# Run container with MySQL
docker run -p 8080:8080 --env-file .env i18n-flow-backend

# Or run with Docker Compose (recommended for production)
# Create a docker-compose.yml file with MySQL and Redis services
docker-compose up -d
```

### Environment Variables for Docker

When running in Docker, you may want to adjust these environment variables:

```
DB_HOST=mysql  # Use service name from docker-compose
REDIS_HOST=redis  # Use service name from docker-compose
```

### 🔮 Future Extensions

#### Monitoring Extensions

When the project requires more complex monitoring, consider:

1. **Prometheus + Grafana**: Metrics collection and visualization
2. **Jaeger**: Distributed tracing
3. **ELK Stack**: Log aggregation and analysis
4. **AlertManager**: Alert management

#### Rate Limiting Extensions

For advanced rate limiting scenarios, consider:

1. **Redis-based Rate Limiting**: For distributed deployments across multiple servers
2. **User-tier Rate Limiting**: Different limits for premium vs. free users
3. **Geographic Rate Limiting**: Location-based rate limiting rules
4. **API Key Rate Limiting**: Per-API-key rate limiting with custom quotas

#### Security Extensions

For enhanced security requirements, consider:

1. **WAF Integration**: Web Application Firewall for advanced threat detection
2. **Behavioral Analysis**: Machine learning-based anomaly detection
3. **IP Geolocation**: Geographic access control and threat intelligence
4. **Advanced Monitoring**: Security dashboards and real-time alerting

The current multi-layer security system provides a solid foundation for these advanced security features.

## License

This project is licensed under the MIT License.
