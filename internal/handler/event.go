package handler

import (
	"net/http"

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
	cfg := c.MustGet("config").(*config.Config)

	if signature == "" || timestamp == "" || nonce == "" || echostr == "" {
		logger.Error(model.ActionWeChatVerify, requestID, model.ErrInvalidParam, nil, map[string]interface{}{
			"signature": signature,
			"timestamp": timestamp,
			"nonce":     nonce,
		})
		c.String(http.StatusBadRequest, "missing parameters")
		return
	}

	if wechat.CheckSignature(cfg, signature, timestamp, nonce) {
		logger.Info(model.ActionWeChatVerify, requestID, map[string]interface{}{
			"signature": signature,
			"timestamp": timestamp,
			"nonce":     nonce,
		})
		c.String(http.StatusOK, echostr)
	} else {
		logger.Error(model.ActionWeChatVerify, requestID, model.ErrUnauthorized, nil, map[string]interface{}{
			"signature": signature,
			"timestamp": timestamp,
			"nonce":     nonce,
		})
		c.String(http.StatusForbidden, "signature verification failed")
	}
}

// HandleWeChatEventPost 处理微信事件推送
func HandleWeChatEventPost(c *gin.Context) {
	// TODO: 处理微信事件推送
	c.String(http.StatusOK, "ok")
}
