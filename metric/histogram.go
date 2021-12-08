package metric

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Histogram struct {
	Title       string
	SeriesName  string
	Columns     []string
	XAxisLabel  string
	Values      []float64
	ChartHeight int
	ChartWidth  int
	PointsWidth float64
	ChartType   string
	BarColor    int
}

func NewHistogram() Histogram {
	c := Histogram{
		BarColor:    2,
		ChartHeight: 10,
		ChartWidth:  10,
		PointsWidth: 20,
	}
	return c
}

func (h Histogram) Draw(outputDir string) (string, error) {
	p := plot.New()
	p.Title.Text = h.Title
	p.X.Label.Text = h.XAxisLabel
	p.NominalX(h.Columns...)
	filename := fmt.Sprintf("%s/%s.png", outputDir, h.SeriesName)
	dataValues := h.generateBarItems()
	chart, err := plotter.NewBarChart(dataValues, vg.Points(h.PointsWidth))
	if err != nil {
		return "", err
	}
	p.Add(chart)
	chart.Color = plotutil.Color(h.BarColor)
	chart.LineStyle.Width = vg.Length(0)
	chartWidth := vg.Length(h.ChartWidth)
	ChartHeight := vg.Length(h.ChartHeight)
	if err := p.Save(chartWidth*vg.Inch, ChartHeight*vg.Inch, filename); err != nil {
		return "", err
	}
	return filename, nil
}

func (h Histogram) generateBarItems() plotter.Values {
	var prevGroup float64
	var items plotter.Values
	for i, bucket := range h.Values {
		var newGroupValue float64
		if i > 0 {
			newGroupValue = bucket - prevGroup
		} else {
			newGroupValue = bucket
		}
		prevGroup = bucket
		items = append(items, newGroupValue)
	}
	h.Values = items
	return items
}
