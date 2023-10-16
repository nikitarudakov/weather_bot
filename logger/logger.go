package logger

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

// InitLogger initializes logger with configurations from variable of type *config.Config
func InitLogger(cfg *config.Config) {
	output := zerolog.ConsoleWriter{Out: os.Stderr}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %s |", i))
	}

	loggingLevel := cfg.Logger.LoggingLevel
	log.Logger = log.
		Output(output).
		Level(zerolog.Level(loggingLevel))
}
