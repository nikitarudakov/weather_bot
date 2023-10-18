package weatherapi

import (
	"encoding/json"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

// DailyTemp stores fields of max and min temperature
type DailyTemp struct {
	Min float64 `json:"Min"`
	Max float64 `json:"max"`
}

// Weather stores field of weather description
type Weather struct {
	Desc string `json:"description,omitempty"`
}

// DailyWeather stores fields of weather features at UnixDt day
type DailyWeather struct {
	UnixDt  int64     `json:"dt"`
	Sunrise int64     `json:"sunrise"`
	Sunset  int64     `json:"sunset"`
	Summary string    `json:"summary"`
	Temp    DailyTemp `json:"temp"`
	Weather []Weather `json:"weather"`
}

// Response represents Weather API response of daily forecasts
type Response struct {
	Daily []DailyWeather `json:"daily"`
}

// WeatherAPI represents the methods for interacting with a weather API
type WeatherAPI interface {
	GetWeatherForecast(lat, lon string) (*Response, error)
}

// WeatherService represents a service that interacts with a weather API
type WeatherService struct {
	cfg *config.WeatherAPICfg
}

// NewWeatherService initialized new WeatherService instance
func NewWeatherService(cfg *config.WeatherAPICfg) WeatherAPI {
	var weatherApi WeatherAPI = &WeatherService{cfg: cfg}
	return weatherApi
}

func (wa *WeatherService) getAPIUrl(lat, lon string) string {
	return wa.cfg.Server +
		"?lat=" + lat +
		"&lon=" + lon +
		"&appid=" + wa.cfg.Token +
		"&exclude=" + wa.cfg.Exclude +
		"&units=" + wa.cfg.Units
}

// GetWeatherForecast fetches API and returns Response - which is weather forecast
func (wa *WeatherService) GetWeatherForecast(lat, lon string) (*Response, error) {
	// Get API url for location coords lat/lon
	apiURL := wa.getAPIUrl(lat, lon)

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
