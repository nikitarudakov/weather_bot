package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherbotdb"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
	"time"
)

var menu = &tele.ReplyMarkup{ResizeKeyboard: true}

func getDtBtnSlice() []tele.Btn {
	dtBtnSlice := make([]tele.Btn, 8)

	dtBtnSlice[0] = menu.Text("All 7 days")

	dtToday := time.Now()

	for dayPlus := 0; dayPlus < 7; dayPlus++ {
		dtStr := dtToday.AddDate(0, 0, dayPlus).Format("02/01/2006")
		dtBtnSlice[dayPlus+1] = menu.Text(dtStr)
	}

	return dtBtnSlice
}

// InitTelegramBot initializes Telegram Weather Bot
func InitTelegramBot(cfg *config.Config, dbClient *weatherbotdb.WeatherBotClientDb) {
	pref := tele.Settings{
		Token:     viper.GetString("telegram.token"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: "HTML",
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Error().
			Str("service", "NewBot").
			Err(err).
			Msg("error initializing bot")

		return
	}

	log.Info().Msg("telegram bot was successfully initialized")

	b.Handle("/start", func(c tele.Context) error {
		if err := insertChatHistoryToDb(c, dbClient, &cfg.Db); err != nil {
			log.Error().
				Err(err).
				Str("service", "dbClient.InsertDocToDbCollection").
				Msg("failed to insert document")
		}

		return c.Send("Welcome to <b>Weather Bot</b>!\nUse this bot to see Weather Forecast in your area :)\nJust send <i>location pin</i> to Weather bot and get accurate forecast for up to <b>7 days</b> forward!")
	})

	b.Handle("/subscribe", func(c tele.Context) error {
		if err := insertChatHistoryToDb(c, dbClient, &cfg.Db); err != nil {
			log.Error().
				Err(err).
				Str("service", "dbClient.InsertDocToDbCollection").
				Msg("failed to insert document")
		}

		return c.Send(`To subscribe for daily weather forecast 
send time (format: 15:04) you wish to receive it at`)
	})

	// Create menu with dates starting today and ending on the day 7 days ahead
	dtBtnSlice := getDtBtnSlice()

	menu.Reply(
		menu.Row(dtBtnSlice[0]),
		menu.Row(dtBtnSlice[1], dtBtnSlice[2]),
		menu.Row(dtBtnSlice[3], dtBtnSlice[4]),
		menu.Row(dtBtnSlice[5], dtBtnSlice[6]),
		menu.Row(dtBtnSlice[7]),
	)
	// --------------------------------------------------------------------------

	b.Handle(tele.OnLocation, func(c tele.Context) error {
		if err := insertChatHistoryToDb(c, dbClient, &cfg.Db); err != nil {
			log.Error().
				Err(err).
				Str("service", "dbClient.InsertDocToDbCollection").
				Msg("failed to insert document")
		}

		err = handleOnLocation(cfg, c)
		if err != nil {
			log.Error().Err(err).Msg("error handling on location message request")
			return err
		}

		return nil
	})

	b.Handle(&dtBtnSlice[0], func(c tele.Context) error {
		if err := insertChatHistoryToDb(c, dbClient, &cfg.Db); err != nil {
			log.Error().
				Err(err).
				Str("service", "dbClient.InsertDocToDbCollection").
				Msg("failed to insert document")
		}

		err = handleWholePeriodBtn(c)
		if err != nil {
			log.Error().Err(err).Msg("error handling on location message request")
			return err
		}

		return nil
	})

	for i := 1; i < 8; i++ {
		dtBtnIndex := i
		b.Handle(&dtBtnSlice[dtBtnIndex], func(c tele.Context) error {
			if err := insertChatHistoryToDb(c, dbClient, &cfg.Db); err != nil {
				log.Error().
					Err(err).
					Str("service", "dbClient.InsertDocToDbCollection").
					Msg("failed to insert document")
			}

			err = handleDateBtn(c, dtBtnIndex)
			if err != nil {
				return err
			}

			return nil
		})
	}

	b.Start()
}
