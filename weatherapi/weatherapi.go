package weatherapi

import (
	"encoding/json"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type DailyTemp struct {
	Min float64 `json:"Min"`
	Max float64 `json:"max"`
}

type Weather struct {
	Desc string `json:"description,omitempty"`
}

type DailyWeather struct {
	UnixDt  int64     `json:"dt"`
	Sunrise int64     `json:"sunrise"`
	Sunset  int64     `json:"sunset"`
	Summary string    `json:"summary"`
	Temp    DailyTemp `json:"temp"`
	Weather []Weather `json:"weather"`
}

type Response struct {
	Daily []DailyWeather `json:"daily"`
}

func GetAPIUrl(cfg *config.Config, lat, lon string) string {
	return cfg.WeatherAPI.Server +
		"?lat=" + lat +
		"&lon=" + lon +
		"&appid=" + cfg.WeatherAPI.Token +
		"&exclude=" + cfg.WeatherAPI.Exclude +
		"&units=" + cfg.WeatherAPI.Units
}

func GetWeatherForecast(apiURL string) (*Response, error) {
	log.Trace().Str("service", "Weather API").Str("api_url", apiURL).Send()

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Error().Str("service", "http.Get").Err(err).Send()
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Str("service", "io.ReadAll").Err(err).Send()
		return nil, err
	}

	weatherResp := Response{}
	if err = json.Unmarshal(body, &weatherResp); err != nil {
		log.Error().Str("service", "json.Unmarshal").Err(err).Send()
		return nil, err
	}

	return &weatherResp, nil
}
