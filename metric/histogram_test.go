package metric

import (
	"os"
	"reflect"
	"testing"
)

func TestHistogram_Draw(t *testing.T) {
	chart := NewHistogram()
	chart.Title = "some metric"
	chart.SeriesName = "node1"
	chart.Columns = []string{
		"Group A",
		"Group B",
		"Group C",
	}
	chart.ChartWidth = 8
	chart.ChartHeight = 8
	chart.PointsWidth = 20
	chart.XAxisLabel = "datagroups"
	chart.Values = []float64{
		10,
		15,
		25,
	}
	chart.ChartWidth = 8
	chart.ChartHeight = 8
	path, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	filename, err := chart.Draw(path)
	if err != nil {
		t.Error(err)
		t.Fail()
	} else {
		t.Log(filename)
	}
}

func TestHistogram_GenerateBarItems(t *testing.T) {
	var chart Histogram
	chart.Values = []float64{
		10,
		15,
		25,
	}
	expected := []float64{10, 5, 10}
	retVal := chart.generateBarItems()
	for i, f := range retVal {
		if f != expected[i] {
			t.Logf("Value for bucket %v is incorrect. Expected: %v Got: %v", i, expected[i], f)
			t.Fail()
		}
	}
}

func TestNewHistogram(t *testing.T) {
	chart := NewHistogram()
	expected := Histogram{
		BarColor:    2,
		ChartHeight: 10,
		ChartWidth:  10,
		PointsWidth: 20,
	}
	if reflect.DeepEqual(chart, expected) {
		t.Logf("Expected %v, %v", expected, chart)
	}
}
