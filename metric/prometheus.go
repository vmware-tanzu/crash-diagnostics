package metric

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"os"
)

////ParseMF uses the prometheus parser library to return a map of metrics from a returned curl from
////an OpenTelemetry metrics endpoint
func ParseMF(path string) (map[string]*dto.MetricFamily, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(reader)

	if err != nil {
		return nil, err
	}
	return mf, nil
}
