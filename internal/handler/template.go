package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zdev0x/wxpush/internal/config"
	"github.com/zdev0x/wxpush/internal/logger"
	"github.com/zdev0x/wxpush/internal/model"
	"github.com/zdev0x/wxpush/internal/service/wechat"
)

// HandleTemplateMsgPush 处理模板消息推送
func HandleTemplateMsgPush(c *gin.Context) {
	requestID := c.GetString("request_id")

	// 获取配置
	cfg := c.MustGet("config").(*config.Config)

	// 验证API密钥
	apiKey := c.Query("api_key")
	if apiKey != cfg.Server.APIKey {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInvalidAPIKey, nil, nil)
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(
			model.ErrInvalidAPIKey,
			"无效的API密钥",
			requestID,
		))
		return
	}

	// 获取必要参数
	templateName := c.Query("template")
	groupName := c.Query("notify_group")
	if templateName == "" || groupName == "" {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInvalidParam, nil, map[string]interface{}{
			"template": templateName,
			"group":    groupName,
		})
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.ErrInvalidParam,
			"缺少模板名称或通知组名称",
			requestID,
		))
		return
	}

	// 验证模板是否存在
	if _, err := config.GetTemplate(cfg, templateName); err != nil {
		logger.Error(model.ActionPushMessage, requestID, model.ErrTemplateNotFound, err, map[string]interface{}{
			"template": templateName,
		})
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.ErrTemplateNotFound,
			"模板不存在",
			requestID,
		))
		return
	}

	// 解析请求体
	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInvalidRequest, err, nil)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.ErrInvalidRequest,
			"请求体格式错误",
			requestID,
		))
		return
	}

	// 发送模板消息
	result, err := wechat.SendTemplateMsg(cfg, templateName, groupName, params, requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.ErrSendFailed,
			err.Error(),
			requestID,
		))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(
		"消息已发送",
		result,
		requestID,
	))
}
