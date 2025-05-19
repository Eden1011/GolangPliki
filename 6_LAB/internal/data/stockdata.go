package data

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type StockData struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func LoadCSV(filename string) ([]StockData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, errors.New("plik ma niewystarczającą ilość danych")
	}

	data := make([]StockData, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 6 {
			continue
		}

		date, err := time.Parse("01/02/2006", record[0])
		if err != nil {
			return nil, fmt.Errorf("błąd parsowania daty w linii %d: %v", i+1, err)
		}

		open, err := parsePrice(record[3])
		if err != nil {
			return nil, fmt.Errorf("błąd parsowania ceny otwarcia w linii %d: %v", i+1, err)
		}

		high, err := parsePrice(record[4])
		if err != nil {
			return nil, fmt.Errorf("błąd parsowania najwyższej ceny w linii %d: %v", i+1, err)
		}

		low, err := parsePrice(record[5])
		if err != nil {
			return nil, fmt.Errorf("błąd parsowania najniższej ceny w linii %d: %v", i+1, err)
		}

		close, err := parsePrice(record[1])
		if err != nil {
			return nil, fmt.Errorf("błąd parsowania ceny zamknięcia w linii %d: %v", i+1, err)
		}

		volume, err := strconv.ParseInt(strings.ReplaceAll(record[2], ",", ""), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("błąd parsowania wolumenu w linii %d: %v", i+1, err)
		}

		stockData := StockData{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		}
		data = append(data, stockData)
	}

	return data, nil
}

func parsePrice(priceStr string) (float64, error) {
	cleanPrice := strings.ReplaceAll(priceStr, "$", "")
	cleanPrice = strings.ReplaceAll(cleanPrice, ",", "")
	return strconv.ParseFloat(cleanPrice, 64)
}
