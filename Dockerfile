# syntax=docker/dockerfile:1

# 第一阶段：构建Go可执行文件
FROM golang:1.21-alpine AS builder
WORKDIR /app

# 安装git，解决 go mod download 依赖私有/远程仓库时失败
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wxpush ./cmd/server

# 第二阶段：精简运行环境
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/wxpush ./
RUN mkdir -p /app/logs

# 设置时区
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 设置日志目录
VOLUME ["/app/logs"]

EXPOSE 8801

HEALTHCHECK --interval=1m --timeout=10s --start-period=30s --retries=3 \
    CMD wget -q --spider http://localhost:8801/wx/event || exit 1

ENTRYPOINT ["/app/wxpush"]