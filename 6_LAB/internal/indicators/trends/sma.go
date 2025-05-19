package trend

import (
	"errors"
	"stock/internal/data"
	"stock/internal/indicators"
)

type SMA struct {
	period int
}

func NewSMA(period int) *SMA {
	return &SMA{period: period}
}

func (s *SMA) Name() string {
	return "Simple Moving Average"
}

func (s *SMA) Calculate(data []data.StockData) (map[string][]float64, error) {
	if len(data) < s.period {
		return nil, errors.New("niewystarczająca ilość danych dla obliczenia SMA")
	}

	result := make([]float64, len(data))

	for i := 0; i < s.period-1; i++ {
		result[i] = 0
	}

	for i := s.period - 1; i < len(data); i++ {
		sum := 0.0
		for j := 0; j < s.period; j++ {
			sum += data[i-j].Close
		}
		result[i] = sum / float64(s.period)
	}

	return map[string][]float64{"SMA": result}, nil
}

func (s *SMA) Description() string {
	return "Simple Moving Average (SMA) to średnia arytmetyczna obliczana przez dodanie ostatnich cen zamknięcia i podzielenie przez liczbę okresów."
}

var _ indicators.Indicator = (*SMA)(nil)
