package volatility

import (
	"errors"
	"math"
	"stock/internal/data"
	"stock/internal/indicators"
)

type BollingerBands struct {
	period     int
	deviations float64
}

func NewBollingerBands(period int, deviations float64) *BollingerBands {
	return &BollingerBands{
		period:     period,
		deviations: deviations,
	}
}

func (b *BollingerBands) Name() string {
	return "Bollinger Bands"
}

func (b *BollingerBands) Calculate(data []data.StockData) (map[string][]float64, error) {
	if len(data) < b.period {
		return nil, errors.New("niewystarczająca ilość danych dla obliczenia Bollinger Bands")
	}

	middle := make([]float64, len(data))
	upper := make([]float64, len(data))
	lower := make([]float64, len(data))

	for i := 0; i < b.period-1; i++ {
		middle[i] = 0
		upper[i] = 0
		lower[i] = 0
	}

	for i := b.period - 1; i < len(data); i++ {
		sum := 0.0
		for j := 0; j < b.period; j++ {
			sum += data[i-j].Close
		}
		middle[i] = sum / float64(b.period)

		variance := 0.0
		for j := 0; j < b.period; j++ {
			diff := data[i-j].Close - middle[i]
			variance += diff * diff
		}
		stdDev := math.Sqrt(variance / float64(b.period))

		upper[i] = middle[i] + (b.deviations * stdDev)
		lower[i] = middle[i] - (b.deviations * stdDev)
	}

	return map[string][]float64{
		"Middle": middle,
		"Upper":  upper,
		"Lower":  lower,
	}, nil
}

func (b *BollingerBands) Description() string {
	return "Bollinger Bands to pasma zmienności umieszczone powyżej i poniżej średniej ruchomej. Zmienność opiera się na odchyleniu standardowym."
}

var _ indicators.Indicator = (*BollingerBands)(nil)
