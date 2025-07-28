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

// HandleTemplateMsgPush 处理模板消息推送
func HandleTemplateMsgPush(c *gin.Context) {
	requestID := c.GetString("request_id")

	// 获取配置，增加容错处理
	cfgInterface, exists := c.Get("config")
	if !exists {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInternal, nil, nil)
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.ErrInternal,
			"服务器配置错误",
			requestID,
		))
		return
	}

	cfg, ok := cfgInterface.(*config.Config)
	if !ok {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInternal, nil, nil)
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.ErrInternal,
			"服务器配置错误",
			requestID,
		))
		return
	}

	// 验证API密钥，增加更详细的错误处理
	apiKey := c.Query("api_key")
	if apiKey == "" {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInvalidAPIKey, nil, map[string]interface{}{
			"error": "missing_api_key",
		})
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(
			model.ErrInvalidAPIKey,
			"缺少API密钥",
			requestID,
		))
		return
	}

	if apiKey != cfg.Server.APIKey {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInvalidAPIKey, nil, map[string]interface{}{
			"provided_key": apiKey,
		})
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(
			model.ErrInvalidAPIKey,
			"无效的API密钥",
			requestID,
		))
		return
	}

	// 获取必要参数，提供更详细的错误信息
	templateName := c.Query("template")
	groupName := c.Query("notify_group")
	
	var missingParams []string
	if templateName == "" {
		missingParams = append(missingParams, "template")
	}
	if groupName == "" {
		missingParams = append(missingParams, "notify_group")
	}
	
	if len(missingParams) > 0 {
		logger.Error(
			model.ActionPushMessage,
			requestID,
			model.ErrInvalidParam,
			nil,
			map[string]interface{}{
				"missing_params": missingParams,
				"template":       templateName,
				"group":          groupName,
			},
		)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.ErrInvalidParam,
			"缺少必需参数: "+strings.Join(missingParams, ", "),
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
			"模板不存在: "+templateName,
			requestID,
		))
		return
	}

	// 解析请求体，增加更好的错误处理
	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInvalidRequest, err, map[string]interface{}{
			"content_type": c.GetHeader("Content-Type"),
			"error":        err.Error(),
		})
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.ErrInvalidRequest,
			"请求体格式错误: "+err.Error(),
			requestID,
		))
		return
	}

	// 验证请求体不为空
	if params == nil {
		logger.Error(model.ActionPushMessage, requestID, model.ErrInvalidRequest, nil, map[string]interface{}{
			"error": "empty_request_body",
		})
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.ErrInvalidRequest,
			"请求体不能为空",
			requestID,
		))
		return
	}

	// 发送模板消息
	result, err := wechat.SendTemplateMsg(cfg, templateName, groupName, params, requestID)
	if err != nil {
		logger.Error(model.ActionPushMessage, requestID, model.ErrSendFailed, err, map[string]interface{}{
			"template": templateName,
			"group":    groupName,
		})
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.ErrSendFailed,
			"消息发送失败: "+err.Error(),
			requestID,
		))
		return
	}

	logger.Info(model.ActionPushMessage, requestID, map[string]interface{}{
		"template": templateName,
		"group":    groupName,
		"status":   "success",
	})

	c.JSON(http.StatusOK, model.NewSuccessResponse(
		"消息已发送",
		result,
		requestID,
	))
}
