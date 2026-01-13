# Docker 部署

使用 Docker 容器化部署 YFlow。

## 前置条件

- Docker Engine 20.10+
- Docker Compose V2

### 使用 Docker Compose

#### 1. 准备环境变量

项目根目录提供了可直接使用的 `docker-compose.yml`，先创建 `.env` 文件：

```env
# 复制 .env.template 为 .env，并至少修改以下值

# MySQL
DB_ROOT_PASSWORD=your_secure_password
DB_NAME=yflow
DB_USERNAME=yflow
DB_PASSWORD=your_secure_password

# JWT（至少 32 位）
JWT_SECRET=your_jwt_secret_min_32_chars
JWT_REFRESH_SECRET=your_refresh_secret_min_32_chars

# CLI API Key（至少 16 位）
CLI_API_KEY=your_api_key_min_16_chars

# 可选：部署域名（Caddy 反代使用）
# DOMAIN=example.com

# 可选：前端请求 API 的 baseURL（默认 /api，经 Caddy 反代到 backend）
VITE_API_URL=/api
```

你也可以直接参考并编辑仓库内的 `.env.template`。

#### 2. 启动服务

```bash
# 构建并启动所有服务
docker compose up -d --build

# 查看日志
docker compose logs -f

# 停止服务
docker compose down

# 停止并删除数据卷
docker compose down -v
```

#### 3. 访问服务

- 生产部署推荐通过反向代理访问（Compose 默认使用 Caddy，对外暴露 `80/443`）
- 后端 API：`/api/*`（由 Caddy 转发到 `backend:8080`）
- 前端站点：根路径（由 Caddy 转发到 `frontend:8080`）

如果你在本地仅想快速体验（没有真实域名/证书），建议直接按开发模式运行前后端，或根据需要调整 `deploy/nginx/Caddyfile` 的 TLS 配置。

### 生产环境优化

#### 1. 健康检查

```yaml
services:
  backend:
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

#### 2. 资源限制

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

#### 3. 日志管理

```yaml
services:
  backend:
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "5"
```

### 验证部署

#### 检查服务状态

```bash
# 查看运行中的容器
docker compose ps

# 检查健康状态
docker inspect --format='{{.State.Health.Status}}' i18n-backend
```

#### 测试 API

```bash
# 登录 API
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}'
```

### 数据备份

#### 备份 MySQL

```bash
docker exec i18n-mysql mysqldump -u root -p i18n_flow > backup.sql
```

#### 备份 Redis

```bash
docker exec i18n-redis redis-cli BGSAVE
docker cp i18n-redis:/data/dump.rdb ./redis_backup.rdb
```

## 下一步

- [快速开始 →](/guide/quick-start)
- [架构概述 →](/architecture/overview)
