package configs

import (
	"github.com/walkerdu/wecom-backend/pkg/chatbot"
	"github.com/walkerdu/wecom-backend/pkg/wecom"
)

type NotifyConfig struct {
	UserID string `json:"user_id"`
	// 群ID
}

// 企业微信配置
type WeComConfig struct {
	AgentConfig  wecom.AgentConfig       `json:"agent_config"`
	Addr         string                  `json:"addr"`
	NotifyConfig map[string]NotifyConfig `json:"notify_config,omitempty"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Config struct {
	WeCom  WeComConfig          `json:"we_com"`
	OpenAI chatbot.OpenAIConfig `json:"open_ai"`
	Redis  RedisConfig          `json:"redis"`
}
