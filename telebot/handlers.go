package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/geoutils"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherapi"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
	"strings"
)

func handleOnLocation(cfg *config.Config, c tele.Context) error {
	lat, lon := c.Message().Location.Lat, c.Message().Location.Lng
	lastLatStored, lastLonStored =
		geoutils.FormatCoordinateToString(lat),
		geoutils.FormatCoordinateToString(lon)

	weatherAPI := weatherapi.NewWeatherService(&cfg.WeatherAPI)

	weatherForecastAtLocation, err := weatherAPI.GetWeatherForecast(lastLatStored, lastLonStored)
	if err != nil {
		log.Error().Err(err).Send()

		return c.Send("Data is unavailable for this location!")
	}

	weatherForecast = weatherForecastAtLocation

	return c.Send("Choose time period to get forecast:", menu)
}

func handleWholePeriodBtn(c tele.Context) error {
	if weatherForecast == nil {
		return c.Send("Data is unavailable!\nSend location pin")
	}

	var dailyWeatherBuilder strings.Builder
	for _, dailyWeather := range weatherForecast.Daily {
		weatherTextMsg, err := dailyWeather.FormatToTextMsg()
		if err != nil {
			log.Warn().Err(err).Send()
		}

		dailyWeatherBuilder.WriteString(weatherTextMsg)
		dailyWeatherBuilder.WriteString("\n")
	}

	return c.Send(dailyWeatherBuilder.String())
}

func handleDateBtn(c tele.Context, dtBtnIndex int) error {
	if weatherForecast == nil {
		return c.Send("Data is unavailable!\nSend location pin")
	}

	weatherTextMsg, err := weatherForecast.Daily[dtBtnIndex-1].FormatToTextMsg()
	if err != nil {
		log.Warn().Err(err).Send()
	}

	return c.Send(weatherTextMsg)
}
