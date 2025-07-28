// Package handler provides HTTP request handlers for the WeChat push service.
package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zdev0x/wxpush/internal/config"
	"github.com/zdev0x/wxpush/internal/logger"
	"github.com/zdev0x/wxpush/internal/model"
	"github.com/zdev0x/wxpush/internal/service/wechat"
)

// HandleWeChatEvent 处理微信服务器验证
func HandleWeChatEvent(c *gin.Context) {
	requestID := c.GetString("request_id")
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	// 获取配置
	cfg, exists := c.Get("config")
	if !exists {
		logger.Error(model.ActionWeChatVerify, requestID, model.ErrInternal, nil, nil)
		c.String(http.StatusInternalServerError, "server configuration error")
		return
	}
	
	config, ok := cfg.(*config.Config)
	if !ok {
		logger.Error(model.ActionWeChatVerify, requestID, model.ErrInternal, nil, nil)
		c.String(http.StatusInternalServerError, "server configuration error")
		return
	}

	// 详细的参数验证和错误信息
	var missingParams []string
	if signature == "" {
		missingParams = append(missingParams, "signature")
	}
	if timestamp == "" {
		missingParams = append(missingParams, "timestamp")
	}
	if nonce == "" {
		missingParams = append(missingParams, "nonce")
	}
	if echostr == "" {
		missingParams = append(missingParams, "echostr")
	}

	if len(missingParams) > 0 {
		logger.Error(model.ActionWeChatVerify, requestID, model.ErrInvalidParam, nil, map[string]interface{}{
			"missing_params": missingParams,
			"signature":      signature,
			"timestamp":      timestamp,
			"nonce":          nonce,
			"echostr":        echostr,
		})
		c.String(http.StatusBadRequest, "missing required parameters: "+strings.Join(missingParams, ", "))
		return
	}

	// 验证签名
	if wechat.CheckSignature(config, signature, timestamp, nonce) {
		logger.Info(model.ActionWeChatVerify, requestID, map[string]interface{}{
			"signature": signature,
			"timestamp": timestamp,
			"nonce":     nonce,
			"status":    "success",
		})
		c.String(http.StatusOK, echostr)
	} else {
		logger.Error(model.ActionWeChatVerify, requestID, model.ErrUnauthorized, nil, map[string]interface{}{
			"signature": signature,
			"timestamp": timestamp,
			"nonce":     nonce,
			"status":    "signature_failed",
		})
		c.String(http.StatusForbidden, "signature verification failed")
	}
}

// HandleWeChatEventPost 处理微信事件推送
func HandleWeChatEventPost(c *gin.Context) {
	// TODO: 处理微信事件推送
	c.String(http.StatusOK, "ok")
}
