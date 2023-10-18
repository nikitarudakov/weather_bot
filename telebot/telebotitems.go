package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherapi"
	tele "gopkg.in/telebot.v3"
	"time"
)

var menu = &tele.ReplyMarkup{ResizeKeyboard: true}

var lastLatStored, lastLonStored string
var weatherForecast *weatherapi.Response

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
