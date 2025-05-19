package momentum

import (
	"errors"
	"stock/internal/data"
	"stock/internal/indicators"
)

type RSI struct {
	period int
}

func NewRSI(period int) *RSI {
	return &RSI{period: period}
}

func (r *RSI) Name() string {
	return "Relative Strength Index"
}

func (r *RSI) Calculate(data []data.StockData) (map[string][]float64, error) {
	if len(data) <= r.period {
		return nil, errors.New("niewystarczająca ilość danych dla obliczenia RSI")
	}

	changes := make([]float64, len(data))
	for i := 1; i < len(data); i++ {
		changes[i] = data[i].Close - data[i-1].Close
	}

	rsiValues := make([]float64, len(data))

	for i := 0; i < r.period; i++ {
		rsiValues[i] = 0
	}

	var sumGain, sumLoss float64
	for i := 1; i <= r.period; i++ {
		change := changes[i]
		if change > 0 {
			sumGain += change
		} else {
			sumLoss -= change
		}
	}

	if sumLoss == 0 {
		rsiValues[r.period] = 100
	} else {
		rs := sumGain / sumLoss
		rsiValues[r.period] = 100 - (100 / (1 + rs))
	}

	for i := r.period + 1; i < len(data); i++ {
		change := changes[i]

		var gain, loss float64
		if change > 0 {
			gain = change
			loss = 0
		} else {
			gain = 0
			loss = -change
		}

		sumGain = ((sumGain * float64(r.period-1)) + gain) / float64(r.period)
		sumLoss = ((sumLoss * float64(r.period-1)) + loss) / float64(r.period)

		if sumLoss == 0 {
			rsiValues[i] = 100
		} else {
			rs := sumGain / sumLoss
			rsiValues[i] = 100 - (100 / (1 + rs))
		}
	}

	return map[string][]float64{"RSI": rsiValues}, nil
}

func (r *RSI) Description() string {
	return "Relative Strength Index (RSI) to oscylator impetu, który mierzy prędkość i zmianę ruchów cenowych. RSI oscyluje między 0 a 100."
}

var _ indicators.Indicator = (*RSI)(nil)
