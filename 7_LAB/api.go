package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	OpenMeteoBaseURL = "https://api.open-meteo.com/v1"
	GeocodingURL     = "https://geocoding-api.open-meteo.com/v1/search"
)

func GetCityCoordinates(cityName string) (*Coordinates, error) {
	params := url.Values{}
	params.Add("name", cityName)
	params.Add("count", "1")
	params.Add("language", "pl")
	params.Add("format", "json")

	fullURL := fmt.Sprintf("%s?%s", GeocodingURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("błąd połączenia z API geocoding: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("błąd HTTP: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("błąd odczytu odpowiedzi: %w", err)
	}

	var geoResponse struct {
		Results []struct {
			Name      string  `json:"name"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &geoResponse)
	if err != nil {
		return nil, fmt.Errorf("błąd parsowania JSON: %w", err)
	}

	if len(geoResponse.Results) == 0 {
		return nil, fmt.Errorf("nie znaleziono miasta: %s", cityName)
	}

	result := geoResponse.Results[0]
	return &Coordinates{
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
	}, nil
}

func GetCurrentWeather(lat, lon float64) (*CurrentWeatherResponse, error) {
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%.6f", lat))
	params.Add("longitude", fmt.Sprintf("%.6f", lon))
	params.Add("current", strings.Join([]string{
		"temperature_2m",
		"relative_humidity_2m",
		"apparent_temperature",
		"precipitation",
		"wind_speed_10m",
		"wind_direction_10m",
		"cloud_cover",
		"surface_pressure",
	}, ","))
	params.Add("timezone", "auto")

	fullURL := fmt.Sprintf("%s/forecast?%s", OpenMeteoBaseURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("błąd połączenia z API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("błąd HTTP: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("błąd odczytu odpowiedzi: %w", err)
	}

	var weather CurrentWeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, fmt.Errorf("błąd parsowania JSON: %w", err)
	}

	return &weather, nil
}

func GetWeatherForecast(lat, lon float64, days int) (*ForecastResponse, error) {
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%.6f", lat))
	params.Add("longitude", fmt.Sprintf("%.6f", lon))
	params.Add("daily", strings.Join([]string{
		"temperature_2m_max",
		"temperature_2m_min",
		"apparent_temperature_max",
		"apparent_temperature_min",
		"precipitation_sum",
		"wind_speed_10m_max",
		"wind_direction_10m_dominant",
		"uv_index_max",
	}, ","))
	params.Add("timezone", "auto")
	params.Add("forecast_days", fmt.Sprintf("%d", days))

	fullURL := fmt.Sprintf("%s/forecast?%s", OpenMeteoBaseURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("błąd połączenia z API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("błąd HTTP: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("błąd odczytu odpowiedzi: %w", err)
	}

	var forecast ForecastResponse
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return nil, fmt.Errorf("błąd parsowania JSON: %w", err)
	}

	return &forecast, nil
}

func GetWeatherHistory(lat, lon float64, startDate, endDate string) (*HistoryResponse, error) {
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%.6f", lat))
	params.Add("longitude", fmt.Sprintf("%.6f", lon))
	params.Add("start_date", startDate)
	params.Add("end_date", endDate)
	params.Add("daily", strings.Join([]string{
		"temperature_2m_max",
		"temperature_2m_min",
		"apparent_temperature_max",
		"apparent_temperature_min",
		"precipitation_sum",
		"wind_speed_10m_max",
		"wind_direction_10m_dominant",
		"uv_index_max",
	}, ","))
	params.Add("timezone", "auto")

	fullURL := fmt.Sprintf("%s/archive?%s", OpenMeteoBaseURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("błąd połączenia z API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("błąd HTTP: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("błąd odczytu odpowiedzi: %w", err)
	}

	var history HistoryResponse
	err = json.Unmarshal(body, &history)
	if err != nil {
		return nil, fmt.Errorf("błąd parsowania JSON: %w", err)
	}

	return &history, nil
}
