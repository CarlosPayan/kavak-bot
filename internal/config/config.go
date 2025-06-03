package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

type OpenAIConfig struct {
	APIKey string `mapstructure:"api_key"`
}

type CatalogConfig struct {
	Path string `mapstructure:"path"`
}

type TwilioConfig struct {
	AccountSID   string `mapstructure:"account_sid"`
	AuthToken    string `mapstructure:"auth_token"`
	WhatsAppFrom string `mapstructure:"whatsapp_from"`
}

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	OpenAI  OpenAIConfig  `mapstructure:"openai"`
	Catalog CatalogConfig `mapstructure:"catalog"`
	Twilio  TwilioConfig  `mapstructure:"twilio"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
