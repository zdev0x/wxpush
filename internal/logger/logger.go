package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/zdev0x/wxpush/internal/model"
)

var (
	logFile string
	mu      sync.Mutex
)

// Init 初始化日志
func Init(file string) error {
	// 只允许绝对路径，防止目录穿越
	if !filepath.IsAbs(file) {
		return fmt.Errorf("日志文件路径必须为绝对路径: %s", file)
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}
	_ = f.Close()
	logFile = file
	return nil
}

// write 写入日志
func write(entry model.LogEntry) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("序列化日志失败: %v", err)
	}

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("写入日志失败: %v", err)
	}

	return nil
}

// Info 写入info级别日志
func Info(action string, requestID string, extra map[string]interface{}) {
	entry := model.NewLogEntry(model.LevelInfo, action, model.StatusSuccess, requestID)
	if extra != nil {
		entry = entry.WithExtra(extra)
	}
	_ = write(entry)
}

// Error 写入error级别日志
func Error(action string, requestID string, code string, err error, extra map[string]interface{}) {
	entry := model.NewLogEntry(model.LevelError, action, model.StatusError, requestID).
		WithCode(code).
		WithMessage(err.Error())
	if extra != nil {
		entry = entry.WithExtra(extra)
	}
	_ = write(entry)
}

// Warn 写入warn级别日志
func Warn(action string, requestID string, message string, extra map[string]interface{}) {
	entry := model.NewLogEntry(model.LevelWarn, action, "", requestID).
		WithMessage(message)
	if extra != nil {
		entry = entry.WithExtra(extra)
	}
	_ = write(entry)
}

// TemplatePush 写入模板推送日志
func TemplatePush(status string, templateName string, groupName string, requestID string, err error, result *model.SendResult) {
	extra := map[string]interface{}{
		"template": templateName,
		"group":    groupName,
	}
	if result != nil {
		extra["success_count"] = result.SuccessCount
		extra["failed_count"] = result.FailedCount
		extra["success_users"] = result.SuccessUsers
		extra["failed_users"] = result.FailedUsers
	}

	entry := model.NewLogEntry(model.LevelInfo, model.ActionPushMessage, status, requestID)
	if err != nil {
		entry = entry.WithCode(model.ErrSendFailed).WithMessage(err.Error())
	}
	entry = entry.WithExtra(extra)
	_ = write(entry)
}
