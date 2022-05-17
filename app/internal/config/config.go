package config

import (
	"github.com/hawkkiller/k121_bot/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Storage  StorageConfig  `yaml:"storage"`
}

type StorageConfig struct {
	Host        string `yaml:"host" env:"DATABASE_HOST" env-required:"true"`
	Port        string `yaml:"port" env:"DATABASE_PORT" env-required:"true"`
	Database    string `yaml:"database" env:"DATABASE_NAME" env-required:"true"`
	Username    string `yaml:"username" env:"DATABASE_USERNAME" env-required:"true"`
	Password    string `yaml:"password" env:"DATABASE_PASSWORD" env-required:"true"`
	MaxAttempts int8   `yaml:"maxAttempts" env:"MAX_ATTEMPTS" env-default:"5"`
}

type TelegramConfig struct {
	Token string `yaml:"token" env:"TELEGRAM_TOKEN" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
