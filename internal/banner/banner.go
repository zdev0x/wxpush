package banner

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

const logo = `
██╗    ██╗██╗  ██╗██████╗ ██╗   ██╗███████╗██╗  ██╗
██║    ██║╚██╗██╔╝██╔══██╗██║   ██║██╔════╝██║  ██║
██║ █╗ ██║ ╚███╔╝ ██████╔╝██║   ██║███████╗███████║
██║███╗██║ ██╔██╗ ██╔═══╝ ██║   ██║╚════██║██╔══██║
╚███╔███╔╝██╔╝ ██╗██║     ╚██████╔╝███████║██║  ██║
 ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝      ╚═════╝ ╚══════╝╚═╝  ╚═╝`

// ShowBanner 显示启动banner
func ShowBanner(version, commit, date string, cfg map[string]string) {
	// 创建颜色输出
	titleColor := color.New(color.FgHiCyan, color.Bold)
	infoColor := color.New(color.FgHiWhite)
	successColor := color.New(color.FgHiGreen)
	urlColor := color.New(color.FgHiBlue, color.Underline)
	separatorColor := color.New(color.FgHiBlack)
	groupColor := color.New(color.FgHiYellow)

	// 显示LOGO
	fmt.Println(titleColor.Sprint(logo))
	fmt.Println()

	// 获取所有配置项
	maxKeyLen := 0
	for k := range cfg {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	// 版本信息
	fmt.Printf("  %s\n", groupColor.Sprint("Version Information:"))
	printConfigItem(cfg, "Version", maxKeyLen, infoColor, separatorColor)
	if commit != "none" {
		printConfigItem(cfg, "Commit", maxKeyLen, infoColor, separatorColor)
	}
	if date != "unknown" {
		printConfigItem(cfg, "Build Time", maxKeyLen, infoColor, separatorColor)
	}
	fmt.Println()

	// 服务配置
	fmt.Printf("  %s\n", groupColor.Sprint("Service Configuration:"))
	printConfigItem(cfg, "Run Mode", maxKeyLen, successColor, separatorColor)
	printConfigItem(cfg, "Listen Address", maxKeyLen, urlColor, separatorColor)
	printConfigItem(cfg, "Log File", maxKeyLen, infoColor, separatorColor)
	fmt.Println()

	// 微信配置
	fmt.Printf("  %s\n", groupColor.Sprint("WeChat Configuration:"))
	printConfigItem(cfg, "WeChat App ID", maxKeyLen, infoColor, separatorColor)
	printConfigItem(cfg, "Message Templates", maxKeyLen, infoColor, separatorColor)
	printConfigItem(cfg, "WeChat Users", maxKeyLen, infoColor, separatorColor)
	printConfigItem(cfg, "Notify Groups", maxKeyLen, infoColor, separatorColor)
	fmt.Println()

	// 显示启动成功信息
	fmt.Printf("  %s %s\n",
		successColor.Sprint("✓"),
		infoColor.Sprint("Server started successfully"),
	)
	fmt.Printf("  %s %s\n",
		successColor.Sprint("✓"),
		infoColor.Sprintf("Start time: %s", time.Now().Format("2006-01-02 15:04:05")),
	)
	fmt.Println()
}

// printConfigItem 打印配置项
func printConfigItem(cfg map[string]string, key string, maxKeyLen int, valueColor, separatorColor *color.Color) {
	if value, ok := cfg[key]; ok {
		padding := strings.Repeat(" ", maxKeyLen-len(key))
		fmt.Printf("  %s%s %s %s\n",
			key,
			padding,
			separatorColor.Sprint("→"),
			valueColor.Sprint(value),
		)
	}
}
