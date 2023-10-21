package weatherapi

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"git.foxminded.ua/foxstudent106092/weather-bot/logger"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/geoutils"
	"testing"
)

func TestStoreWeatherForecastForUser(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Logf("Error parsing config %s", err)
	}

	dbClient, _ := db.NewDatabaseClient(&cfg.Db)

	var userID int64 = 12345

	latStr := geoutils.FormatCoordinateToString(50.447731)
	lonStr := geoutils.FormatCoordinateToString(30.542721)

	weatherAPI := NewWeatherAPIService(&cfg.WeatherAPI)

	weatherForecastAtLocation, err := weatherAPI.GetWeatherForecast(latStr, lonStr)
	if err != nil {
		t.Error(err)
	}

	if err = weatherForecastAtLocation.StoreUpdateWeatherForecast(dbClient, &cfg.Db, userID); err != nil {
		t.Error(err)
	}
}

func TestWeatherAPI(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Logf("Error parsing config %s", err)
	}

	logger.InitLogger(cfg)

	testCases := []struct {
		lat float32
		lon float32
	}{
		{40.712776, -74.005974},
		{51.5073516, -0.127758},
		{50.447731, 30.542721},
	}

	for _, testCase := range testCases {
		testCase := testCase

		latStr := geoutils.FormatCoordinateToString(testCase.lat)
		lonStr := geoutils.FormatCoordinateToString(testCase.lon)
		testName := fmt.Sprintf("API test for coords: (%s, %s)", latStr, lonStr)

		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			weatherAPI := NewWeatherAPIService(&cfg.WeatherAPI)

			weatherForecastAtLocation, err := weatherAPI.GetWeatherForecast(latStr, lonStr)
			if err != nil {
				t.Error(err)
			}

			t.Log(fmt.Sprintf("%+v\n", *weatherForecastAtLocation))
		})
	}
}
