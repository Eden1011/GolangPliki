package indicators

import "stock/internal/data"

type Indicator interface {
	Name() string
	Calculate(data []data.StockData) (map[string][]float64, error)
	Description() string
}
