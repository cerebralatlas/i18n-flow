# 安装说明

本指南将帮助您安装和设置 i18n Flow 的所有组件，包括后端服务器、前端管理界面和 CLI 工具。

## 系统要求

在开始安装之前，请确保您的系统满足以下要求：

- **操作系统**：Linux、macOS 或 Windows
- **Go**：1.20 或更高版本
- **Node.js**：18 或更高版本
- **MySQL**：8.0 或更高版本
- **包管理器**：pnpm（推荐）或 npm
- **Docker**（可选，用于容器化部署）：20.10+ 和 Docker Compose v2+

## 基础安装

### 1. 克隆仓库

首先，从 GitHub 克隆 i18n Flow 仓库：

```bash
git clone https://github.com/ilukemagic/i18n-flow.git
cd i18n-flow
```

### 2. 后端安装

进入后端目录并设置环境：

```bash
cd admin-backend

# 创建环境配置文件
cp .env.example .env

# 编辑 .env 文件，配置数据库连接和其他设置
# 根据您的环境修改以下内容：
# - 数据库凭证
# - JWT 密钥
# - API 密钥
# - 其他配置选项
```

安装依赖并启动后端服务：

```bash
go mod download
go run main.go
```

后端 API 将在 `http://localhost:8080` 可用，Swagger 文档在 `http://localhost:8080/swagger/index.html`。

### 3. 前端安装

进入前端目录并安装依赖：

```bash
cd ../admin-frontend
pnpm install
```

启动开发服务器：

```bash
pnpm dev
```

管理界面将在 `http://localhost:5173` 可用。

### 4. CLI 工具安装

CLI 工具可以全局安装，以便在任何项目中使用：

```bash
npm install -g i18n-flow-cli
```

## Docker 部署

i18n Flow 可以使用 Docker 和 Docker Compose 快速部署：

```bash
# 在项目根目录中
cp .env.example .env

# 编辑 .env 文件，配置必要的环境变量

# 启动服务
docker compose up -d
```

这将启动三个服务：

- MySQL 数据库（端口 3306）
- 后端 API（端口 8080）
- 前端管理面板（端口 80）

### 访问应用

- 管理界面：http://localhost
- API 和 Swagger 文档：http://localhost/swagger/index.html

### 初始登录

- 用户名：`admin`
- 密码：`admin123`
- **重要提示**：首次登录后请修改默认密码！

## 开发环境设置

如果您想参与开发，请按照以下步骤设置开发环境：

### 后端开发

后端使用 Go 和 Gin 框架开发：

```bash
cd admin-backend

# 设置开发环境
go mod download

# 运行测试
go test ./...

# 启动开发服务器
go run main.go
```

### 前端开发

前端使用 React、TypeScript 和 Vite 开发：

```bash
cd admin-frontend

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build
```

### CLI 工具开发

CLI 工具使用 TypeScript/JavaScript 和 Node.js 开发：

```bash
cd i18n-flow-cli

# 安装依赖
pnpm install

# 构建工具
pnpm build

# 链接到全局
npm link
```

## 下一步

- [前端使用指南](/zh/guide/frontend-guide) - 学习如何使用管理界面
- [后端使用指南](/zh/guide/backend-guide) - 了解后端 API 和配置
- [CLI 工具指南](/zh/guide/cli-guide) - 学习如何使用命令行工具
- [Docker 部署](/zh/deployment/docker) - 详细的 Docker 部署说明
