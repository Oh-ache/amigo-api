package miniapp

import (
	"sync"

	"amigo-api/common/utils/plug/miniapp/config"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

// 线程安全的小程序客户端存储
var (
	clients sync.Map
)

// MiniAppConfig 小程序配置
type MiniAppConfig struct {
	AppID    string `json:"app_id"`
	Secret   string `json:"secret"`
	Token    string `json:"token"`
	AESKey   string `json:"aes_key"`
	MchID    string `json:"mch_id"`
	KeyPath  string `json:"key_path"`
	CertPath string `json:"cert_path"`
	Key      string `json:"key"`
}

// Client 小程序客户端类型
type Client struct {
	*miniProgram.MiniProgram
}

// InitClient 初始化小程序客户端
func InitClient(miniappName string, cfg *MiniAppConfig) (*Client, error) {
	// 尝试从缓存中获取
	if client, ok := clients.Load(miniappName); ok {
		return client.(*Client), nil
	}

	// 创建新的小程序客户端配置
	mpConfig := &config.MiniProgram{
		AppID:    cfg.AppID,
		Secret:   cfg.Secret,
		Token:    cfg.Token,
		AESKey:   cfg.AESKey,
		MchID:    cfg.MchID,
		KeyPath:  cfg.KeyPath,
		CertPath: cfg.CertPath,
		Key:      cfg.Key,
	}

	// 初始化小程序客户端
	mp, err := miniProgram.NewMiniProgram(mpConfig.ToPowerWeChatConfig())
	if err != nil {
		return nil, err
	}

	// 存储到缓存中
	client := &Client{mp}
	clients.Store(miniappName, client)

	return client, nil
}

// GetClient 获取小程序客户端
func GetClient(miniappName string) (*Client, error) {
	if client, ok := clients.Load(miniappName); ok {
		return client.(*Client), nil
	}

	return nil, nil
}

// RemoveClient 移除小程序客户端
func RemoveClient(miniappName string) {
	clients.Delete(miniappName)
}

// ListClients 列出所有小程序客户端
func ListClients() []string {
	var names []string
	clients.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}
