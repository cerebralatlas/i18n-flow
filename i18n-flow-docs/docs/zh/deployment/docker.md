# Docker 部署

本指南将帮助您使用 Docker 和 Docker Compose 部署 i18n Flow 系统，这是推荐的部署方式，适用于开发环境和生产环境。

## 系统要求

在开始之前，请确保您的服务器满足以下要求：

- **操作系统**：Linux、macOS 或 Windows
- **Docker**：20.10+ 版本
- **Docker Compose**：v2+ 版本
- **可用内存**：至少 2GB 内存（推荐 4GB+）
- **CPU**：至少 2 核心
- **存储**：至少 10GB 可用磁盘空间
- **网络**：可访问外部网络，用于下载 Docker 镜像

## 快速部署

### 步骤 1: 获取代码

首先，克隆 i18n Flow 仓库到您的服务器：

```bash
git clone https://github.com/ilukemagic/i18n-flow.git
cd i18n-flow
```

### 步骤 2: 配置环境变量

复制环境变量模板并根据您的需求进行配置：

```bash
cp .env.example .env
```

使用您喜欢的文本编辑器编辑 `.env` 文件：

```bash
nano .env
```

主要需要配置以下环境变量：

```properties
# 数据库配置
DB_ROOT_PASSWORD=secure_root_password  # MySQL root 密码
DB_USERNAME=i18nflow                   # 数据库用户名
DB_PASSWORD=secure_password            # 数据库密码
DB_NAME=i18n_flow                      # 数据库名称

# JWT 配置
JWT_SECRET=your_secure_jwt_secret                 # JWT 密钥
JWT_EXPIRATION_HOURS=24                           # JWT 过期时间（小时）
JWT_REFRESH_SECRET=your_secure_refresh_secret     # JWT 刷新密钥
JWT_REFRESH_EXPIRATION_HOURS=168                  # JWT 刷新过期时间（小时）

# CLI API 配置
CLI_API_KEY=your_secure_api_key                   # CLI 工具的 API 密钥
```

::: warning 安全提示
在生产环境中，请务必使用强密码和密钥，并保管好您的 `.env` 文件。
:::

### 步骤 3: 启动服务

使用 Docker Compose 启动所有服务：

```bash
docker compose up -d
```

这将启动三个服务：

1. **MySQL 数据库**：用于存储所有 i18n 数据
2. **后端 API 服务器**：提供 RESTful API
3. **前端管理界面**：提供用户界面

首次启动可能需要几分钟，因为 Docker 需要下载镜像并构建容器。

### 步骤 4: 验证部署

启动后，您可以访问以下地址来验证部署是否成功：

- **管理界面**：`http://your-server-ip`（或 `http://localhost` 如果在本地部署）
- **API 文档**：`http://your-server-ip/swagger/index.html`

初始登录凭据：

- **用户名**：`admin`
- **密码**：`admin123`

::: danger 重要提示
首次登录后，请立即修改默认密码以确保系统安全！
:::

## Docker Compose 配置详解

i18n Flow 的 Docker Compose 配置文件 (`docker-compose.yml`) 定义了三个服务：

### 数据库服务

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

### 后端服务

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

### 前端服务

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

## 自定义配置

### 修改端口映射

如果您需要更改默认端口映射，可以编辑 `docker-compose.yml` 文件中的 `ports` 部分：

```yaml
# 修改前端服务端口（从 80 改为 8000）
frontend:
  ports:
    - "8000:80"

# 修改后端 API 端口（从 8080 改为 9000）
backend:
  ports:
    - "9000:8080"
```

### 使用外部数据库

如果您想使用外部数据库而不是 Docker 容器中的数据库，可以：

1. 修改 `.env` 文件中的数据库配置：

```properties
DB_HOST=your-external-db-host
DB_PORT=3306
DB_USERNAME=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=your-db-name
```

2. 在 `docker-compose.yml` 中移除数据库服务和依赖：

```yaml
# 移除 db 服务定义
# 修改 backend 服务的 depends_on 部分
backend:
  depends_on: [] # 移除 db 依赖
```

### 设置 HTTPS

在生产环境中，建议使用 HTTPS。您可以通过以下两种方式之一进行配置：

1. **使用反向代理**（推荐）：在前面配置 Nginx 或 Traefik 作为反向代理，并使用 Let's Encrypt 证书

2. **更新前端 Nginx 配置**：修改 `admin-frontend/nginx.conf` 文件，添加 SSL 配置，然后更新 Docker Compose 中的前端服务配置，挂载证书文件

## 数据持久化

Docker Compose 配置使用命名卷 `mysql_data` 来持久化数据库数据。这确保即使在容器被删除后，数据仍然存在。

如果您需要备份数据，可以：

1. 使用 Docker 命令备份数据库：

```bash
docker exec i18n_flow_db mysqldump -u root -p<root-password> i18n_flow > i18n_flow_backup.sql
```

2. 备份 Docker 卷：

```bash
docker run --rm -v i18n-flow_mysql_data:/source -v $(pwd):/backup alpine tar -czvf /backup/mysql_data_backup.tar.gz /source
```

## 更新系统

更新 i18n Flow 系统到新版本：

```bash
# 拉取最新代码
git pull

# 重新构建并启动服务
docker compose down
docker compose up -d --build
```

## 监控和日志

### 查看服务状态

```bash
docker compose ps
```

### 查看日志

```bash
# 查看所有服务的日志
docker compose logs

# 查看特定服务的日志（例如后端）
docker compose logs backend

# 实时查看日志
docker compose logs -f backend
```

### 监控容器资源使用情况

```bash
docker stats
```

## 故障排除

### 服务启动失败

如果某个服务无法启动，请检查日志：

```bash
docker compose logs <service-name>
```

### 数据库连接问题

如果后端无法连接到数据库：

1. 确认数据库服务是否健康：`docker compose ps db`
2. 检查数据库连接设置：`docker compose logs backend | grep "database"`
3. 尝试重启服务：`docker compose restart backend`

### 网络问题

如果服务间无法通信：

1. 检查网络是否正确创建：`docker network ls`
2. 验证容器是否在同一网络：`docker network inspect i18n_flow_network`

### 权限问题

如果遇到权限错误：

```bash
# 确保数据卷有正确的权限
docker compose down
docker volume rm i18n-flow_mysql_data
docker compose up -d
```

## 性能优化

对于生产环境，可以考虑以下优化措施：

1. **增加数据库资源**：在 `docker-compose.yml` 中配置更多资源给数据库服务

```yaml
db:
  deploy:
    resources:
      limits:
        cpus: "2"
        memory: 2G
```

2. **启用数据库缓存**：在 `.env` 文件中添加 MySQL 缓存配置

3. **使用 Redis 缓存**：添加 Redis 服务用于缓存

4. **配置负载均衡**：在前面添加负载均衡器，部署多个后端和前端服务实例

## 多环境部署

对于不同环境（开发、测试、生产），您可以创建不同的环境文件和 Docker Compose 配置：

```bash
# 开发环境
cp .env.example .env.dev
cp docker-compose.yml docker-compose.dev.yml

# 测试环境
cp .env.example .env.test
cp docker-compose.yml docker-compose.test.yml

# 生产环境
cp .env.example .env.prod
cp docker-compose.yml docker-compose.prod.yml
```

然后使用 `-f` 选项指定要使用的配置文件：

```bash
# 开发环境
docker compose -f docker-compose.dev.yml --env-file .env.dev up -d

# 测试环境
docker compose -f docker-compose.test.yml --env-file .env.test up -d

# 生产环境
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
```
