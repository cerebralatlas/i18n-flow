# 后端使用指南

i18n Flow 的后端是一个基于 Go 语言开发的强大 API 服务器，负责管理所有翻译数据，并为前端和 CLI 工具提供接口。本指南将帮助您了解后端服务的配置和使用方法。

## 配置后端

### 环境配置

后端服务使用 `.env` 文件进行配置。默认的配置模板位于 `.env.example` 文件中，您需要复制并修改它：

```bash
cp .env.example .env
```

以下是主要配置选项：

```properties
# 数据库配置
DB_USERNAME=i18nflow           # 数据库用户名
DB_PASSWORD=password           # 数据库密码
DB_HOST=localhost              # 数据库主机
DB_PORT=3306                   # 数据库端口
DB_NAME=i18n_flow              # 数据库名称

# JWT认证配置
JWT_SECRET=your_secure_jwt_secret                 # JWT密钥
JWT_EXPIRATION_HOURS=24                           # JWT过期时间（小时）
JWT_REFRESH_SECRET=your_secure_refresh_secret     # JWT刷新密钥
JWT_REFRESH_EXPIRATION_HOURS=168                  # JWT刷新过期时间（小时）

# CLI API配置
CLI_API_KEY=your_secure_api_key                   # CLI工具的API密钥

# 服务器配置
SERVER_PORT=8080               # 服务器监听端口
SERVER_HOST=0.0.0.0            # 服务器监听地址
```

### 数据库设置

后端使用 MySQL 作为数据库。确保您已经创建了相应的数据库：

```sql
CREATE DATABASE i18n_flow CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'i18nflow'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON i18n_flow.* TO 'i18nflow'@'localhost';
FLUSH PRIVILEGES;
```

::: tip 提示
首次启动时，系统会自动创建所需的数据库表结构。
:::

## 启动服务

启动后端服务非常简单：

```bash
# 进入后端目录
cd admin-backend

# 下载依赖
go mod download

# 启动服务
go run main.go
```

服务默认将在 `http://localhost:8080` 上启动。您可以通过访问 `http://localhost:8080/swagger/index.html` 来查看 API 文档。

## API 认证

i18n Flow 后端使用两种认证方式：

1. **JWT 认证**：用于前端用户访问
2. **API 密钥认证**：用于 CLI 工具访问

### JWT 认证

获取 JWT 令牌：

```http
POST /api/login

{
  "username": "admin",
  "password": "admin123"
}
```

成功后返回：

```json
{
  "token": "eyJhbGciOiJI...",
  "refresh_token": "eyJhbGciOiJI...",
  "expires_at": "2023-03-01T12:00:00Z"
}
```

在后续请求中，将令牌添加到 Authorization 头：

```
Authorization: Bearer eyJhbGciOiJI...
```

### API 密钥认证

CLI 工具使用预设的 API 密钥进行认证，该密钥在 `.env` 文件中配置（`CLI_API_KEY`）。

在 CLI 请求中，需要添加以下头部：

```
X-API-Key: your_secure_api_key
```

## 主要 API 端点

### 项目管理

- `POST /api/projects` - 创建项目
- `GET /api/projects` - 获取项目列表
- `GET /api/projects/detail/:id` - 获取项目详情
- `PUT /api/projects/update/:id` - 更新项目
- `DELETE /api/projects/delete/:id` - 删除项目

### 语言管理

- `GET /api/languages` - 获取语言列表
- `POST /api/languages` - 创建语言
- `PUT /api/languages/:id` - 更新语言
- `DELETE /api/languages/:id` - 删除语言

### 翻译管理

- `POST /api/translations` - 创建翻译
- `POST /api/translations/batch` - 批量创建翻译
- `GET /api/translations/by-project/:project_id` - 获取项目的所有翻译
- `GET /api/translations/matrix/by-project/:project_id` - 获取项目的翻译矩阵
- `GET /api/translations/:id` - 获取翻译详情
- `PUT /api/translations/:id` - 更新翻译
- `DELETE /api/translations/:id` - 删除翻译
- `POST /api/translations/batch-delete` - 批量删除翻译

### 导入/导出

- `GET /api/exports/project/:project_id` - 导出项目翻译
- `POST /api/imports/project/:project_id` - 导入项目翻译

### CLI 工具集成

- `GET /api/cli/translations` - 获取翻译（CLI 使用）
- `POST /api/cli/keys` - 推送翻译键（CLI 使用）

### 仪表板

- `GET /api/dashboard/stats` - 获取系统统计信息

## 数据模型

i18n Flow 后端使用以下主要数据模型：

### 项目 (Project)

```go
type Project struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    Slug        string    `json:"slug" gorm:"uniqueIndex;not null"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 语言 (Language)

```go
type Language struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Code      string    `json:"code" gorm:"uniqueIndex;not null"`
    Locale    string    `json:"locale" gorm:"uniqueIndex"`
    IsRTL     bool      `json:"is_rtl" gorm:"default:false"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 翻译键 (Key)

```go
type Key struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    ProjectID   uint      `json:"project_id" gorm:"not null"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 翻译值 (Translation)

```go
type Translation struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    KeyID      uint      `json:"key_id" gorm:"not null"`
    LanguageID uint      `json:"language_id" gorm:"not null"`
    Value      string    `json:"value"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

## 日志和监控

i18n Flow 后端使用结构化日志记录系统活动。日志文件默认输出到标准输出，但您可以配置将其重定向到文件：

```bash
# 将日志输出到文件
go run main.go > i18n-flow.log 2>&1
```

## 备份与恢复

定期备份数据库是一个良好的实践：

```bash
# 备份数据库
mysqldump -u i18nflow -p i18n_flow > i18n_flow_backup.sql

# 恢复数据库
mysql -u i18nflow -p i18n_flow < i18n_flow_backup.sql
```

## 性能优化

对于大型项目，可以考虑以下优化：

1. 启用数据库连接池
2. 配置合适的缓存策略
3. 使用反向代理（如 Nginx）处理静态资源

这些配置可以在生产环境的 `.env` 文件中调整。

## 故障排除

### 数据库连接问题

如果遇到数据库连接问题：

1. 检查数据库凭证是否正确
2. 确认数据库服务是否运行
3. 检查网络连接和防火墙设置

### API 错误响应

系统返回的错误响应格式如下：

```json
{
  "error": "错误信息",
  "code": "错误代码",
  "details": {} // 可选的错误详情
}
```
