package config

import (
	"github.com/spf13/viper"
	"log"
)

type TelegramCfg struct {
	Token string `mapstructure:"token"`
}

// DbCfg stores values for db connection
type DbCfg struct {
	Name                   string `mapstructure:"name"`
	SubsCollectionName     string `mapstructure:"subscription_collection_name"`
	ForecastCollectionName string `mapstructure:"forecast_collection_name"`
	ConnectionURL          string `mapstructure:"connection_url"`
}

// WeatherAPICfg stores values for api requests query params
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
	Telegram   TelegramCfg   `mapstructure:"telegram"`
	Db         DbCfg         `mapstructure:"db"`
	WeatherAPI WeatherAPICfg `mapstructure:"weather_api"`
	Logger     loggerCfg     `mapstructure:"logger"`
}

// GetConfig parses .json file to Config struct
func GetConfig() (*Config, error) {
	viper.SetConfigName(".config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	// read config
	if err := viper.ReadInConfig(); err != nil {
		log.Println("config could not be loaded, ERROR")
		return nil, err
	}

	cfg := &Config{}

	// parse config to struct
	if err := viper.Unmarshal(cfg); err != nil {
		log.Println("config could not be parsed, ERROR")
		return nil, err
	}

	return cfg, nil
}
