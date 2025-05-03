# 后端 API 参考

i18n Flow 提供了一套完整的 RESTful API，用于与系统进行交互。本文档详细介绍了可用的 API 端点、参数和响应格式。

## API 基础

### 基本 URL

所有 API 请求的基本 URL 为：

```
http://your-i18n-flow-server.com/api
```

### 认证

i18n Flow API 支持两种认证方式：

1. **JWT 认证**：用于前端管理界面
2. **API 密钥认证**：用于 CLI 工具和其他自动化集成

#### JWT 认证

对于需要用户身份验证的请求，需要在 HTTP 头部添加 JWT 令牌：

```
Authorization: Bearer <jwt-token>
```

#### API 密钥认证

对于 CLI 工具和自动化集成，需要在 HTTP 头部添加 API 密钥：

```
X-API-Key: <api-key>
```

API 密钥在后端 `.env` 文件中的 `CLI_API_KEY` 设置。

### 响应格式

所有 API 响应都使用 JSON 格式，遵循以下结构：

**成功响应**:

```json
{
  "success": true,
  "data": {
    // 响应数据
  }
}
```

**错误响应**:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "错误描述信息",
    "details": {} // 可选的额外错误信息
  }
}
```

### 分页

支持分页的端点使用以下查询参数：

- `page`: 页码，默认为 1
- `limit`: 每页项目数，默认为 20

分页响应包含以下元数据：

```json
{
  "success": true,
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 20,
    "pages": 5
  }
}
```

## 认证 API

### 登录

获取 JWT 令牌。

**请求**:

```
POST /login
```

**请求体**:

```json
{
  "username": "admin",
  "password": "password"
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2023-06-01T12:00:00Z",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

### 刷新令牌

使用刷新令牌获取新的 JWT 令牌。

**请求**:

```
POST /refresh
```

**请求体**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2023-06-01T12:00:00Z"
  }
}
```

### 验证 CLI API 密钥

验证 CLI API 密钥是否有效。

**请求**:

```
GET /cli/auth
```

**请求头**:

```
X-API-Key: your-api-key
```

**响应**:

```json
{
  "success": true,
  "data": {
    "valid": true
  }
}
```

## 项目 API

### 创建项目

创建新项目。

**请求**:

```
POST /projects
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**请求体**:

```json
{
  "name": "我的项目",
  "description": "项目描述",
  "slug": "my-project"
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "我的项目",
    "description": "项目描述",
    "slug": "my-project",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取项目列表

获取所有项目的列表，支持分页。

**请求**:

```
GET /projects?page=1&limit=20
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "我的项目",
      "description": "项目描述",
      "slug": "my-project",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "stats": {
        "languages": 3,
        "keys": 150,
        "completion": 75
      }
    }
    // 更多项目...
  ],
  "meta": {
    "total": 10,
    "page": 1,
    "limit": 20,
    "pages": 1
  }
}
```

### 获取项目详情

获取单个项目的详细信息。

**请求**:

```
GET /projects/detail/:id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "我的项目",
    "description": "项目描述",
    "slug": "my-project",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "stats": {
      "languages": 3,
      "keys": 150,
      "completion": 75,
      "language_stats": [
        {
          "language_code": "en",
          "language_name": "English",
          "total_keys": 150,
          "translated_keys": 150,
          "completion": 100
        },
        {
          "language_code": "zh-CN",
          "language_name": "简体中文",
          "total_keys": 150,
          "translated_keys": 125,
          "completion": 83
        },
        {
          "language_code": "ja",
          "language_name": "日本語",
          "total_keys": 150,
          "translated_keys": 75,
          "completion": 50
        }
      ]
    }
  }
}
```

### 更新项目

更新现有项目的信息。

**请求**:

```
PUT /projects/update/:id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**请求体**:

```json
{
  "name": "我的项目（更新）",
  "description": "更新后的项目描述",
  "slug": "my-project-updated"
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "我的项目（更新）",
    "description": "更新后的项目描述",
    "slug": "my-project-updated",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z"
  }
}
```

### 删除项目

删除项目及其所有相关数据。

**请求**:

```
DELETE /projects/delete/:id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": {
    "message": "项目已成功删除"
  }
}
```

## 语言 API

### 获取语言列表

获取所有支持的语言列表。

**请求**:

```
GET /languages
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "English",
      "code": "en",
      "locale": "en-US",
      "is_rtl": false,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "简体中文",
      "code": "zh-CN",
      "locale": "zh-Hans",
      "is_rtl": false,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
    // 更多语言...
  ]
}
```

### 创建语言

添加新的支持语言。

**请求**:

