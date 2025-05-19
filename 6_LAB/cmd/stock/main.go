package main

import (
	"flag"
	"fmt"
	"os"
	"stock/internal/data"
	"stock/internal/indicators"
	"stock/internal/indicators/momentum"
	trend "stock/internal/indicators/trends"
	"stock/internal/indicators/volatility"
)

func main() {
	dataFile := flag.String("data", "data.csv", "Ścieżka do pliku CSV z danymi giełdowymi")
	trendPeriod := flag.Int("tperiod", 14, "Okres dla wskaźnika trendu")
	momentumPeriod := flag.Int("mperiod", 14, "Okres dla wskaźnika impetu")
	volatilityPeriod := flag.Int("vperiod", 20, "Okres dla wskaźnika zmienności")
	deviations := flag.Float64("dev", 2.0, "Odchylenia standardowe dla Bollinger Bands")
	flag.Parse()

	fmt.Printf("Wczytywanie danych z %s...\n", *dataFile)
	stockData, err := data.LoadCSV(*dataFile)
	if err != nil {
		fmt.Printf("Błąd wczytywania danych: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Wczytano %d punktów danych\n", len(stockData))

	indList := []indicators.Indicator{
		trend.NewSMA(*trendPeriod),
		momentum.NewRSI(*momentumPeriod),
		volatility.NewBollingerBands(*volatilityPeriod, *deviations),
	}

	for _, ind := range indList {
		fmt.Printf("\n%s:\n", ind.Name())
		fmt.Printf("%s\n", ind.Description())

		results, err := ind.Calculate(stockData)
		if err != nil {
			fmt.Printf("Błąd obliczania %s: %v\n", ind.Name(), err)
			continue
		}

		printCount := 5
		if len(stockData) < printCount {
			printCount = len(stockData)
		}

		if ind.Name() == "Bollinger Bands" {
			fmt.Printf("%-12s%-12s%-12s%-12s\n", "Data", "Middle", "Upper", "Lower")

			for i := len(stockData) - printCount; i < len(stockData); i++ {
				fmt.Printf("%-12s", stockData[i].Date.Format("2006-01-02"))

				middleValues := results["Middle"]
				upperValues := results["Upper"]
				lowerValues := results["Lower"]

				if i < len(middleValues) && middleValues[i] != 0 {
					fmt.Printf("%-12.4f", middleValues[i])
				} else {
					fmt.Printf("%-12s", "N/A")
				}

				if i < len(upperValues) && upperValues[i] != 0 {
					fmt.Printf("%-12.4f", upperValues[i])
				} else {
					fmt.Printf("%-12s", "N/A")
				}

				if i < len(lowerValues) && lowerValues[i] != 0 {
					fmt.Printf("%-12.4f", lowerValues[i])
				} else {
					fmt.Printf("%-12s", "N/A")
				}

				fmt.Println()
			}
		} else {
			fmt.Printf("%-12s", "Data")
			for key := range results {
				fmt.Printf("%-12s", key)
			}
			fmt.Println()

			for i := len(stockData) - printCount; i < len(stockData); i++ {
				fmt.Printf("%-12s", stockData[i].Date.Format("2006-01-02"))
				for _, values := range results {
					if i < len(values) && values[i] != 0 {
						fmt.Printf("%-12.4f", values[i])
					} else {
						fmt.Printf("%-12s", "N/A")
					}
				}
				fmt.Println()
			}
		}
	}
}
