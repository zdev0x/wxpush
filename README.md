<div align="center">

# WxPush

[![Docker](https://github.com/zdev0x/wxpush/actions/workflows/docker.yml/badge.svg)](https://github.com/zdev0x/wxpush/actions/workflows/docker.yml)
[![Release](https://github.com/zdev0x/wxpush/actions/workflows/release.yml/badge.svg)](https://github.com/zdev0x/wxpush/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/zdev0x/wxpush/graph/badge.svg?token=YOUR_CODECOV_TOKEN)](https://codecov.io/gh/zdev0x/wxpush)
[![GoDoc](https://godoc.org/github.com/zdev0x/wxpush?status.svg)](https://godoc.org/github.com/zdev0x/wxpush)
[![License](https://img.shields.io/github/license/zdev0x/wxpush.svg)](https://github.com/zdev0x/wxpush/blob/main/LICENSE)

**基于 Go 语言的微信公众号消息推送服务**

支持多模板、多用户组，适合自动化、监控、通知等场景

✨ **支持微信公众平台测试号，可自建个人微信通知系统** ✨

---

## 🚀 **推荐使用 Release 版本**

> **⚠️ 重要提示：** 强烈建议使用 [**最新 Release 版本**](https://github.com/zdev0x/wxpush/releases/latest) 而非主分支代码，以获得更稳定的使用体验和完整的功能支持。

[![Latest Release](https://img.shields.io/github/v/release/zdev0x/wxpush?style=for-the-badge&logo=github&color=brightgreen)](https://github.com/zdev0x/wxpush/releases/latest)

---

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

### 方式一：二进制文件 (强烈推荐)

> 📦 **推荐选择：** 下载预编译的二进制文件，无需安装Go环境，开箱即用

```bash
# 下载最新 Release 版本
wget https://github.com/zdev0x/wxpush/releases/latest/download/wxpush_Linux_x86_64.tar.gz
tar -xzf wxpush_Linux_x86_64.tar.gz

# 配置并运行
cp config.example.yaml config.yaml
vim config.yaml
./wxpush -c config.yaml
```

### 方式二：Docker 运行

```bash
# 拉取最新 Release 版本镜像 (推荐)
docker pull ghcr.io/zdev0x/wxpush:latest

# 或使用指定版本号
# docker pull ghcr.io/zdev0x/wxpush:v1.0.0

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
  ghcr.io/zdev0x/wxpush:latest
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
  - name: "sms_forward"      # 模板名称，用于API调用
    id: "TEMPLATE_ID"        # 微信模板消息 ID
    title: "短信转发"        # 模板说明
    content: "短信来源：{{SOURCE.DATA}} 消息内容：{{CONTENT.DATA}} 发送时间：{{DATETIME.DATA}}"
  
  - name: "github_monitor"   # GitHub监控模板
    id: "XYZ1234567890abc"   # 微信模板消息 ID
    title: "GitHub监控"      # 模板说明  
    content: "监控人员：{{MONITOR.DATA}} 操作时间：{{CREATED_AT.DATA}} 操作类型：{{EVENT_TYPE.DATA}} 项目名称：{{PROJECT_NAME.DATA}} 项目作者：{{DEVELOPER.DATA}} 告警时间：{{DATETIME.DATA}}"
```

**模板变量说明：**
- 模板 `content` 中使用 `{{VARIABLE_NAME.DATA}}` 格式定义变量占位符
- API 请求时，请求体中的字段名为 `VARIABLE_NAME`（不包含 `.DATA` 后缀）
- 系统会自动将请求体中的字段值替换到对应的模板占位符中

**如何获取微信模板 ID：**
1. 登录 [微信公众平台](https://mp.weixin.qq.com/)
2. 进入"功能" -> "模板消息"
3. 添加或选择现有模板，获取模板 ID
4. 模板格式必须与配置文件中的 `content` 字段匹配

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

根据模板变量动态构建请求体，例如 `sms_forward` 模板：

```bash
# 使用 sms_forward 模板发送消息
curl -X POST "http://localhost:8801/wx/push?api_key=your_api_key&template=sms_forward&notify_group=default" \
  -H "Content-Type: application/json" \
  -d '{
    "SOURCE": "移动10086",
    "CONTENT": "您的验证码是123456，5分钟内有效。",
    "DATETIME": "2025-07-28 15:30:45"
  }'
```

使用 `github_monitor` 模板：

```bash
curl -X POST "http://localhost:8801/wx/push?api_key=your_api_key&template=github_monitor&notify_group=admin" \
  -H "Content-Type: application/json" \
  -d '{
    "MONITOR": "GitHub Bot",
    "CREATED_AT": "2025-07-28T15:30:45Z",
    "EVENT_TYPE": "push",
    "PROJECT_NAME": "awesome-project",
    "DEVELOPER": "zdev0x",
    "DATETIME": "2025-07-28 15:30:45"
  }'
```

**请求体参数说明：**
- 请求体中的字段名必须与模板中定义的变量名完全一致（不包含 `.DATA` 后缀）
- 每个字段的值将替换模板内容中对应的 `{{FIELD_NAME.DATA}}` 占位符
- 所有在模板 `content` 中使用的变量都必须在请求体中提供

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
    image: ghcr.io/zdev0x/wxpush:latest  # 使用最新 Release 版本
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

使用一键管理脚本：

```bash
# 下载管理脚本
curl -fsSL https://raw.githubusercontent.com/zdev0x/wxpush/main/scripts/manage.sh -o manage.sh
chmod +x manage.sh

# 安装服务
sudo ./manage.sh install

# 编辑配置文件
sudo ./manage.sh config

# 启动服务
sudo ./manage.sh start

# 查看状态
./manage.sh status

# 查看日志
./manage.sh logs
```

**管理脚本支持的命令：**
- `install` - 安装服务
- `uninstall [--keep]` - 卸载服务（--keep 保留配置）
- `start/stop/restart` - 启动/停止/重启服务
- `enable/disable` - 启用/禁用开机启动
- `status` - 查看服务状态
- `logs` - 查看实时日志
- `update` - 更新到最新版本
- `config` - 编辑配置文件
- `help` - 显示帮助信息

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
<summary><strong>Q: 如何获取模板变量列表？</strong></summary>

A: 
1. 查看配置文件中的模板定义，如 `config.example.yaml`
2. 模板 `content` 中的 `{{VARIABLE_NAME.DATA}}` 对应请求参数 `VARIABLE_NAME`
3. 例如模板内容为 "消息：{{CONTENT.DATA}} 时间：{{TIME.DATA}}"，则请求体需要包含 `CONTENT` 和 `TIME` 字段
4. 可以通过微信公众平台查看模板消息的字段要求
</details>

<details>
<summary><strong>Q: 模板变量不匹配怎么办？</strong></summary>

A: 
1. 确保请求体中的字段名与模板变量名一致（去掉 `.DATA` 后缀）
2. 检查配置文件中模板的 `content` 字段格式
3. 所有模板中使用的变量都必须在请求体中提供
4. 变量名区分大小写
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