# i18n Flow 文档网站

简体中文 | [English](README.md)

这是 i18n Flow 项目的官方文档网站，基于 VitePress 构建，支持中英文双语。

## 项目结构

- `docs/`: 文档内容
  - `en/`: 英文文档
  - `zh/`: 中文文档
- `Dockerfile`: 用于构建文档网站的 Docker 镜像
- `nginx.conf`: Nginx 配置文件

## 开发环境

### 本地开发

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建静态网站
pnpm build

# 预览构建结果
pnpm preview
```

## Docker 部署

### 集成到主项目 docker-compose.yml

将以下服务配置添加到项目根目录的 `docker-compose.yml` 文件中：

```yaml
# 文档网站服务
docs:
  build:
    context: ./i18n-flow-docs
    dockerfile: Dockerfile
  container_name: i18n_flow_docs
  restart: unless-stopped
  ports:
    - "8000:80" # 文档网站端口映射到主机的 8000 端口
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

同时，在 `volumes` 部分添加：

```yaml
volumes:
  # 其他已有的卷
  docs_nginx_logs:
```

修改完成后，使用以下命令启动整个应用：

```bash
# 在项目根目录下执行
docker compose up -d
```

文档网站将在 `http://localhost:8000` 可访问。

## 访问文档网站

部署完成后，您可以通过以下方式访问文档网站：

- **网址**：`http://localhost:8000`

### 可用内容

- **英文文档**：`http://localhost:8000/en/`
- **中文文档**：`http://localhost:8000/zh/`

### 文档章节

- 开始使用：i18n Flow 简介
- 安装说明：安装和配置步骤
- 前端使用指南：管理界面使用说明
- 后端使用指南：API 和配置详情
- CLI 工具指南：命令行工具用法
- 使用教程：完整工作流示例

您可以使用网站右上角的语言选择器在中英文版本之间切换。

## 自定义配置

### 修改端口

如果需要修改端口，编辑主项目 `docker-compose.yml` 文件中文档服务的端口映射：

```yaml
ports:
  - "自定义端口:80" # 例如 "9000:80"
```

### 使用自定义域名

1. 修改 `nginx.conf` 文件中的 `server_name` 指令：

```nginx
server_name docs.your-domain.com;
```

2. 设置域名解析到服务器 IP

3. 配置 HTTPS（建议生产环境使用）：
   - 使用 Let's Encrypt 获取免费证书
   - 配置 Nginx SSL 配置
