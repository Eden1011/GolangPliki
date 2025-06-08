package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func DisplayCurrentWeather(weather *CurrentWeatherResponse) {
	fmt.Printf("AKTUALNA POGODA\n")
	fmt.Printf("Lokalizacja: %.2f, %.2f\n", weather.Latitude, weather.Longitude)
	fmt.Printf("Czas: %s\n\n", weather.Current.GetFormattedTime())

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Parametr", "Wartość", "Jednostka"})
	table.SetBorder(true)
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")

	data := [][]string{
		{" Temperatura", fmt.Sprintf("%.1f", weather.Current.Temperature), "°C"},
		{" Temperatura odczuwalna", fmt.Sprintf("%.1f", weather.Current.ApparentTemperature), "°C"},
		{" Wilgotność", fmt.Sprintf("%d", weather.Current.RelativeHumidity), "%"},
		{" Opady", fmt.Sprintf("%.1f", weather.Current.Precipitation), "mm"},
		{" Prędkość wiatru", fmt.Sprintf("%.1f", weather.Current.WindSpeed), "km/h"},
		{" Kierunek wiatru", weather.Current.GetWindDirectionText(), fmt.Sprintf("%.0f°", weather.Current.WindDirection)},
		{" Zachmurzenie", fmt.Sprintf("%d", weather.Current.CloudCover), "%"},
		{" Ciśnienie", fmt.Sprintf("%.0f", weather.Current.Pressure), "hPa"},
	}

	for _, row := range data {
		table.Append(row)
	}

	table.Render()
	fmt.Printf("\n%s Ogólna charakterystyka: %s\n", weather.Current.GetWeatherIcon(), getWeatherDescription(weather.Current.Temperature, weather.Current.Precipitation, weather.Current.CloudCover))
}

func DisplayForecast(forecast *ForecastResponse) {
	fmt.Printf("PROGNOZA POGODY\n")
	fmt.Printf("Lokalizacja: %.2f, %.2f\n\n", forecast.Latitude, forecast.Longitude)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Data",
		"Temp. Min",
		"Temp. Max",
		"Odczuwalna Min",
		"Odczuwalna Max",
		"Opady",
		"Wiatr Max",
		"UV Index",
	})
	table.SetBorder(true)

	for i := 0; i < len(forecast.Daily.Time); i++ {
		row := []string{
			forecast.Daily.Time[i],
			fmt.Sprintf("%.1f°C", forecast.Daily.TemperatureMin[i]),
			fmt.Sprintf("%.1f°C", forecast.Daily.TemperatureMax[i]),
			fmt.Sprintf("%.1f°C", forecast.Daily.ApparentTempMin[i]),
			fmt.Sprintf("%.1f°C", forecast.Daily.ApparentTempMax[i]),
			fmt.Sprintf("%.1f mm", forecast.Daily.Precipitation[i]),
			fmt.Sprintf("%.1f km/h", forecast.Daily.WindSpeedMax[i]),
			fmt.Sprintf("%.1f", forecast.Daily.UVIndexMax[i]),
		}
		table.Append(row)
	}

	table.Render()
}

func DisplayHistory(history *HistoryResponse) {
	fmt.Printf("DANE HISTORYCZNE\n")
	fmt.Printf("Lokalizacja: %.2f, %.2f\n\n", history.Latitude, history.Longitude)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Data",
		"Temp. Min",
		"Temp. Max",
		"Odczuwalna Min",
		"Odczuwalna Max",
		"Opady",
		"Wiatr Max",
		"UV Index",
	})
	table.SetBorder(true)

	for i := 0; i < len(history.Daily.Time); i++ {
		row := []string{
			history.Daily.Time[i],
			fmt.Sprintf("%.1f°C", history.Daily.TemperatureMin[i]),
			fmt.Sprintf("%.1f°C", history.Daily.TemperatureMax[i]),
			fmt.Sprintf("%.1f°C", history.Daily.ApparentTempMin[i]),
			fmt.Sprintf("%.1f°C", history.Daily.ApparentTempMax[i]),
			fmt.Sprintf("%.1f mm", history.Daily.Precipitation[i]),
			fmt.Sprintf("%.1f km/h", history.Daily.WindSpeedMax[i]),
			fmt.Sprintf("%.1f", history.Daily.UVIndexMax[i]),
		}
		table.Append(row)
	}

	table.Render()

	if len(history.Daily.Time) > 0 {
		fmt.Println("\nPODSUMOWANIE STATYSTYCZNE:")

		statTable := tablewriter.NewWriter(os.Stdout)
		statTable.SetHeader([]string{"Parametr", "Średnia", "Minimum", "Maksimum"})
		statTable.SetBorder(true)

		tempMinAvg, tempMinMin, tempMinMax := calculateStats(history.Daily.TemperatureMin)
		tempMaxAvg, tempMaxMin, tempMaxMax := calculateStats(history.Daily.TemperatureMax)
		precipAvg, precipMin, precipMax := calculateStats(history.Daily.Precipitation)
		windAvg, windMin, windMax := calculateStats(history.Daily.WindSpeedMax)

		statData := [][]string{
			{"Temperatura Min", fmt.Sprintf("%.1f°C", tempMinAvg), fmt.Sprintf("%.1f°C", tempMinMin), fmt.Sprintf("%.1f°C", tempMinMax)},
			{"Temperatura Max", fmt.Sprintf("%.1f°C", tempMaxAvg), fmt.Sprintf("%.1f°C", tempMaxMin), fmt.Sprintf("%.1f°C", tempMaxMax)},
			{"Opady", fmt.Sprintf("%.1f mm", precipAvg), fmt.Sprintf("%.1f mm", precipMin), fmt.Sprintf("%.1f mm", precipMax)},
			{"Wiatr Max", fmt.Sprintf("%.1f km/h", windAvg), fmt.Sprintf("%.1f km/h", windMin), fmt.Sprintf("%.1f km/h", windMax)},
		}

		for _, row := range statData {
			statTable.Append(row)
		}

		statTable.Render()
	}
}

func getWeatherDescription(temp, precipitation float64, clouds int) string {
	if precipitation > 1.0 {
		return "Deszczowo"
	}
	if clouds > 80 {
		return "Bardzo pochmurno"
	}
	if clouds > 50 {
		return "Pochmurno"
	}
	if temp > 25 {
		return "Słonecznie i ciepło"
	}
	if temp < 0 {
		return "Mroźnie"
	}
	if temp < 10 {
		return "Chłodno"
	}
	return "Przyjemnie"
}

func calculateStats(data []float64) (avg, min, max float64) {
	if len(data) == 0 {
		return 0, 0, 0
	}

	sum := 0.0
	min = data[0]
	max = data[0]

	for _, value := range data {
		sum += value
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	avg = sum / float64(len(data))
	return avg, min, max
}
