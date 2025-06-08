package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var city string
	var days int
	var dates string

	flag.StringVar(&city, "miasto", "", "Nazwa miasta")
	flag.IntVar(&days, "dni", 7, "Liczba dni prognozy (domyślnie 7)")
	flag.StringVar(&dates, "daty", "", "Daty w formacie YYYY-MM-DD,YYYY-MM-DD")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Użycie:")
		fmt.Println("  go run . aktualna <miasto>")
		fmt.Println("  go run . prognoza <miasto> [dni]")
		fmt.Println("  go run . historia <miasto> [daty]")
		os.Exit(1)
	}

	command := args[0]
	cityArg := args[1]

	if cityArg != "" {
		city = cityArg
	}

	if city == "" {
		log.Fatal("Miasto jest wymagane")
	}

	if len(args) > 2 {
		switch command {
		case "prognoza":
			if d, err := strconv.Atoi(args[2]); err == nil {
				days = d
			}
		case "historia":
			dates = args[2]
		}
	}

	config, err := LoadConfig("config.json")
	if err != nil {
		log.Printf("Nie można załadować konfiguracji: %v. Używam domyślnych wartości.", err)
		config = GetDefaultConfig()
	}

	coord, err := GetCityCoordinates(city)
	if err != nil {
		log.Fatalf("Nie można znaleźć miasta %s: %v", city, err)
	}

	fmt.Printf(" Dane pogodowe dla: %s (%.2f, %.2f)\n\n", city, coord.Latitude, coord.Longitude)

	switch command {
	case "aktualna":
		handleCurrentWeather(coord)
	case "prognoza":
		handleForecast(coord, days, config)
	case "historia":
		handleHistory(coord, dates)
	default:
		fmt.Printf("Nieznana komenda: %s\n", command)
		os.Exit(1)
	}
}

func handleCurrentWeather(coord *Coordinates) {
	weather, err := GetCurrentWeather(coord.Latitude, coord.Longitude)
	if err != nil {
		log.Fatalf("Błąd pobierania aktualnej pogody: %v", err)
	}

	DisplayCurrentWeather(weather)
}

func handleForecast(coord *Coordinates, days int, config *Config) {
	if days > 16 {
		days = 16
	}

	forecast, err := GetWeatherForecast(coord.Latitude, coord.Longitude, days)
	if err != nil {
		log.Fatalf("Błąd pobierania prognozy: %v", err)
	}

	DisplayForecast(forecast)

	err = GenerateWeatherChart(forecast, "prognoza_temperatury.png")
	if err != nil {
		log.Printf("Błąd generowania wykresu: %v", err)
	} else {
		fmt.Println(" Wykres zapisano do pliku: prognoza_temperatury.png")
	}

	threats := AnalyzeWeatherThreats(forecast, config)
	if len(threats) > 0 {
		fmt.Println("\n  OSTRZEŻENIA POGODOWE:")
		for _, threat := range threats {
			fmt.Printf("%s - %s: %s\n", threat.Date, threat.Type, threat.Description)
		}
	} else {
		fmt.Println("\nBrak ostrzeżeń pogodowych")
	}
}

func handleHistory(coord *Coordinates, dates string) {
	if dates == "" {

		endDate := time.Now().AddDate(0, 0, -1)
		startDate := endDate.AddDate(0, 0, -6)
		dates = fmt.Sprintf("%s,%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	dateParts := strings.Split(dates, ",")
	if len(dateParts) != 2 {
		log.Fatal("Daty powinny być w formacie: YYYY-MM-DD,YYYY-MM-DD")
	}

	startDate := dateParts[0]
	endDate := dateParts[1]

	history, err := GetWeatherHistory(coord.Latitude, coord.Longitude, startDate, endDate)
	if err != nil {
		log.Fatalf("Błąd pobierania danych historycznych: %v", err)
	}

	DisplayHistory(history)

	err = GenerateHistoryChart(history, "historia_temperatury.png")
	if err != nil {
		log.Printf("Błąd generowania wykresu: %v", err)
	} else {
		fmt.Println("📊 Wykres zapisano do pliku: historia_temperatury.png")
	}
}
