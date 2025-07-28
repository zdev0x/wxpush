# WxPush

<div align="center">

[![CI](https://github.com/zdev0x/wxpush/actions/workflows/ci.yml/badge.svg)](https://github.com/zdev0x/wxpush/actions/workflows/ci.yml)
[![Docker](https://github.com/zdev0x/wxpush/actions/workflows/docker.yml/badge.svg)](https://github.com/zdev0x/wxpush/actions/workflows/docker.yml)
[![Release](https://github.com/zdev0x/wxpush/actions/workflows/release.yml/badge.svg)](https://github.com/zdev0x/wxpush/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/zdev0x/wxpush/graph/badge.svg?token=YOUR_CODECOV_TOKEN)](https://codecov.io/gh/zdev0x/wxpush)
[![Go Report Card](https://goreportcard.com/badge/github.com/zdev0x/wxpush)](https://goreportcard.com/report/github.com/zdev0x/wxpush)
[![GoDoc](https://godoc.org/github.com/zdev0x/wxpush?status.svg)](https://godoc.org/github.com/zdev0x/wxpush)
[![License](https://img.shields.io/github/license/zdev0x/wxpush.svg)](https://github.com/zdev0x/wxpush/blob/main/LICENSE)

**基于 Go 语言的微信公众号消息推送服务**

支持多模板、多用户组，适合自动化、监控、通知等场景

✨ **支持微信公众平台测试号，可自建个人微信通知系统** ✨

</div>

## 📋 Table of Contents

- [特性](#-特性)
- [快速开始](#-快速开始)
- [配置说明](#-配置说明)
- [API 文档](#-api-文档)
- [部署方式](#-部署方式)
- [开发指南](#-开发指南)
- [常见问题](#-常见问题)
- [贡献指南](#-贡献指南)
- [许可证](#-许可证)

## ✨ 特性

- 🚀 **高性能并发推送** - 基于 Go 协程的高并发消息推送
- 🎯 **灵活模板与分组配置** - 支持多模板、多用户组管理
- 🔒 **API 密钥安全校验** - 确保接口调用安全性
- 📦 **多种部署方式** - 支持 Docker、二进制、系统服务部署
- 🛠️ **易于二次开发和扩展** - 清晰的代码结构，完善的文档
- 🧪 **微信测试号支持** - 支持微信公众平台测试号，快速搭建个人通知系统
- 📊 **完善的监控和日志** - 结构化日志输出，便于监控和调试
- 🌐 **跨平台支持** - 支持 Linux、Windows、macOS 多平台

## 🚀 快速开始

### 方式一：Docker 运行 (推荐)

```bash
# 拉取镜像
docker pull ghcr.io/zdev0x/wxpush:main

# 下载配置文件模板
curl -O https://raw.githubusercontent.com/zdev0x/wxpush/main/config.example.yaml
cp config.example.yaml config.yaml

# 编辑配置文件
vim config.yaml

# 运行容器
docker run -d \
  --name wxpush \
  -p 8801:8801 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -v $(pwd)/logs:/app/logs \
  --restart unless-stopped \
  ghcr.io/zdev0x/wxpush:main
```

### 方式二：二进制文件

```bash
# 下载最新版本
wget https://github.com/zdev0x/wxpush/releases/latest/download/wxpush_Linux_x86_64.tar.gz
tar -xzf wxpush_Linux_x86_64.tar.gz

# 配置并运行
cp config.example.yaml config.yaml
vim config.yaml
./wxpush -c config.yaml
```

### 方式三：源码编译

```bash
# 克隆代码
git clone https://github.com/zdev0x/wxpush.git
cd wxpush

# 安装依赖
go mod download

# 编译运行
go build -o wxpush ./cmd/server
./wxpush -c config.yaml
```

## ⚙️ 配置说明

配置文件采用 YAML 格式，主要包含以下部分：

### 微信公众号配置
```yaml
wechat:
  app_id: "YOUR_APPID"        # 微信公众平台 AppID
  app_secret: "YOUR_SECRET"   # 微信公众平台 AppSecret  
  token: "YOUR_TOKEN"         # 用于服务器验证的 Token
```

### 模板配置
```yaml
templates:
  - name: "notification"      # 模板名称
    id: "TEMPLATE_ID"         # 微信模板消息 ID
    title: "通知消息"         # 模板说明
    content: "内容：{{CONTENT.DATA}} 时间：{{TIME.DATA}}"
```

### 用户和分组配置
```yaml
users:
  - name: "user1"            # 用户别名
    openid: "USER_OPENID"    # 微信用户 OpenID

notify_groups:
  - name: "default"          # 分组名称
    users: ["user1"]         # 分组用户列表
```

### 服务器配置
```yaml
server:
  api_key: "YOUR_API_KEY"    # API 访问密钥
  listen_addr: ":8801"       # 监听地址
  log_file: "logs/push.log"  # 日志文件路径
  mode: "release"            # 运行模式：release/debug/test
```

详细配置示例请查看 [`config.example.yaml`](config.example.yaml)

## 📡 API 文档

### 发送模板消息

**POST** `/wx/push`

**参数：**
- `api_key` (query, required): API密钥
- `template` (query, required): 模板名称
- `notify_group` (query, required): 通知组名称

**请求体示例：**
```json
{
  "source": "系统通知",
  "content": "您有新的消息需要处理",
  "datetime": "2024-03-21 15:04:05"
}
```

**响应示例：**
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

### 服务器验证

**GET** `/wx/event`

用于微信服务器验证，无需手动调用。

## 🐳 部署方式

### Docker Compose (推荐)

创建 `docker-compose.yml`：

```yaml
version: '3.8'
services:
  wxpush:
    image: ghcr.io/zdev0x/wxpush:main
    container_name: wxpush
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

```bash
docker-compose up -d
```

### 系统服务 (Linux)

使用一键安装脚本：

```bash
curl -fsSL https://raw.githubusercontent.com/zdev0x/wxpush/main/scripts/install.sh | sudo bash
```

或手动安装：

```bash
# 下载二进制文件
sudo wget -O /usr/local/bin/wxpush https://github.com/zdev0x/wxpush/releases/latest/download/wxpush_Linux_x86_64.tar.gz

# 创建配置目录
sudo mkdir -p /etc/wxpush
sudo cp config.yaml /etc/wxpush/

# 创建 systemd 服务
sudo tee /etc/systemd/system/wxpush.service > /dev/null <<EOF
[Unit]
Description=WxPush Service
After=network.target

[Service]
Type=simple
User=wxpush
Group=wxpush
ExecStart=/usr/local/bin/wxpush -c /etc/wxpush/config.yaml
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable wxpush
sudo systemctl start wxpush
```

## 🛠️ 开发指南

### 环境要求

- Go 1.21 或更高版本
- Git

### 本地开发

```bash
# 克隆项目
git clone https://github.com/zdev0x/wxpush.git
cd wxpush

# 安装依赖
go mod download

# 运行测试
go test ./...

# 代码检查
golangci-lint run

# 构建
go build -o wxpush ./cmd/server
```

### 项目结构

```
wxpush/
├── cmd/                    # 程序入口点
│   └── server/
│       └── main.go
├── internal/               # 内部包
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP 处理器
│   ├── logger/            # 日志模块
│   ├── model/             # 数据模型
│   └── service/           # 业务逻辑
├── scripts/               # 部署脚本
├── docs/                  # 文档
├── .github/               # GitHub Actions
├── config.example.yaml    # 配置文件模板
└── README.md
```

### 贡献流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## ❓ 常见问题

<details>
<summary><strong>Q: 如何获取微信公众号的 AppID 和 AppSecret？</strong></summary>

A: 
1. 登录 [微信公众平台](https://mp.weixin.qq.com/)
2. 进入"开发" -> "基本配置"
3. 获取 AppID 和 AppSecret
4. 也可以使用[微信测试号](https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Requesting_an_API_Test_Account.html)进行测试
</details>

<details>
<summary><strong>Q: 如何获取用户的 OpenID？</strong></summary>

A: 
1. 用户关注公众号后，可通过微信开发者工具获取
2. 或通过接口获取关注者列表
3. 测试号环境下，关注后会显示 OpenID
</details>

<details>
<summary><strong>Q: 容器启动失败怎么办？</strong></summary>

A: 
1. 检查配置文件路径和权限：`docker logs wxpush`
2. 确保端口 8801 未被占用：`netstat -tlnp | grep 8801`
3. 检查配置文件格式是否正确
</details>

<details>
<summary><strong>Q: API 调用返回认证失败？</strong></summary>

A: 
1. 检查 `api_key` 参数是否正确
2. 确认配置文件中的 `api_key` 已正确设置
3. 注意 API 密钥的大小写敏感
</details>

## 🤝 贡献指南

我们欢迎各种形式的贡献！

### 报告 Bug
- 使用 [GitHub Issues](https://github.com/zdev0x/wxpush/issues) 报告 bug
- 请提供详细的错误信息和重现步骤

### 功能建议
- 使用 [GitHub Issues](https://github.com/zdev0x/wxpush/issues) 提出功能建议
- 详细描述建议的功能和使用场景

### 代码贡献
- 遵循现有的代码风格
- 确保所有测试通过
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证 - 详情请查看 [LICENSE](LICENSE) 文件。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

<div align="center">

**如果这个项目对你有帮助，请给它一个 ⭐️**

Made with ❤️ by [zdev0x](https://github.com/zdev0x)

</div>