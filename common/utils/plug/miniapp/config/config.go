package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

// MiniProgram 小程序配置，与 PowerWeChat 库中的 UserConfig 兼容
type MiniProgram struct {
	AppID     string `json:"app_id"`
	Secret    string `json:"secret"`
	Token     string `json:"token"`
	AESKey    string `json:"aes_key"`
	MchID     string `json:"mch_id"`
	KeyPath   string `json:"key_path"`
	CertPath  string `json:"cert_path"`
	Key       string `json:"key"`
}

// ToPowerWeChatConfig 将自定义配置转换为 PowerWeChat 库的 UserConfig 类型
func (c *MiniProgram) ToPowerWeChatConfig() *miniProgram.UserConfig {
	return &miniProgram.UserConfig{
		AppID:  c.AppID,
		Secret: c.Secret,
		Token:  c.Token,
		AESKey: c.AESKey,
		// 注意：PowerWeChat 的 miniProgram.UserConfig 结构体中不包含 MchID、KeyPath、CertPath、Key 字段
		// 这些字段通常用于支付功能，可能在其他模块中配置
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(filePath string) (*MiniProgram, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config MiniProgram
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() (*MiniProgram, error) {
	config := &MiniProgram{
		AppID:     os.Getenv("MINIAPP_APPID"),
		Secret:    os.Getenv("MINIAPP_SECRET"),
		Token:     os.Getenv("MINIAPP_TOKEN"),
		AESKey:    os.Getenv("MINIAPP_AESKEY"),
		MchID:     os.Getenv("MINIAPP_MCHID"),
		KeyPath:   os.Getenv("MINIAPP_KEYPATH"),
		CertPath:  os.Getenv("MINIAPP_CERTPATH"),
		Key:       os.Getenv("MINIAPP_KEY"),
	}

	return config, nil
}

