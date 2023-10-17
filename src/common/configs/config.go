package configs

import (
	"github.com/spf13/viper"
	"kits/api/src/core/enums"
	"time"
)

type Config struct {
	Mode   enums.AppMode `mapstructure:"mode"`
	Server struct {
		Name   string `mapstructure:"name"`
		Port   string `mapstructure:"port"`
		Prefix string `mapstructure:"prefix"`
	} `mapstructure:"server"`
	Postgresql struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		User        string `mapstructure:"user"`
		Password    string `mapstructure:"password"`
		DbName      string `mapstructure:"db_name"`
		SslMode     string `mapstructure:"ssl_mode"`
		AutoMigrate bool   `mapstructure:"auto_migrate"`
		MaxLifeTime int    `mapstructure:"max_life_time"`
	} `mapstructure:"postgresql"`
	Redis struct {
		Hosts    []string `mapstructure:"hosts"`
		Username string   `mapstructure:"username"`
		Password string   `mapstructure:"password"`
	} `mapstructure:"redis"`
	Aes map[string]struct {
		Key    string        `mapstructure:"key"`
		Expire time.Duration `mapstructure:"expire"`
	} `mapstructure:"aes"`
}

func NewAppConfig(pathConfig string) (*Config, error) {
	var common *Config
	viper.SetConfigFile(pathConfig)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&common)

	return common, nil
}
