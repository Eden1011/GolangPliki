package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

type Exchange struct {
	Period   string
	Currency string
	Rate     float64
}

func LoadExchangeRates(filepath string) []Exchange {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Encountered fatal error during file reading: %v", err)
	}
	var exchangeRates []Exchange
	for i, v := range records {
		if i == 0 || v[2] == "" {
			continue
		}
		rate, err := strconv.ParseFloat(v[2], 64)
		if err != nil {
			log.Printf("Warning: Could not parse rate for row %d: %v", i, err)
			continue
		}
		exchangeRates = append(exchangeRates, Exchange{
			Period:   v[0],
			Currency: v[1],
			Rate:     rate,
		})
	}
	return exchangeRates
}

func SortExchangeRatesByRate(data []Exchange, ascending bool) []Exchange {
	result := make([]Exchange, len(data))
	copy(result, data)
	if ascending {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Rate < result[j].Rate
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Rate > result[j].Rate
		})
	}
	return result
}

func SortExchangeRatesByCurrency(data []Exchange, ascending bool) []Exchange {
	result := make([]Exchange, len(data))
	copy(result, data)
	if ascending {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Currency < result[j].Currency
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Currency > result[j].Currency
		})
	}
	return result
}

func GetCurrencyStats(data []Exchange, currency string) (float64, Exchange, Exchange) {
	var filteredData []Exchange

	for _, exchange := range data {
		if exchange.Currency == currency {
			filteredData = append(filteredData, exchange)
		}
	}

	if len(filteredData) == 0 {
		log.Fatalf("Currency %s not found in the data", currency)
	}

	var sum float64
	for _, exchange := range filteredData {
		sum += exchange.Rate
	}
	average := sum / float64(len(filteredData))

	highest := Exchange{Rate: -math.MaxFloat64}
	lowest := Exchange{Rate: math.MaxFloat64}

	for _, exchange := range filteredData {
		if exchange.Rate > highest.Rate {
			highest = exchange
		}
		if exchange.Rate < lowest.Rate {
			lowest = exchange
		}
	}

	return average, highest, lowest
}

func main() {
	exchangeRates := LoadExchangeRates("./euro-exchange-rates.csv")
	average, highest, lowest := GetCurrencyStats(exchangeRates, "USD")
	fmt.Printf("Average Rate: %.4f\n", average)
	fmt.Printf("Highest Rate: %.4f (on %s)\n", highest.Rate, highest.Period)
	fmt.Printf("Lowest Rate: %.4f (on %s)\n", lowest.Rate, lowest.Period)

	// fmt.Println("\nSorted by Rate (Ascending):")
	// rateAscending := SortExchangeRatesByRate(exchangeRates, true)
	// fmt.Println(rateAscending)
	//
	// fmt.Println("\nSorted by Rate (Descending):")
	// rateDescending := SortExchangeRatesByRate(exchangeRates, false)
	// fmt.Println(rateDescending)
}
