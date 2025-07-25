package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zdev0x/wxpush/internal/config"
)

// Config 配置中间件
func Config(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	}
}
