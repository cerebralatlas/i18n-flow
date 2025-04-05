FROM node:22-bookworm AS builder

# 安装 pnpm 和系统依赖
RUN npm install -g pnpm && \
    apt-get update && apt-get install -y python3 make g++

WORKDIR /app

# 复制依赖文件
COPY package.json pnpm-lock.yaml ./

# 安装依赖并修复权限
RUN pnpm install --frozen-lockfile && \
    chmod +x node_modules/typescript/bin/tsc

# 设置 PATH
ENV PATH="/app/node_modules/.bin:${PATH}"

# 复制源代码（排除 node_modules）
COPY . .

# 验证路径
RUN which tsc && ls -la node_modules/typescript/bin

# 构建
RUN pnpm run build

# 后续步骤保持不变
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]