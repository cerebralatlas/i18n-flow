# 构建阶段
FROM node:18-alpine as build

# 安装 pnpm
RUN npm install -g pnpm

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制源代码
COPY . .

# 构建静态站点
RUN pnpm run build

# 生产阶段
FROM nginx:alpine

# 复制构建的静态文件到 nginx 服务目录
COPY --from=build /app/docs/.vitepress/dist /usr/share/nginx/html

# 复制自定义 nginx 配置文件
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露 80 端口
EXPOSE 80

# 启动 nginx 服务
CMD ["nginx", "-g", "daemon off;"] 