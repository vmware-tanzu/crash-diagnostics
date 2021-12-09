package metric

import (
	"fmt"
	dto "github.com/prometheus/client_model/go"
)

const (
	EtcdDefaultClientCert = "/etc/kubernetes/pki/apiserver-etcd-client.crt"
	EtcdDefaultKeyFile    = "/etc/kubernetes/pki/apiserver-etcd-client.key"
	EtcdDefaultEndpoint   = "https://localhost:2379/metrics"
	EtcdCurlCmd           = "sudo curl -sk"
)

func etcdKnownMetrics() map[string]struct{} {
	return map[string]struct{}{
		"etcd_disk_backend_commit_duration_seconds":                        struct{}{},
		"etcd_network_peer_round_trip_time_seconds":                        struct{}{},
		"etcd_debugging_disk_backend_commit_rebalance_duration_seconds":    struct{}{},
		"etcd_debugging_disk_backend_commit_spill_duration_seconds":        struct{}{},
		"etcd_debugging_disk_backend_commit_write_duration_seconds":        struct{}{},
		"etcd_debugging_lease_ttl_total":                                   struct{}{},
		"etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds":    struct{}{},
		"etcd_debugging_mvcc_db_compaction_total_duration_milliseconds":    struct{}{},
		"etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds": struct{}{},
		"etcd_debugging_snap_save_marshalling_duration_seconds":            struct{}{},
		"etcd_debugging_snap_save_total_duration_seconds":                  struct{}{},
		"etcd_disk_backend_defrag_duration_seconds":                        struct{}{},
		"etcd_disk_backend_snapshot_duration_seconds":                      struct{}{},
		"etcd_disk_wal_fsync_duration_seconds":                             struct{}{},
		"etcd_mvcc_hash_duration_seconds":                                  struct{}{},
		"etcd_mvcc_hash_rev_duration_seconds":                              struct{}{},
		"etcd_snap_db_fsync_duration_seconds":                              struct{}{},
		"etcd_snap_db_save_total_duration_seconds":                         struct{}{},
		"etcd_snap_fsync_duration_seconds":                                 struct{}{},
	}

}

type EtcdMetricClient struct {
	Name             string
	SupportedMetrics map[string]struct{}
	Command          string
	ServerKey        string
	ClientCert       string
	Endpoint         string
	WorkDir          string
}

func (m EtcdMetricClient) generateHistogram(metricFamily *dto.MetricFamily) ([]string, error) {
	var files []string
	for _, metric := range metricFamily.GetMetric() {
		chart := NewHistogram()
		switch metricFamily.GetName() {
		case "etcd_network_peer_round_trip_time_seconds":
			// one graph per node round trip
			node := metric.GetLabel()
			seriesName := fmt.Sprintf("%s_%s", node[0].GetName(), node[0].GetValue())
			chart.SeriesName = fmt.Sprintf("%s_%s", metricFamily.GetName(), seriesName)
			chart.BarColor = 3
			chart.Title = fmt.Sprintf("%s to %s", metricFamily.GetHelp(), node[0].GetValue())
			chart.XAxisLabel = chart.SeriesName
		default:
			// just one graph
			chart.Title = fmt.Sprintf("%s", metricFamily.GetHelp())
			chart.SeriesName = metricFamily.GetName()
		}
		chart.XAxisLabel = metricFamily.GetName()
		for _, bucket := range metric.Histogram.GetBucket() {
			chart.Columns = append(chart.Columns, fmt.Sprintf("%v", bucket.GetUpperBound()))
			chart.Values = append(chart.Values, float64(bucket.GetCumulativeCount()))
		}
		chartFile, err := chart.Draw(m.WorkDir)
		if err != nil {
			return nil, err
		}
		files = append(files, chartFile)
	}
	return files, nil
}

func (m EtcdMetricClient) Plot(filename string, filter []string) ([]string, error) {
	// get prom
	var charts []string
	promMetrics, err := ParseMF(filename)
	if err != nil {
		return nil, err
	}
	for _, metric := range promMetrics {
		var chartFiles []string
		if m.IsKnownMetric(metric.GetName()) {
			switch metric.GetType() {
			case dto.MetricType_HISTOGRAM:
				chartFiles, err = m.generateHistogram(metric)
				if err != nil {
					return nil, err
				}
			}
			charts = append(charts, chartFiles...)
		} else {
			fmt.Printf("Skipping Unsupported Metric: %s\n", metric.GetName())
		}
	}
	return charts, nil
}

func (m EtcdMetricClient) GetCommandOutput() string {
	if m.ServerKey == "" {
		m.ServerKey = EtcdDefaultKeyFile
	}
	if m.ClientCert == "" {
		m.ClientCert = EtcdDefaultClientCert
	}
	if m.Endpoint == "" {
		m.Endpoint = EtcdDefaultEndpoint
	}
	m.Command = fmt.Sprintf("%s --cert %s --key %s %s", EtcdCurlCmd, m.ClientCert, m.ServerKey, m.Endpoint)
	return m.Command
}

func (m EtcdMetricClient) IsKnownMetric(metric string) bool {
	_, ok := m.SupportedMetrics[metric]
	return ok
}
