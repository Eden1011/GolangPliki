package main

import (
	"fmt"
)

func AnalyzeWeatherThreats(forecast *ForecastResponse, config *Config) []WeatherThreat {
	var threats []WeatherThreat

	for i := 0; i < len(forecast.Daily.Time); i++ {
		date := forecast.Daily.Time[i]

		if forecast.Daily.TemperatureMax[i] > config.Threats.HighTemperature {
			threats = append(threats, WeatherThreat{
				Date: date,
				Type: "Wysoka temperatura",
				Description: fmt.Sprintf("Temperatura maksymalna %.1f°C przekracza próg %.1f°C",
					forecast.Daily.TemperatureMax[i], config.Threats.HighTemperature),
				Value: forecast.Daily.TemperatureMax[i],
			})
		}

		if forecast.Daily.TemperatureMin[i] < config.Threats.LowTemperature {
			threats = append(threats, WeatherThreat{
				Date: date,
				Type: "Niska temperatura",
				Description: fmt.Sprintf("Temperatura minimalna %.1f°C poniżej progu %.1f°C",
					forecast.Daily.TemperatureMin[i], config.Threats.LowTemperature),
				Value: forecast.Daily.TemperatureMin[i],
			})
		}

		if forecast.Daily.WindSpeedMax[i] > config.Threats.HighWindSpeed {
			threats = append(threats, WeatherThreat{
				Date: date,
				Type: "Silny wiatr",
				Description: fmt.Sprintf("Prędkość wiatru %.1f km/h przekracza próg %.1f km/h",
					forecast.Daily.WindSpeedMax[i], config.Threats.HighWindSpeed),
				Value: forecast.Daily.WindSpeedMax[i],
			})
		}

		if forecast.Daily.Precipitation[i] > config.Threats.HighPrecipitation {
			threats = append(threats, WeatherThreat{
				Date: date,
				Type: "Intensywne opady",
				Description: fmt.Sprintf("Opady %.1f mm przekraczają próg %.1f mm",
					forecast.Daily.Precipitation[i], config.Threats.HighPrecipitation),
				Value: forecast.Daily.Precipitation[i],
			})
		}

		if len(forecast.Daily.UVIndexMax) > i && forecast.Daily.UVIndexMax[i] > config.Threats.HighUVIndex {
			threats = append(threats, WeatherThreat{
				Date: date,
				Type: "Wysokie promieniowanie UV",
				Description: fmt.Sprintf("Indeks UV %.1f przekracza próg %.1f",
					forecast.Daily.UVIndexMax[i], config.Threats.HighUVIndex),
				Value: forecast.Daily.UVIndexMax[i],
			})
		}
	}

	consecutiveHotDays := 0
	consecutiveColdDays := 0

	for i := 0; i < len(forecast.Daily.Time); i++ {
		if forecast.Daily.TemperatureMax[i] > config.Threats.HighTemperature {
			consecutiveHotDays++
		} else {
			if consecutiveHotDays >= 3 {
				threats = append(threats, WeatherThreat{
					Date:        forecast.Daily.Time[i-consecutiveHotDays],
					Type:        "Fala upałów",
					Description: fmt.Sprintf("Przewidywane %d kolejnych dni z wysoką temperaturą", consecutiveHotDays),
					Value:       float64(consecutiveHotDays),
				})
			}
			consecutiveHotDays = 0
		}

		if forecast.Daily.TemperatureMin[i] < config.Threats.LowTemperature {
			consecutiveColdDays++
		} else {
			if consecutiveColdDays >= 3 {
				threats = append(threats, WeatherThreat{
					Date:        forecast.Daily.Time[i-consecutiveColdDays],
					Type:        "Fala mrozów",
					Description: fmt.Sprintf("Przewidywane %d kolejnych dni z niską temperaturą", consecutiveColdDays),
					Value:       float64(consecutiveColdDays),
				})
			}
			consecutiveColdDays = 0
		}
	}

	if consecutiveHotDays >= 3 {
		threats = append(threats, WeatherThreat{
			Date:        forecast.Daily.Time[len(forecast.Daily.Time)-consecutiveHotDays],
			Type:        "Fala upałów",
			Description: fmt.Sprintf("Przewidywane %d kolejnych dni z wysoką temperaturą", consecutiveHotDays),
			Value:       float64(consecutiveHotDays),
		})
	}

	if consecutiveColdDays >= 3 {
		threats = append(threats, WeatherThreat{
			Date:        forecast.Daily.Time[len(forecast.Daily.Time)-consecutiveColdDays],
			Type:        "Fala mrozów",
			Description: fmt.Sprintf("Przewidywane %d kolejnych dni z niską temperaturą", consecutiveColdDays),
			Value:       float64(consecutiveColdDays),
		})
	}

	return threats
}

