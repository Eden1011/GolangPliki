package main

import (
	"time"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CurrentWeather struct {
	Time                string  `json:"time"`
	Temperature         float64 `json:"temperature_2m"`
	RelativeHumidity    int     `json:"relative_humidity_2m"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	Precipitation       float64 `json:"precipitation"`
	WindSpeed           float64 `json:"wind_speed_10m"`
	WindDirection       float64 `json:"wind_direction_10m"`
	CloudCover          int     `json:"cloud_cover"`
	Pressure            float64 `json:"surface_pressure"`
}

type CurrentWeatherResponse struct {
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	Timezone  string         `json:"timezone"`
	Current   CurrentWeather `json:"current"`
}

type DailyWeather struct {
	Time            []string  `json:"time"`
	TemperatureMax  []float64 `json:"temperature_2m_max"`
	TemperatureMin  []float64 `json:"temperature_2m_min"`
	ApparentTempMax []float64 `json:"apparent_temperature_max"`
	ApparentTempMin []float64 `json:"apparent_temperature_min"`
	Precipitation   []float64 `json:"precipitation_sum"`
	WindSpeedMax    []float64 `json:"wind_speed_10m_max"`
	WindDirection   []float64 `json:"wind_direction_10m_dominant"`
	UVIndexMax      []float64 `json:"uv_index_max"`
}

type ForecastResponse struct {
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
	Timezone  string       `json:"timezone"`
	Daily     DailyWeather `json:"daily"`
}

type HistoryResponse struct {
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
	Timezone  string       `json:"timezone"`
	Daily     DailyWeather `json:"daily"`
}

// Struktury dla analizy zagrożeń
type WeatherThreat struct {
	Date        string
	Type        string
	Description string
	Value       float64
}

type Config struct {
	Threats struct {
		HighTemperature   float64 `json:"high_temperature"`
		LowTemperature    float64 `json:"low_temperature"`
		HighWindSpeed     float64 `json:"high_wind_speed"`
		HighPrecipitation float64 `json:"high_precipitation"`
		HighUVIndex       float64 `json:"high_uv_index"`
	} `json:"threats"`
}

type City struct {
	Name      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func (cw *CurrentWeather) GetFormattedTime() string {
	t, err := time.Parse(time.RFC3339, cw.Time)
	if err != nil {
		return cw.Time
	}
	return t.Format("2006-01-02 15:04")
}

func (cw *CurrentWeather) GetWindDirectionText() string {
	direction := cw.WindDirection
	switch {
	case direction >= 337.5 || direction < 22.5:
		return "Północny"
	case direction >= 22.5 && direction < 67.5:
		return "Północno-wschodni"
	case direction >= 67.5 && direction < 112.5:
		return "Wschodni"
	case direction >= 112.5 && direction < 157.5:
		return "Południowo-wschodni"
	case direction >= 157.5 && direction < 202.5:
		return "Południowy"
	case direction >= 202.5 && direction < 247.5:
		return "Południowo-zachodni"
	case direction >= 247.5 && direction < 292.5:
		return "Zachodni"
	case direction >= 292.5 && direction < 337.5:
		return "Północno-zachodni"
	default:
		return "Nieznany"
	}
}

func (cw *CurrentWeather) GetWeatherIcon() string {
	return ""
}