```
POST /languages
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**请求体**:

```json
{
  "name": "Español",
  "code": "es",
  "locale": "es-ES",
  "is_rtl": false
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "id": 3,
    "name": "Español",
    "code": "es",
    "locale": "es-ES",
    "is_rtl": false,
    "created_at": "2023-01-02T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z"
  }
}
```

### 更新语言

更新现有语言的信息。

**请求**:

```
PUT /languages/:id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**请求体**:

```json
{
  "name": "Español (España)",
  "locale": "es-ES",
  "is_rtl": false
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "id": 3,
    "name": "Español (España)",
    "code": "es",
    "locale": "es-ES",
    "is_rtl": false,
    "created_at": "2023-01-02T00:00:00Z",
    "updated_at": "2023-01-03T00:00:00Z"
  }
}
```

### 删除语言

删除语言及其所有相关翻译。

**请求**:

```
DELETE /languages/:id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": {
    "message": "语言已成功删除"
  }
}
```

## 翻译 API

### 创建翻译

添加新的翻译。

**请求**:

```
POST /translations
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**请求体**:

```json
{
  "project_id": 1,
  "key": "common.buttons.save",
  "description": "保存按钮文本",
  "translations": [
    {
      "language_id": 1,
      "value": "Save"
    },
    {
      "language_id": 2,
      "value": "保存"
    }
  ]
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "key_id": 1,
    "key": "common.buttons.save",
    "description": "保存按钮文本",
    "project_id": 1,
    "translations": [
      {
        "id": 1,
        "language_id": 1,
        "language_code": "en",
        "value": "Save"
      },
      {
        "id": 2,
        "language_id": 2,
        "language_code": "zh-CN",
        "value": "保存"
      }
    ],
    "created_at": "2023-01-05T00:00:00Z",
    "updated_at": "2023-01-05T00:00:00Z"
  }
}
```

### 批量创建翻译

批量添加多个翻译。

**请求**:

```
POST /translations/batch
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**请求体**:

```json
{
  "project_id": 1,
  "items": [
    {
      "key": "common.buttons.save",
      "description": "保存按钮文本",
      "translations": [
        {
          "language_id": 1,
          "value": "Save"
        },
        {
          "language_id": 2,
          "value": "保存"
        }
      ]
    },
    {
      "key": "common.buttons.cancel",
      "description": "取消按钮文本",
      "translations": [
        {
          "language_id": 1,
          "value": "Cancel"
        },
        {
          "language_id": 2,
          "value": "取消"
        }
      ]
    }
  ]
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "created": 2,
    "updated": 0,
    "failed": 0,
    "items": [
      {
        "key_id": 1,
        "key": "common.buttons.save",
        "success": true
      },
      {
        "key_id": 2,
        "key": "common.buttons.cancel",
        "success": true
      }
    ]
  }
}
```

### 按项目获取翻译

获取项目的所有翻译。

**请求**:

```
GET /translations/by-project/:project_id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": [
    {
      "key_id": 1,
      "key": "common.buttons.save",
      "description": "保存按钮文本",
      "project_id": 1,
      "translations": [
        {
          "id": 1,
          "language_id": 1,
          "language_code": "en",
          "value": "Save"
        },
        {
          "id": 2,
          "language_id": 2,
          "language_code": "zh-CN",
          "value": "保存"
        }
      ],
      "created_at": "2023-01-05T00:00:00Z",
      "updated_at": "2023-01-05T00:00:00Z"
    }
    // 更多翻译...
  ]
}
```

### 获取翻译矩阵

获取项目的翻译矩阵（表格视图）。

**请求**:

```
GET /translations/matrix/by-project/:project_id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": {
    "languages": [
      {
        "id": 1,
        "code": "en",
        "name": "English"
      },
      {
        "id": 2,
        "code": "zh-CN",
        "name": "简体中文"
      }
    ],
    "keys": [
      {
        "id": 1,
        "key": "common.buttons.save",
        "description": "保存按钮文本",
        "values": {
          "1": "Save",
          "2": "保存"
        }
      },
      {
        "id": 2,
        "key": "common.buttons.cancel",
        "description": "取消按钮文本",
        "values": {
          "1": "Cancel",
          "2": "取消"
        }
      }
    ]
  }
}
```

### 更新翻译

更新翻译值。

**请求**:

