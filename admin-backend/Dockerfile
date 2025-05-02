FROM golang:1.22-alpine AS builder

WORKDIR /app

# 设置GOTOOLCHAIN以允许使用较低版本的依赖
ENV GOTOOLCHAIN=auto

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o i18n-flow-backend .

# 创建最终的轻量级镜像
FROM alpine:latest

WORKDIR /app

# 安装必要的CA证书 (可能需要HTTPS请求)
RUN apk --no-cache add ca-certificates

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/i18n-flow-backend .

# 创建.env文件存放位置
RUN mkdir -p /app/config
VOLUME /app/config

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./i18n-flow-backend"]