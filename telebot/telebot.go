package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
	"time"
)

// InitTelegramBot initializes Telegram Weather Bot
func InitTelegramBot(cfg *config.Config, dbClient db.DatabaseAccessor) {
	pref := tele.Settings{
		Token:     viper.GetString("telegram.token"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: "HTML",
	}

	// connect to telegram bot
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	// create menu with dates starting today and ending on the day 7 days ahead
	menuDateBtnSlice := getMenuDateBtnSlice()

	menu.Reply(
		menu.Row(menuDateBtnSlice[0]),
		menu.Row(menuDateBtnSlice[1], menuDateBtnSlice[2]),
		menu.Row(menuDateBtnSlice[3], menuDateBtnSlice[4]),
		menu.Row(menuDateBtnSlice[5], menuDateBtnSlice[6]),
		menu.Row(menuDateBtnSlice[7]),
	)

	// handles event triggered by start command
	b.Handle("/start", func(context tele.Context) error {
		return handleStartCmd(context)
	})

	// handles event triggered by subscribe command
	b.Handle("/subscribe", func(context tele.Context) error {
		return handleSubscriptionCmd(context, dbClient)
	})

	// handles event triggered by any sent text after subscribe cmd was triggered
	b.Handle(tele.OnText, func(context tele.Context) error {
		return handleTimeMessageForSubscription(context, dbClient)
	})

	// handles event triggered when Location Pin is sent
	b.Handle(tele.OnLocation, func(context tele.Context) error {
		return handleLocationPinMessage(context, cfg)
	})

	// handles event triggered when whole period btn is pressed on menu
	b.Handle(&menuDateBtnSlice[0], func(context tele.Context) error {
		return handleDateWholePeriodBtn(context)
	})

	// handles event triggered when day btn is pressed on menu
	for dtBtnIndex := 1; dtBtnIndex < 8; dtBtnIndex++ {
		b.Handle(&menuDateBtnSlice[dtBtnIndex], func(context tele.Context) error {
			return handleDateBtn(context, dtBtnIndex)
		})
	}

	b.Start()

	log.Info().Msg("Telegram bot is running ...")
}
