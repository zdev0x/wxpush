package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/zdev0x/wxpush/internal/config"
	"github.com/zdev0x/wxpush/internal/logger"
	"github.com/zdev0x/wxpush/internal/model"
)

// SendTemplateMsg 发送模板消息
func SendTemplateMsg(
	cfg *config.Config,
	templateName string,
	groupName string,
	params map[string]interface{},
	requestID string,
) (*model.SendResult, error) {
	// 获取模板配置
	tmpl, err := config.GetTemplate(cfg, templateName)
	if err != nil {
		return nil, fmt.Errorf(
			"获取模板失败: %v",
			err,
		)
	}

	// 获取用户列表
	users, err := config.GetGroupUsers(cfg, groupName)
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %v", err)
	}

	// 获取access_token
	token, err := GetAccessToken(cfg)
	if err != nil {
		return nil, fmt.Errorf("获取access_token失败: %v", err)
	}

	// 准备模板数据
	data := make(map[string]map[string]string)
	for key, value := range params {
		data[key] = map[string]string{
			"value": fmt.Sprint(value),
		}
	}

	// 并发发送消息
	var wg sync.WaitGroup
	var mu sync.Mutex
	result := &model.SendResult{
		SuccessUsers: make([]string, 0, len(users)),
		FailedUsers:  make([]string, 0),
	}

	for _, openid := range users {
		wg.Add(1)
		go func(openid string) {
			defer wg.Done()

			msg := map[string]interface{}{
				"touser":      openid,
				"template_id": tmpl.ID,
				"data":        data,
			}

			jsonData, err := json.Marshal(msg)
			if err != nil {
				mu.Lock()
				result.FailedUsers = append(result.FailedUsers, openid)
				mu.Unlock()
				return
			}

			// url 来源于微信官方API，安全
			url := fmt.Sprintf(
				"https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s",
				token,
			)
			// #nosec G107
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				mu.Lock()
				result.FailedUsers = append(result.FailedUsers, openid)
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				mu.Lock()
				result.FailedUsers = append(result.FailedUsers, openid)
				mu.Unlock()
				return
			}

			var wxResp struct {
				ErrCode int    `json:"errcode"`
				ErrMsg  string `json:"errmsg"`
				MsgID   int64  `json:"msgid"`
			}

			if err := json.Unmarshal(body, &wxResp); err != nil {
				mu.Lock()
				result.FailedUsers = append(result.FailedUsers, openid)
				mu.Unlock()
				return
			}

			if wxResp.ErrCode != 0 {
				mu.Lock()
				result.FailedUsers = append(result.FailedUsers, openid)
				mu.Unlock()
				return
			}

			mu.Lock()
			result.SuccessUsers = append(result.SuccessUsers, openid)
			mu.Unlock()
		}(openid)
	}

	wg.Wait()

	// 更新统计
	result.SuccessCount = len(result.SuccessUsers)
	result.FailedCount = len(result.FailedUsers)

	// 记录发送结果
	if result.FailedCount > 0 {
		logger.Error(
			model.ActionPushMessage,
			requestID,
			model.ErrSendFailed,
			fmt.Errorf("部分用户发送失败"),
			map[string]interface{}{
				"template":      templateName,
				"group":         groupName,
				"success":       result.SuccessUsers,
				"failed":        result.FailedUsers,
				"success_count": result.SuccessCount,
				"failed_count":  result.FailedCount,
			},
		)
		return result, fmt.Errorf("部分用户发送失败: %v", result.FailedUsers)
	}

	logger.Info(model.ActionPushMessage, requestID, map[string]interface{}{
		"template":      templateName,
		"group":         groupName,
		"success":       result.SuccessUsers,
		"success_count": result.SuccessCount,
	})

	return result, nil
}
