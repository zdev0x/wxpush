# WxPush

[![CI](https://github.com/zdev0x/wxpush/actions/workflows/ci.yml/badge.svg)](https://github.com/zdev0x/wxpush/actions/workflows/ci.yml)
[![Docker](https://github.com/zdev0x/wxpush/actions/workflows/docker.yml/badge.svg)](https://github.com/zdev0x/wxpush/actions/workflows/docker.yml)
[![Release](https://github.com/zdev0x/wxpush/actions/workflows/release.yml/badge.svg)](https://github.com/zdev0x/wxpush/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/zdev0x/wxpush/branch/main/graph/badge.svg)](https://codecov.io/gh/zdev0x/wxpush)
[![Go Report Card](https://goreportcard.com/badge/github.com/zdev0x/wxpush)](https://goreportcard.com/report/github.com/zdev0x/wxpush)
[![GoDoc](https://godoc.org/github.com/zdev0x/wxpush?status.svg)](https://godoc.org/github.com/zdev0x/wxpush)
[![License](https://img.shields.io/github/license/zdev0x/wxpush.svg)](https://github.com/zdev0x/wxpush/blob/main/LICENSE)

> 基于 Go 语言的微信公众号消息推送服务，支持多模板、多用户组，适合自动化、监控、通知等场景。**支持微信公众平台测试号，可自建个人微信通知系统。**

---

## 特性

- 🚀 高性能并发推送
- 🎯 灵活模板与分组配置
- 🔒 API 密钥安全校验
- 📦 支持 Docker、二进制、系统服务部署
- 🛠️ 易于二次开发和扩展
- 📑 结构清晰，文档完善
- 🧪 **支持微信公众平台测试号，自建个人微信通知系统**

---

## 快速开始

### 1. 拉取代码或镜像

```bash
git clone https://github.com/zdev0x/wxpush.git
cd wxpush
```

或直接使用 Docker：

```bash
docker pull ghcr.io/zdev0x/wxpush:main
```

### 2. 配置

复制并编辑配置文件：

```bash
cp config.example.yaml config.yaml
vim config.yaml
```

### 3. 本地运行

```bash
go build -o wxpush ./cmd/server
./wxpush -c config.yaml
```

### 4. Docker 运行

```bash
docker run -d \
  -p 8801:8801 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -v $(pwd)/logs:/app/logs \
  ghcr.io/zdev0x/wxpush:main
```
> 建议挂载 config.yaml 和 logs 目录，确保配置和日志持久化。

### Docker Compose

```yaml
version: '3'
services:
  wxpush:
    image: ghcr.io/zdev0x/wxpush:main
    restart: unless-stopped
    ports:
      - "8801:8801"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./logs:/app/logs
    environment:
      - TZ=Asia/Shanghai
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8801/wx/event"]
      interval: 1m
      timeout: 10s
      retries: 3
      start_period: 30s
```

---

## 配置说明

配置文件采用 YAML 格式，主要包含：

- **wechat**：公众号 appid、secret、token
- **templates**：消息模板配置
- **users**：用户 openid 配置
- **notify_groups**：通知分组
- **server**：服务参数（API 密钥、监听端口、日志、模式）

详见 [`config.example.yaml`](config.example.yaml)。

---

## API 示例

发送模板消息：

```http
POST /wx/push?api_key=xxx&template=sms_forward&notify_group=default
Content-Type: application/json

{
  "source": "中国移动",
  "content": "您的验证码是：123456",
  "datetime": "2024-03-21 15:04:05"
}
```

返回：

```json
{
  "status": "success",
  "message": "消息已发送",
  "data": {
    "success_count": 2,
    "failed_count": 0,
    "success_users": ["user1", "user2"],
    "failed_users": []
  }
}
```

---

## 部署方式

### 一键安装（Linux）

```bash
curl -O https://raw.githubusercontent.com/zdev0x/wxpush/main/scripts/install.sh
chmod +x install.sh
sudo ./install.sh
```

### Docker Compose

```yaml
version: '3'
services:
  wxpush:
    image: ghcr.io/zdev0x/wxpush:main
    restart: unless-stopped
    ports:
      - "8801:8801"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./logs:/app/logs
    environment:
      - TZ=Asia/Shanghai
```

---

## 目录结构

```
wxpush/
├── cmd/            # 程序入口
├── internal/       # 业务与核心模块
├── scripts/        # 部署脚本
├── docs/           # 文档
├── config.yaml     # 配置文件
└── ...
```

---

## 开发与贡献

- Go 1.21+
- `go mod download`
- `go test ./...`
- 代码风格遵循 golangci-lint 检查

欢迎 PR 和 Issue！

---

## 常见问题

- **端口冲突**：请确保 8801 端口未被占用。
- **日志无输出**：请检查 logs 目录挂载和权限。
- **配置未生效**：确认 config.yaml 路径正确并已挂载到容器内 `/app/config.yaml`。

---

## 许可证

MIT License © [zdev0x](https://github.com/zdev0x)

---

如果本项目对你有帮助，欢迎 Star ⭐️