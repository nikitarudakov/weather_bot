package weatherbotdb

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"testing"
)

func TestDbConnection(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Errorf("error parsing config %s", err)
	}

	_, err = NewWeatherBotDbClient(&cfg.Db)
	if err != nil {
		t.Error(err)
	}
}
