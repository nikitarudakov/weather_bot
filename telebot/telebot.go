package telebot

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/geoutils"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherapi"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
	"strings"
	"time"
)

func InitTelegramBot(cfg *config.Config) {
	pref := tele.Settings{
		Token:     viper.GetString("telegram.token"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: "Markdown",
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Error().
			Str("service", "InitTelegramBot").
			Err(err).
			Msg("error initializing bot")

		return
	}

	log.Info().Msg("telegram bot was successfully initialized")

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Welcome to *Weather Bot*!\nUse this bot to see Weather Forecast in your area :)")
	})

	b.Handle(tele.OnLocation, func(c tele.Context) error {
		lat, lon := c.Message().Location.Lat, c.Message().Location.Lng
		latStr, lonStr := geoutils.FormatCoordinateToString(lat), geoutils.FormatCoordinateToString(lon)

		apiURL := weatherapi.GetAPIUrl(cfg, latStr, lonStr)

		resp, err := weatherapi.GetWeatherForecast(apiURL)
		if err != nil {
			log.Error().
				Str("service", "weatherapi.GetWeatherForecast").
				Err(err).
				Msg("failed to get weather forecast")
		}

		fmt.Println(resp.Daily)
		if len(resp.Daily) == 0 {
			fmt.Println("HERERERERERER")
			return c.Send("Data unavailable. Try again")
		}

		var weatherFormattedTextMsg []string
		for _, dailyWeather := range resp.Daily {
			weatherTextMsg := dailyWeather.FormatToTextMsg()

			weatherFormattedTextMsg = append(weatherFormattedTextMsg, weatherTextMsg)
		}

		weatherStringsJoined := strings.Join(weatherFormattedTextMsg, "\n")

		return c.Send(weatherStringsJoined)
	})

	b.Start()
}
