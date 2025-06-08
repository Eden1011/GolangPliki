package main

import (
	"fmt"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func GenerateWeatherChart(forecast *ForecastResponse, filename string) error {
	p := plot.New()
	p.Title.Text = "Prognoza Temperatury"
	p.X.Label.Text = "Data"
	p.Y.Label.Text = "Temperatura (°C)"

	tempMinPoints := make(plotter.XYs, len(forecast.Daily.Time))
	tempMaxPoints := make(plotter.XYs, len(forecast.Daily.Time))

	for i := 0; i < len(forecast.Daily.Time); i++ {

		date, err := time.Parse("2006-01-02", forecast.Daily.Time[i])
		if err != nil {
			return fmt.Errorf("błąd parsowania daty: %w", err)
		}

		x := float64(date.Unix())
		tempMinPoints[i].X = x
		tempMinPoints[i].Y = forecast.Daily.TemperatureMin[i]
		tempMaxPoints[i].X = x
		tempMaxPoints[i].Y = forecast.Daily.TemperatureMax[i]
	}

	minLine, err := plotter.NewLine(tempMinPoints)
	if err != nil {
		return fmt.Errorf("błąd tworzenia linii minimum: %w", err)
	}
	minLine.LineStyle.Color = plotutil.Color(1)
	minLine.LineStyle.Width = vg.Points(2)

	maxLine, err := plotter.NewLine(tempMaxPoints)
	if err != nil {
		return fmt.Errorf("błąd tworzenia linii maksimum: %w", err)
	}
	maxLine.LineStyle.Color = plotutil.Color(0)
	maxLine.LineStyle.Width = vg.Points(2)

	minPoints, err := plotter.NewScatter(tempMinPoints)
	if err != nil {
		return fmt.Errorf("błąd tworzenia punktów minimum: %w", err)
	}
	minPoints.GlyphStyle.Color = plotutil.Color(1)
	minPoints.GlyphStyle.Shape = draw.CircleGlyph{}

	maxPoints, err := plotter.NewScatter(tempMaxPoints)
	if err != nil {
		return fmt.Errorf("błąd tworzenia punktów maksimum: %w", err)
	}
	maxPoints.GlyphStyle.Color = plotutil.Color(0)
	maxPoints.GlyphStyle.Shape = draw.CircleGlyph{}

	p.Add(minLine, minPoints, maxLine, maxPoints)

	p.Legend.Add("Temperatura Min", minLine, minPoints)
	p.Legend.Add("Temperatura Max", maxLine, maxPoints)
	p.Legend.Top = true

	p.X.Tick.Marker = &DateTicker{}

	if err := p.Save(12*vg.Inch, 8*vg.Inch, filename); err != nil {
		return fmt.Errorf("błąd zapisywania wykresu: %w", err)
	}

	return nil
}

func GenerateHistoryChart(history *HistoryResponse, filename string) error {
	p := plot.New()
	p.Title.Text = "Historia Temperatury"
	p.X.Label.Text = "Data"
	p.Y.Label.Text = "Temperatura (°C)"

	tempMinPoints := make(plotter.XYs, len(history.Daily.Time))
	tempMaxPoints := make(plotter.XYs, len(history.Daily.Time))

	for i := 0; i < len(history.Daily.Time); i++ {

		date, err := time.Parse("2006-01-02", history.Daily.Time[i])
		if err != nil {
			return fmt.Errorf("błąd parsowania daty: %w", err)
		}

		x := float64(date.Unix())
		tempMinPoints[i].X = x
		tempMinPoints[i].Y = history.Daily.TemperatureMin[i]
		tempMaxPoints[i].X = x
		tempMaxPoints[i].Y = history.Daily.TemperatureMax[i]
	}

	minLine, err := plotter.NewLine(tempMinPoints)
	if err != nil {
		return fmt.Errorf("błąd tworzenia linii minimum: %w", err)
	}
	minLine.LineStyle.Color = plotutil.Color(1)
	minLine.LineStyle.Width = vg.Points(2)

	maxLine, err := plotter.NewLine(tempMaxPoints)
	if err != nil {
		return fmt.Errorf("błąd tworzenia linii maksimum: %w", err)
	}
	maxLine.LineStyle.Color = plotutil.Color(0)
	maxLine.LineStyle.Width = vg.Points(2)

	minPoints, err := plotter.NewScatter(tempMinPoints)
	if err != nil {
		return fmt.Errorf("błąd tworzenia punktów minimum: %w", err)
	}
	minPoints.GlyphStyle.Color = plotutil.Color(1)
	minPoints.GlyphStyle.Shape = draw.CircleGlyph{}

	maxPoints, err := plotter.NewScatter(tempMaxPoints)
	if err != nil {
		return fmt.Errorf("błąd tworzenia punktów maksimum: %w", err)
	}
	maxPoints.GlyphStyle.Color = plotutil.Color(0)
	maxPoints.GlyphStyle.Shape = draw.CircleGlyph{}

	p.Add(minLine, minPoints, maxLine, maxPoints)

	p.Legend.Add("Temperatura Min", minLine, minPoints)
	p.Legend.Add("Temperatura Max", maxLine, maxPoints)
	p.Legend.Top = true

	p.X.Tick.Marker = &DateTicker{}

	if err := p.Save(12*vg.Inch, 8*vg.Inch, filename); err != nil {
		return fmt.Errorf("błąd zapisywania wykresu: %w", err)
	}

	return nil
}

type DateTicker struct{}

func (dt *DateTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick

	start := time.Unix(int64(min), 0)
	end := time.Unix(int64(max), 0)

	duration := end.Sub(start)
	var interval time.Duration
	var format string

	switch {
	case duration <= 7*24*time.Hour:
		interval = 24 * time.Hour
		format = "01-02"
	case duration <= 30*24*time.Hour:
		interval = 3 * 24 * time.Hour
		format = "01-02"
	default:
		interval = 7 * 24 * time.Hour
		format = "01-02"
	}

	for t := start.Truncate(interval); t.Before(end.Add(interval)); t = t.Add(interval) {
		if t.After(start.Add(-interval)) && t.Before(end.Add(interval)) {
			ticks = append(ticks, plot.Tick{
				Value: float64(t.Unix()),
				Label: t.Format(format),
			})
		}
	}

	return ticks
}
