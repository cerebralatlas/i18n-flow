# 后端架构

了解 Go 后端的技术实现和架构设计。

## 技术栈

- **Go 1.23+**
- **Gin** - HTTP Web 框架
- **GORM** - ORM 框架
- **MySQL 8.0** - 主数据库
- **Redis 7.0** - 缓存
- **Uber FX** - 依赖注入

## 分层架构

```
┌─────────────────────────────────────┐
│           Handlers (HTTP 层)        │
│    接收请求、参数验证、响应格式化     │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│           Services (业务层)          │
│         业务逻辑处理、编排            │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│         Repositories (数据层)        │
│         数据访问、缓存管理            │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│            Database & Cache          │
│            MySQL + Redis             │
└─────────────────────────────────────┘
```

## 目录结构

```
admin-backend/internal/
├── api/
│   ├── handlers/        # HTTP 处理器
│   ├── middleware/      # 中间件
│   ├── response/        # 响应格式化
│   └── routes/          # 路由定义
│
├── config/              # 配置加载
│   └── config.go
│
├── container/           # FX 依赖注入
│   └── container.go
│
├── domain/
│   ├── models.go        # 实体定义
│   └── repository.go    # 仓储接口
│
├── dto/                 # 数据传输对象
│   ├── requests/
│   └── responses/
│
├── repository/          # 数据访问实现
│   ├── base.go
│   ├── project.go
│   ├── translation.go
│   └── user.go
│
├── service/             # 业务逻辑
│   ├── auth.go
│   ├── project.go
│   └── translation.go
│
└── utils/               # 工具函数
    ├── jwt.go
    ├── password.go
    └── datetime.go
```

## 核心组件

### 依赖注入容器

使用 Uber FX 管理所有依赖：

```go
// container/container.go
func ProvideContainer() *fx.App {
    return fx.New(
        fx.Provide(
            config.NewConfig,
            repository.NewDB,
            repository.NewRedis,
            service.NewProjectService,
            service.NewTranslationService,
        ),
        fx.Invoke(func(svc *service.ProjectService) {
            // 初始化后执行
        }),
    )
}
```

### 中间件

| 中间件 | 功能 |
|--------|------|
| `CORS` | 跨域支持 |
| `JWT Auth` | JWT 认证 |
| `API Key Auth` | CLI 认证 |
| `Rate Limit` | API 限流 |
| `Logger` | 请求日志 |

### 缓存策略

- **项目列表**：缓存 5 分钟
- **翻译矩阵**：缓存 1 分钟
- **用户会话**：Redis 存储，过期时间 7 天

## 数据库模型

### 用户

```go
type User struct {
    ID        uint   `gorm:"primarykey"`
    Email     string `gorm:"uniqueIndex"`
    Password  string
    Role      string // admin, member, viewer
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 项目

```go
type Project struct {
    ID              uint   `gorm:"primarykey"`
    Name            string
    Description     string
    DefaultLanguage string
    ApiKey          string
    OwnerID         uint
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 翻译

```go
type Translation struct {
    ID        uint   `gorm:"primarykey"`
    ProjectID uint
    Key       string `gorm:"index"`
    Language  string `gorm:"index"`
    Value     string `gorm:"type:text"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## API 路由

```
/api
├── /login           # 登录
├── /refresh         # 刷新 Token
├── /register        # 注册
├── /projects        # 项目 CRUD
├── /languages       # 语言管理
├── /translations    # 翻译管理
├── /exports         # 导出
├── /imports         # 导入
├── /users           # 用户管理
├── /invitations     # 邀请管理
└── /cli             # CLI 专用
    ├── /auth
    ├── /translations
    └── /keys
```

## 下一步

- [前端架构 →](/architecture/frontend)
- [CLI 架构 →](/architecture/cli)
