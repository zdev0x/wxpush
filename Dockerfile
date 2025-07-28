# syntax=docker/dockerfile:1

# 第一阶段：构建Go可执行文件
FROM golang:1.23-alpine AS builder
WORKDIR /app

# 安装git和build基础依赖，确保go mod download可用
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod tidy && go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wxpush ./cmd/server

# 第二阶段：精简运行环境
FROM alpine:latest

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates wget tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 复制可执行文件并设置权限
COPY --from=builder /app/wxpush ./
RUN chmod +x /app/wxpush

# 创建日志目录并设置权限
RUN mkdir -p /app/logs && \
    chown -R appuser:appgroup /app && \
    chmod 755 /app/logs

# 切换到非 root 用户
USER appuser

# 设置日志目录
VOLUME ["/app/logs"]

EXPOSE 8801

HEALTHCHECK --interval=1m --timeout=10s --start-period=30s --retries=3 \
    CMD wget -q --spider http://localhost:8801/wx/event || exit 1

ENTRYPOINT ["/app/wxpush"]