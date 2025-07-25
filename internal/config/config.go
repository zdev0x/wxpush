package config

// Config 配置结构体
type Config struct {
	WeChat       WeChatConfig  `yaml:"wechat"`
	Templates    []Template    `yaml:"templates"`
	Users        []User        `yaml:"users"`
	NotifyGroups []NotifyGroup `yaml:"notify_groups"`
	Server       ServerConfig  `yaml:"server"`
}

// WeChatConfig 微信配置
type WeChatConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
	Token     string `yaml:"token"`
}

// Template 模板配置
type Template struct {
	Name    string `yaml:"name"`    // 模板名称
	ID      string `yaml:"id"`      // 模板ID
	Title   string `yaml:"title"`   // 模板说明
	Content string `yaml:"content"` // 模板格式
}

// User 用户配置
type User struct {
	Name   string `yaml:"name"`
	OpenID string `yaml:"openid"`
}

// NotifyGroup 通知组配置
type NotifyGroup struct {
	Name  string   `yaml:"name"`
	Users []string `yaml:"users"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	APIKey     string `yaml:"api_key"`     // API访问密钥
	ListenAddr string `yaml:"listen_addr"` // 监听地址
	LogFile    string `yaml:"log_file"`    // 日志文件路径
	Mode       string `yaml:"mode"`        // 运行模式
}

// 运行模式常量
const (
	ModeRelease = "release" // 生产模式
	ModeDebug   = "debug"   // 调试模式
	ModeTest    = "test"    // 测试模式
)
