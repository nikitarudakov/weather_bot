package config

import (
	"github.com/spf13/viper"
	"log"
)

type telegramCfg struct {
	Token string `mapstructure:"token"`
}

type WeatherAPICfg struct {
	Server  string `mapstructure:"server"`
	Token   string `mapstructure:"token"`
	Exclude string `mapstructure:"exclude"`
	Units   string `mapstructure:"units"`
}

type loggerCfg struct {
	LoggingLevel int8 `mapstructure:"logging_level"`
}

// Config Create private data struct to hold config options.
type Config struct {
	Telegram   telegramCfg   `mapstructure:"telegram"`
	WeatherAPI WeatherAPICfg `mapstructure:"weather_api"`
	Logger     loggerCfg     `mapstructure:"logger"`
}

// GetConfig parses .json file to Config struct
func GetConfig() (*Config, error) {
	viper.SetConfigName(".config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		log.Println("config could not be loaded, ERROR")
		return nil, err
	}

	cfg := &Config{}

	// Parse config to struct
	if err := viper.Unmarshal(cfg); err != nil {
		log.Println("config could not be parsed, ERROR")
		return nil, err
	}

	return cfg, nil
}
