package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zdev0x/wxpush/internal/config"
	"github.com/zdev0x/wxpush/internal/middleware"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	// 添加配置中间件
	r.Use(middleware.Config(cfg))

	// 注册路由
	r.GET("/wx/event", HandleWeChatEvent)
	r.POST("/wx/event", HandleWeChatEventPost)
	r.POST("/wx/push", HandleTemplateMsgPush)
}
