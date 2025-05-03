# i18n-flow

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

简体中文 | [English](README.md)

i18n-flow 是一个全面的国际化(i18n)管理平台，专为开发团队、内容创作者和本地化专家设计，旨在简化翻译工作流程。

## 🌟 功能特点

- **项目管理**：通过项目组织翻译内容，实现更好的工作流程管理
- **多语言支持**：管理无限数量的语言翻译，支持自定义语言设置
- **基于键值的翻译系统**：通过键值对方式保持应用程序中翻译的一致性
- **批量操作**：支持批量导入、导出和更新翻译，节省时间
- **上下文支持**：添加上下文信息，提高翻译准确性
- **Excel 导入/导出**：支持标准格式，便于与现有工作流程集成
- **RESTful API**：提供完善文档的 API，方便与您的系统集成
- **用户认证**：安全的基于 JWT 的认证系统
- **响应式 UI**：使用 React 和 Ant Design 构建的现代管理仪表板
- **CLI 集成**：命令行工具，实现与开发工作流程的无缝集成

## 📦 组件构成

i18n-flow 由三个主要组件组成：

1. **管理后端**：基于 Go 的 API 服务器，管理所有 i18n 数据并提供 RESTful 接口
2. **管理前端**：基于 React 的仪表板，用于可视化管理翻译项目
3. **CLI 工具**：命令行界面，供开发人员将翻译与代码库同步

## 🚀 快速开始

### 前提条件

- Go 1.20 或更高版本
- Node.js 18 或更高版本
- MySQL 8.0
- pnpm(推荐)或 npm

### 后端设置

1. 克隆仓库：

   ```bash
   git clone https://github.com/ilukemagic/i18n-flow.git
   cd i18n-flow/admin-backend
   ```

2. 配置环境：

   ```bash
   cp .env.example .env
   # 编辑.env文件，填入您的数据库凭证和其他设置
   ```

3. 运行后端：

   ```bash
   go mod download
   go run main.go
   ```

4. API 将在`http://localhost:8080`可用，Swagger 文档在`http://localhost:8080/swagger/index.html`

### 前端设置

1. 导航到前端目录：

   ```bash
   cd ../admin-frontend
   ```

2. 安装依赖：

   ```bash
   pnpm install
   ```

3. 启动开发服务器：

   ```bash
   pnpm dev
   ```

4. 管理界面将在`http://localhost:5173`可用

### CLI 工具设置

1. 全局安装 CLI 工具：

   ```bash
   npm install -g i18n-flow-cli
   ```

2. 在您的项目中初始化 i18n-flow：

   ```bash
   i18n-flow init
   ```

3. 按照交互式设置配置您的项目。

## 🐳 Docker 部署

您可以使用 Docker 和 Docker Compose 轻松部署 i18n-flow。

### 前提条件

- Docker 20.10+
- Docker Compose v2+

### 快速部署

1. 克隆仓库：

   ```bash
   git clone https://github.com/ilukemagic/i18n-flow.git
   cd i18n-flow
   ```

2. 配置环境变量：

   ```bash
   cp .env.example .env
   ```

   编辑`.env`文件设置：

   - 数据库凭证
   - JWT 密钥
   - API 密钥
   - 其他配置选项

3. 启动服务：

   ```bash
   docker compose up -d
   ```

   这将启动三个服务：

   - **MySQL**数据库在端口 3306
   - **后端 API**在端口 8080
   - **前端管理面板**在端口 80

4. 访问应用程序：

   - 管理界面：<http://localhost>
   - API 和 Swagger 文档：<http://localhost/swagger/index.html>

5. 初始登录：
   - 用户名：`admin`
   - 密码：`admin123`
   - **重要**：首次登录后请修改默认密码！

### Docker Compose 配置

默认的`docker-compose.yml`包括：

- **数据库**：MySQL 8.0，通过 Docker 卷实现数据持久化
- **后端**：Go 1.23 API 服务器，连接 MySQL
- **前端**：通过 Nginx 提供的 React SPA，具有 API 代理功能

### 自定义配置

您可以通过以下方式调整部署：

1. 修改`.env`文件中的环境变量
2. 更改`docker-compose.yml`中的端口映射
3. 根据需要更新 Docker 构建上下文或卷

### 更新应用程序

要更新到新版本：

```bash
git pull
docker compose down
docker compose up -d --build
```

### 故障排除

- **数据库连接问题**：检查`.env`文件中的数据库凭证
- **前端无法加载**：验证 nginx 代理配置是否正确指向后端
- **后端无法启动**：使用`docker compose logs backend`检查后端日志

## 📚 文档

### 文档网站

i18n Flow 的文档以专门的网站形式提供：

- **网址**：通过 Docker 部署后，访问 `http://localhost:8000/zh/`
- **语言**：同时提供中文和英文版本
- **内容**：
  - 所有组件的综合指南
  - 安装和设置说明
  - 包含实际示例的使用教程
  - API 参考文档
  - 部署指南

访问特定语言版本：

- 中文：`http://localhost:8000/zh/`
- 英文：`http://localhost:8000/en/`

### API 文档

运行后端服务器时，API 文档通过 Swagger UI 提供：

- 在浏览器中打开`http://localhost:8080/swagger/index.html`

### 管理仪表板

管理仪表板提供直观的界面，用于：

- 管理项目及其翻译键
- 添加和更新语言定义
- 输入和编辑翻译
- 以各种格式导入和导出翻译
- 用户管理和访问控制
- 监控翻译状态和进度

### CLI 工具

CLI 工具提供以下命令：

- `init`：在您的项目中初始化 i18n-flow
- `sync`：将翻译从服务器同步到本地项目
- `push`：将翻译键推送到服务器
- `status`：检查项目的翻译状态

完整的 CLI 文档可通过运行`i18n-flow --help`获取。

## 🏗️ 架构

i18n-flow 采用现代技术栈构建，遵循清晰的架构方法：

### 后端

- **语言**：Go
- **框架**：Gin
- **数据库**：MySQL 与 GORM ORM
- **认证**：JWT 和 API 密钥
- **文档**：Swagger/OpenAPI

### 前端

- **语言**：TypeScript
- **框架**：React 19
- **UI 库**：Ant Design 5
- **状态管理**：React Context API
- **路由**：React Router 7
- **样式**：Tailwind CSS 4
- **构建工具**：Vite

### CLI 工具

- **语言**：TypeScript/JavaScript
- **运行时**：Node.js
- **分发**：npm 包

## 🔄 工作流程

i18n-flow 中的典型工作流程：

1. 通过管理仪表板创建项目
2. 定义您需要支持的语言
3. 使用 CLI 扫描源代码中的翻译键
4. 使用`i18n-flow push --scan`将新键推送到翻译服务器
5. 通过管理界面添加翻译
6. 使用`i18n-flow sync`将最新翻译同步到您的项目
7. 使用翻译函数在应用程序中集成翻译

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

1. Fork 仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m '添加一些很棒的功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

详细指南请参阅[贡献指南](CONTRIBUTING.md)。

## 📄 许可证

本项目采用 MIT 许可证 - 详情请参阅[LICENSE](LICENSE)文件。

## 📞 联系方式

如有问题或反馈，请提交 GitHub 问题或联系项目维护者。