```
PUT /translations/:id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**请求体**:

```json
{
  "value": "保存更改"
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "id": 2,
    "language_id": 2,
    "language_code": "zh-CN",
    "key_id": 1,
    "value": "保存更改",
    "updated_at": "2023-01-06T00:00:00Z"
  }
}
```

### 删除翻译键

删除翻译键及其所有相关翻译值。

**请求**:

```
DELETE /translations/:key_id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": {
    "message": "翻译键已成功删除"
  }
}
```

## CLI API

### 获取翻译（CLI）

获取项目的翻译，供 CLI 工具使用。

**请求**:

```
GET /cli/translations?project_id=my-project&locales=en,zh-CN
```

**请求头**:

```
X-API-Key: your-api-key
```

**响应**:

```json
{
  "success": true,
  "data": {
    "project": {
      "id": 1,
      "name": "我的项目",
      "slug": "my-project"
    },
    "languages": [
      {
        "id": 1,
        "code": "en",
        "name": "English"
      },
      {
        "id": 2,
        "code": "zh-CN",
        "name": "简体中文"
      }
    ],
    "translations": {
      "en": {
        "common.buttons.save": "Save",
        "common.buttons.cancel": "Cancel"
      },
      "zh-CN": {
        "common.buttons.save": "保存",
        "common.buttons.cancel": "取消"
      }
    }
  }
}
```

### 推送翻译键（CLI）

从 CLI 工具推送新的翻译键。

**请求**:

```
POST /cli/keys
```

**请求头**:

```
X-API-Key: your-api-key
```

**请求体**:

```json
{
  "project_id": "my-project",
  "keys": [
    "common.buttons.save",
    "common.buttons.cancel",
    "common.buttons.edit"
  ]
}
```

**响应**:

```json
{
  "success": true,
  "data": {
    "added": 1,
    "existing": 2,
    "failed": 0,
    "keys": [
      {
        "key": "common.buttons.save",
        "status": "existing"
      },
      {
        "key": "common.buttons.cancel",
        "status": "existing"
      },
      {
        "key": "common.buttons.edit",
        "status": "added"
      }
    ]
  }
}
```

## 导入/导出 API

### 导出项目翻译

导出项目的翻译。

**请求**:

```
GET /exports/project/:project_id?format=json&locales=en,zh-CN
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

文件下载，格式取决于 `format` 参数（json、xlsx、csv）。

### 导入项目翻译

导入项目的翻译。

**请求**:

```
POST /imports/project/:project_id
```

**请求头**:

```
Authorization: Bearer <jwt-token>
Content-Type: multipart/form-data
```

**请求参数**:

- `file`: 要导入的文件
- `format`: 文件格式 (json、xlsx、csv)
- `locale`: 语言代码 (仅用于单语言导入)
- `merge_strategy`: 合并策略 (overwrite、keep_existing)

**响应**:

```json
{
  "success": true,
  "data": {
    "imported": 50,
    "updated": 25,
    "skipped": 5,
    "project_id": 1,
    "locale": "zh-CN"
  }
}
```

## 仪表板 API

### 获取系统统计信息

获取系统统计信息，用于仪表板显示。

**请求**:

```
GET /dashboard/stats
```

**请求头**:

```
Authorization: Bearer <jwt-token>
```

**响应**:

```json
{
  "success": true,
  "data": {
    "projects": {
      "total": 5,
      "recent": [
        {
          "id": 1,
          "name": "我的项目",
          "slug": "my-project",
          "updated_at": "2023-01-06T00:00:00Z"
        }
        // 更多项目...
      ]
    },
    "languages": {
      "total": 6,
      "list": [
        {
          "code": "en",
          "name": "English",
          "projects_count": 5
        },
        {
          "code": "zh-CN",
          "name": "简体中文",
          "projects_count": 4
        }
        // 更多语言...
      ]
    },
    "translations": {
      "total_keys": 500,
      "total_translations": 2500,
      "completion_rate": 83.5
    },
    "recent_activity": [
      {
        "type": "translation_update",
        "project_id": 1,
        "project_name": "我的项目",
        "user": "admin",
        "timestamp": "2023-01-06T12:30:00Z",
        "details": {
          "key": "common.buttons.save",
          "language": "zh-CN"
        }
      }
      // 更多活动...
    ]
  }
}
```

## 错误代码

以下是系统可能返回的常见错误代码：

| 代码                            | 描述             |
| ------------------------------- | ---------------- |
| `AUTH_INVALID_CREDENTIALS`      | 无效的认证凭据   |
| `AUTH_TOKEN_EXPIRED`            | 认证令牌已过期   |
| `AUTH_TOKEN_INVALID`            | 无效的认证令牌   |
| `AUTH_INSUFFICIENT_PERMISSIONS` | 权限不足         |
| `RESOURCE_NOT_FOUND`            | 请求的资源不存在 |
| `RESOURCE_ALREADY_EXISTS`       | 资源已存在       |
| `VALIDATION_ERROR`              | 请求数据验证错误 |
| `SERVER_ERROR`                  | 服务器内部错误   |

## 速率限制

API 实施速率限制以防止滥用。限制如下：

- 认证端点：每 IP 地址每分钟 10 次请求
- 其他端点：每 API 密钥或用户每分钟 100 次请求

超过限制会返回 HTTP 429 状态码，并附带以下响应：

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "已超过速率限制，请稍后再试",
    "details": {
      "retry_after": 30
    }
  }
}
```
