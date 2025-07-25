package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zdev0x/wxpush/internal/banner"
	"github.com/zdev0x/wxpush/internal/config"
	"github.com/zdev0x/wxpush/internal/handler"
	"github.com/zdev0x/wxpush/internal/logger"
	"github.com/zdev0x/wxpush/internal/middleware"
	"github.com/zdev0x/wxpush/internal/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// 版本信息，在编译时通过 -ldflags 设置
	version = "dev"     // -X main.version
	commit  = "none"    // -X main.commit
	date    = "unknown" // -X main.date

	// 帮助信息
	helpText = `Usage: wxpush [options]

Options:
  -c string
        配置文件路径 (default "config.yaml")
  -v    显示版本信息
`
)

func main() {
	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, helpText)
	}

	// 解析命令行参数
	configFile := flag.String("c", "config.yaml", "配置文件路径")
	showVersion := flag.Bool("v", false, "显示版本信息")
	flag.Parse()

	// 显示版本信息
	if *showVersion {
		if commit == "none" {
			fmt.Printf("wxpush version %s\n", version)
		} else {
			fmt.Printf("wxpush version %s (%s)\n", version, strings.ToLower(commit[:7]))
		}
		if date != "unknown" {
			fmt.Printf("build time: %s\n", date)
		}
		os.Exit(0)
	}

	// 确保配置文件存在
	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在: %s", *configFile)
	}

	// 加载配置文件
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 确保日志目录存在
	logDir := filepath.Dir(cfg.Server.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Server.LogFile); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 设置gin模式
	switch cfg.Server.Mode {
	case config.ModeDebug:
		gin.SetMode(gin.DebugMode)
	case config.ModeTest:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建gin引擎
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())

	// 注册路由
	handler.RegisterRoutes(r, cfg)

	// 显示启动banner
	banner.ShowBanner(version, commit, date, map[string]string{
		// 版本信息
		"Version":    version,
		"Commit":     commit,
		"Build Time": date,

		// 服务配置
		"Run Mode":       cases.Title(language.Und).String(cfg.Server.Mode),
		"Listen Address": cfg.Server.ListenAddr,
		"Log File":       cfg.Server.LogFile,

		// 微信配置
		"WeChat App ID":     cfg.WeChat.AppID,
		"Message Templates": fmt.Sprintf("%d template(s)", len(cfg.Templates)),
		"WeChat Users":      fmt.Sprintf("%d user(s)", len(cfg.Users)),
		"Notify Groups":     fmt.Sprintf("%d group(s)", len(cfg.NotifyGroups)),
	})

	// 启动服务
	logger.Info(model.ActionServerStart, "", map[string]interface{}{
		"version":     version,
		"commit":      commit,
		"build_time":  date,
		"listen_addr": cfg.Server.ListenAddr,
		"config_file": *configFile,
		"mode":        cfg.Server.Mode,
	})

	if err := r.Run(cfg.Server.ListenAddr); err != nil {
		logger.Error(model.ActionServerStart, "", model.ErrInternal, err, map[string]interface{}{
			"listen_addr": cfg.Server.ListenAddr,
		})
		os.Exit(1)
	}
}
