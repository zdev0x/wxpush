package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load 加载配置文件
func Load(filename string) (*Config, error) {
	// 只允许绝对路径或当前目录，防止目录穿越
	if !filepath.IsAbs(filename) && !strings.HasPrefix(filename, ".") {
		return nil, fmt.Errorf("配置文件路径不安全: %s", filename)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置默认值
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = ModeRelease
	}

	// 验证运行模式
	if err := validateMode(cfg.Server.Mode); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validateMode 验证运行模式
func validateMode(mode string) error {
	switch mode {
	case ModeRelease, ModeDebug, ModeTest:
		return nil
	default:
		return fmt.Errorf("无效的运行模式 %q，可选值: release, debug, test", mode)
	}
}

// GetTemplate 根据名称获取模板
func GetTemplate(cfg *Config, name string) (*Template, error) {
	for _, t := range cfg.Templates {
		if t.Name == name {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("template not found: %s", name)
}

// GetGroupUsers 根据通知组名获取用户OpenID列表
func GetGroupUsers(cfg *Config, groupName string) ([]string, error) {
	// 查找通知组
	var group *NotifyGroup
	for _, g := range cfg.NotifyGroups {
		if g.Name == groupName {
			group = &g
			break
		}
	}
	if group == nil {
		return nil, fmt.Errorf("notify group not found: %s", groupName)
	}

	// 获取用户OpenID
	userMap := make(map[string]string)
	for _, u := range cfg.Users {
		userMap[u.Name] = u.OpenID
	}

	var openids []string
	var notFoundUsers []string

	for _, username := range group.Users {
		if openid, ok := userMap[username]; ok {
			openids = append(openids, openid)
		} else {
			notFoundUsers = append(notFoundUsers, username)
		}
	}

	if len(notFoundUsers) > 0 {
		return openids, fmt.Errorf("users not found: %v", notFoundUsers)
	}

	if len(openids) == 0 {
		return nil, fmt.Errorf("no valid users found in notify group: %s", groupName)
	}

	return openids, nil
}
