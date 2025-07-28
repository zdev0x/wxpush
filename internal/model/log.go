// Package model defines data structures and models for the WeChat push service.
package model

import "time"

// LogEntry 日志条目
type LogEntry struct {
	Time      string                 `json:"time"`                 // 时间，RFC3339格式
	Level     string                 `json:"level"`                // 日志级别：info/error/warn
	Action    string                 `json:"action"`               // 操作类型
	Status    string                 `json:"status,omitempty"`     // 状态：success/error
	Code      string                 `json:"code,omitempty"`       // 错误代码
	Message   string                 `json:"message,omitempty"`    // 错误消息
	RequestID string                 `json:"request_id,omitempty"` // 请求ID
	Extra     map[string]interface{} `json:"extra,omitempty"`      // 附加信息
}

// 日志级别常量
const (
	LevelInfo  = "info"
	LevelError = "error"
	LevelWarn  = "warn"
)

// 操作类型常量
const (
	// 服务操作
	ActionServerStart  = "server_start"  // 服务启动
	ActionServerStop   = "server_stop"   // 服务停止
	ActionTokenRefresh = "token_refresh" // 刷新Token

	// 微信操作
	ActionWeChatVerify = "wechat_verify" // 微信验证
	ActionPushMessage  = "push_message"  // 消息推送
)

// NewLogEntry 创建日志条目
func NewLogEntry(level, action, status string, requestID string) LogEntry {
	return LogEntry{
		Time:      time.Now().Format(time.RFC3339),
		Level:     level,
		Action:    action,
		Status:    status,
		RequestID: requestID,
	}
}

// WithCode 设置错误代码
func (e LogEntry) WithCode(code string) LogEntry {
	e.Code = code
	return e
}

// WithMessage 设置错误消息
func (e LogEntry) WithMessage(message string) LogEntry {
	e.Message = message
	return e
}

// WithExtra 设置附加信息
func (e LogEntry) WithExtra(extra map[string]interface{}) LogEntry {
	e.Extra = extra
	return e
}
