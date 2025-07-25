package wechat

import (
	// #nosec G505
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/zdev0x/wxpush/internal/config"
	"github.com/zdev0x/wxpush/internal/logger"
	"github.com/zdev0x/wxpush/internal/model"
)

var (
	accessToken    string
	accessTokenExp time.Time
	tokenMutex     sync.Mutex
)

// GetAccessToken 获取access_token
func GetAccessToken(cfg *config.Config) (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// 检查缓存的token是否有效
	if accessToken != "" && time.Now().Before(accessTokenExp) {
		return accessToken, nil
	}

	// 请求新token
	// url 来源于微信官方API，安全
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		cfg.WeChat.AppID,
		cfg.WeChat.AppSecret,
	)
	// #nosec G107
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("token_refresh", "", model.ErrInternal, err, nil)
		return "", fmt.Errorf("请求access_token失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("token_refresh", "", model.ErrInternal, err, nil)
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("token_refresh", "", model.ErrInternal, err, nil)
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		logger.Error("token_refresh", "", model.ErrInternal, fmt.Errorf(result.ErrMsg), map[string]interface{}{
			"errcode": result.ErrCode,
			"errmsg":  result.ErrMsg,
		})
		return "", fmt.Errorf("获取access_token失败: %s", result.ErrMsg)
	}

	// 更新缓存
	accessToken = result.AccessToken
	accessTokenExp = time.Now().Add(time.Duration(result.ExpiresIn-300) * time.Second)

	logger.Info("token_refresh", "", map[string]interface{}{
		"expires_in": result.ExpiresIn,
		"expire_at":  accessTokenExp.Format(time.RFC3339),
	})

	return accessToken, nil
}

// CheckSignature 检查微信服务器签名
func CheckSignature(cfg *config.Config, signature, timestamp, nonce string) bool {
	// 1. 将token、timestamp、nonce三个参数进行字典序排序
	params := []string{cfg.WeChat.Token, timestamp, nonce}
	sort.Strings(params)

	// 2. 将三个参数字符串拼接成一个字符串进行sha1加密
	// 微信官方要求sha1算法，风险可接受
	str := strings.Join(params, "")
	h := sha1.New() // #nosec G401
	h.Write([]byte(str))
	sign := hex.EncodeToString(h.Sum(nil))

	// 3. 开发者获得加密后的字符串可与signature对比，标识该请求来源于微信
	return sign == signature
}