func AnalyzeTemperatureTrends(forecast *ForecastResponse) map[string]interface{} {
	if len(forecast.Daily.TemperatureMax) == 0 {
		return nil
	}

	analysis := make(map[string]interface{})

	var sumMax, sumMin float64
	for i := 0; i < len(forecast.Daily.TemperatureMax); i++ {
		sumMax += forecast.Daily.TemperatureMax[i]
		sumMin += forecast.Daily.TemperatureMin[i]
	}

	avgMax := sumMax / float64(len(forecast.Daily.TemperatureMax))
	avgMin := sumMin / float64(len(forecast.Daily.TemperatureMin))

	analysis["average_max_temp"] = avgMax
	analysis["average_min_temp"] = avgMin

	maxTemp := forecast.Daily.TemperatureMax[0]
	minTemp := forecast.Daily.TemperatureMin[0]
	maxTempDate := forecast.Daily.Time[0]
	minTempDate := forecast.Daily.Time[0]

	for i := 1; i < len(forecast.Daily.TemperatureMax); i++ {
		if forecast.Daily.TemperatureMax[i] > maxTemp {
			maxTemp = forecast.Daily.TemperatureMax[i]
			maxTempDate = forecast.Daily.Time[i]
		}
		if forecast.Daily.TemperatureMin[i] < minTemp {
			minTemp = forecast.Daily.TemperatureMin[i]
			minTempDate = forecast.Daily.Time[i]
		}
	}

	analysis["max_temperature"] = maxTemp
	analysis["max_temperature_date"] = maxTempDate
	analysis["min_temperature"] = minTemp
	analysis["min_temperature_date"] = minTempDate

	trend := calculateTemperatureTrend(forecast.Daily.TemperatureMax)
	analysis["temperature_trend"] = trend

	return analysis
}

func calculateTemperatureTrend(temps []float64) float64 {
	if len(temps) < 2 {
		return 0
	}

	n := float64(len(temps))
	var sumX, sumY, sumXY, sumX2 float64

	for i, temp := range temps {
		x := float64(i)
		y := temp
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	numerator := n*sumXY - sumX*sumY
	denominator := n*sumX2 - sumX*sumX

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

func GetThreatSeverity(threat *WeatherThreat, config *Config) string {
	switch threat.Type {
	case "Wysoka temperatura":
		if threat.Value > config.Threats.HighTemperature+10 {
			return "Krytyczne"
		} else if threat.Value > config.Threats.HighTemperature+5 {
			return "Wysokie"
		}
		return "Średnie"
	case "Niska temperatura":
		if threat.Value < config.Threats.LowTemperature-10 {
			return "Krytyczne"
		} else if threat.Value < config.Threats.LowTemperature-5 {
			return "Wysokie"
		}
		return "Średnie"
	case "Silny wiatr":
		if threat.Value > config.Threats.HighWindSpeed+20 {
			return "Krytyczne"
		} else if threat.Value > config.Threats.HighWindSpeed+10 {
			return "Wysokie"
		}
		return "Średnie"
	case "Intensywne opady":
		if threat.Value > config.Threats.HighPrecipitation*2 {
			return "Krytyczne"
		} else if threat.Value > config.Threats.HighPrecipitation*1.5 {
			return "Wysokie"
		}
		return "Średnie"
	}
	return "Nieznane"
}
