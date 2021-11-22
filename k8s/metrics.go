// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"fmt"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"os"
	"strings"
)

const (
	EtcdDiskBackendCommitDurationSeconds                      = "etcd_disk_backend_commit_duration_seconds"
	EtcdNetworkPeerRoundTripTimeSeconds                       = "etcd_network_peer_round_trip_time_seconds"
	EtcdDebuggingDiskBackendCommitRebalanceDurationSeconds    = "etcd_debugging_disk_backend_commit_rebalance_duration_seconds"
	EtcdDebuggingDiskBackendCommitSpillDurationSeconds        = "etcd_debugging_disk_backend_commit_spill_duration_seconds"
	EtcdDebuggingDiskBackendCommitWriteDurationSeconds        = "etcd_debugging_disk_backend_commit_write_duration_seconds"
	EtcdDebuggingLeaseTtlTotal                                = "etcd_debugging_lease_ttl_total"
	EtcdDebuggingMvccDbCompactionPauseDurationMilliseconds    = "etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds"
	EtcdDebuggingMvccDbCompactionTotalDurationMilliseconds    = "etcd_debugging_mvcc_db_compaction_total_duration_milliseconds"
	EtcdDebuggingMvccIndexCompactionPauseDurationMilliseconds = "etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds"
	EtcdDebuggingSnapSaveMarshallingDurationSeconds           = "etcd_debugging_snap_save_marshalling_duration_seconds"
	EtcdDebuggingSnapSaveTotalDurationSeconds                 = "etcd_debugging_snap_save_total_duration_seconds"
	EtcdDiskBackendDefragDurationSeconds                      = "etcd_disk_backend_defrag_duration_seconds"
	EtcdDiskBackendSnapshotDurationSeconds                    = "etcd_disk_backend_snapshot_duration_seconds"
	EtcdDiskWalFsyncDuratioSeconds                            = "etcd_disk_wal_fsync_duration_seconds"
	EtcdMvccHashDurationSeconds                               = "etcd_mvcc_hash_duration_seconds"
	EtcdMvccHashRevDurationSeconds                            = "etcd_mvcc_hash_rev_duration_seconds"
	EtcdSnapDbFsyncDurationSeconds                            = "etcd_snap_db_fsync_duration_seconds"
	EtcdSnapDbSaveTotalDurationSeconds                        = "etcd_snap_db_save_total_duration_seconds"
	EtcdSnapFsyncDurationSeconds                              = "etcd_snap_fsync_duration_seconds"

	EtcdDefaultClientCert      = "/etc/kubernetes/pki/apiserver-etcd-client.crt"
	EtcdDefaultKeyFile         = "/etc/kubernetes/pki/apiserver-etcd-client.key"
	EtcdDefaultEndpoint        = "https://localhost:2379/metrics"
	ApiServerDefaultClientCert = "/etc/kubernetes/pki/apiserver-kubelet-client.crt"
	ApiServerDefaultKeyFile    = "/etc/kubernetes/pki/apiserver-kubelet-client.key"
	ApiServerDefaultEndpoint   = "https://localhost:6443/metrics"
)

var (
	curlCmd = "sudo curl -sk"

	clientCerts = map[string]string{
		EtcdDiskBackendCommitDurationSeconds:                      EtcdDefaultClientCert,
		EtcdNetworkPeerRoundTripTimeSeconds:                       EtcdDefaultClientCert,
		EtcdDebuggingDiskBackendCommitRebalanceDurationSeconds:    EtcdDefaultKeyFile,
		EtcdDebuggingDiskBackendCommitSpillDurationSeconds:        EtcdDefaultKeyFile,
		EtcdDebuggingDiskBackendCommitWriteDurationSeconds:        EtcdDefaultKeyFile,
		EtcdDebuggingLeaseTtlTotal:                                EtcdDefaultKeyFile,
		EtcdDebuggingMvccDbCompactionPauseDurationMilliseconds:    EtcdDefaultKeyFile,
		EtcdDebuggingMvccDbCompactionTotalDurationMilliseconds:    EtcdDefaultKeyFile,
		EtcdDebuggingMvccIndexCompactionPauseDurationMilliseconds: EtcdDefaultKeyFile,
		EtcdDebuggingSnapSaveMarshallingDurationSeconds:           EtcdDefaultKeyFile,
		EtcdDebuggingSnapSaveTotalDurationSeconds:                 EtcdDefaultKeyFile,
		EtcdDiskBackendDefragDurationSeconds:                      EtcdDefaultKeyFile,
		EtcdDiskBackendSnapshotDurationSeconds:                    EtcdDefaultKeyFile,
		EtcdDiskWalFsyncDuratioSeconds:                            EtcdDefaultKeyFile,
		EtcdMvccHashDurationSeconds:                               EtcdDefaultKeyFile,
		EtcdMvccHashRevDurationSeconds:                            EtcdDefaultKeyFile,
		EtcdSnapDbFsyncDurationSeconds:                            EtcdDefaultKeyFile,
		EtcdSnapDbSaveTotalDurationSeconds:                        EtcdDefaultKeyFile,
		EtcdSnapFsyncDurationSeconds:                              EtcdDefaultKeyFile,
	}
	keyfiles = map[string]string{
		EtcdDiskBackendCommitDurationSeconds:                      EtcdDefaultKeyFile,
		EtcdNetworkPeerRoundTripTimeSeconds:                       EtcdDefaultKeyFile,
		EtcdDebuggingDiskBackendCommitRebalanceDurationSeconds:    EtcdDefaultKeyFile,
		EtcdDebuggingDiskBackendCommitSpillDurationSeconds:        EtcdDefaultKeyFile,
		EtcdDebuggingDiskBackendCommitWriteDurationSeconds:        EtcdDefaultKeyFile,
		EtcdDebuggingLeaseTtlTotal:                                EtcdDefaultKeyFile,
		EtcdDebuggingMvccDbCompactionPauseDurationMilliseconds:    EtcdDefaultKeyFile,
		EtcdDebuggingMvccDbCompactionTotalDurationMilliseconds:    EtcdDefaultKeyFile,
		EtcdDebuggingMvccIndexCompactionPauseDurationMilliseconds: EtcdDefaultKeyFile,
		EtcdDebuggingSnapSaveMarshallingDurationSeconds:           EtcdDefaultKeyFile,
		EtcdDebuggingSnapSaveTotalDurationSeconds:                 EtcdDefaultKeyFile,
		EtcdDiskBackendDefragDurationSeconds:                      EtcdDefaultKeyFile,
		EtcdDiskBackendSnapshotDurationSeconds:                    EtcdDefaultKeyFile,
		EtcdDiskWalFsyncDuratioSeconds:                            EtcdDefaultKeyFile,
		EtcdMvccHashDurationSeconds:                               EtcdDefaultKeyFile,
		EtcdMvccHashRevDurationSeconds:                            EtcdDefaultKeyFile,
		EtcdSnapDbFsyncDurationSeconds:                            EtcdDefaultKeyFile,
		EtcdSnapDbSaveTotalDurationSeconds:                        EtcdDefaultKeyFile,
		EtcdSnapFsyncDurationSeconds:                              EtcdDefaultKeyFile,
	}
	endpoints = map[string]string{
		EtcdDiskBackendCommitDurationSeconds:                      EtcdDefaultEndpoint,
		EtcdNetworkPeerRoundTripTimeSeconds:                       EtcdDefaultEndpoint,
		EtcdDebuggingDiskBackendCommitRebalanceDurationSeconds:    EtcdDefaultKeyFile,
		EtcdDebuggingDiskBackendCommitSpillDurationSeconds:        EtcdDefaultKeyFile,
		EtcdDebuggingDiskBackendCommitWriteDurationSeconds:        EtcdDefaultKeyFile,
		EtcdDebuggingLeaseTtlTotal:                                EtcdDefaultKeyFile,
		EtcdDebuggingMvccDbCompactionPauseDurationMilliseconds:    EtcdDefaultKeyFile,
		EtcdDebuggingMvccDbCompactionTotalDurationMilliseconds:    EtcdDefaultKeyFile,
		EtcdDebuggingMvccIndexCompactionPauseDurationMilliseconds: EtcdDefaultKeyFile,
		EtcdDebuggingSnapSaveMarshallingDurationSeconds:           EtcdDefaultKeyFile,
		EtcdDebuggingSnapSaveTotalDurationSeconds:                 EtcdDefaultKeyFile,
		EtcdDiskBackendDefragDurationSeconds:                      EtcdDefaultKeyFile,
		EtcdDiskBackendSnapshotDurationSeconds:                    EtcdDefaultKeyFile,
		EtcdDiskWalFsyncDuratioSeconds:                            EtcdDefaultKeyFile,
		EtcdMvccHashDurationSeconds:                               EtcdDefaultKeyFile,
		EtcdMvccHashRevDurationSeconds:                            EtcdDefaultKeyFile,
		EtcdSnapDbFsyncDurationSeconds:                            EtcdDefaultKeyFile,
		EtcdSnapDbSaveTotalDurationSeconds:                        EtcdDefaultKeyFile,
		EtcdSnapFsyncDurationSeconds:                              EtcdDefaultKeyFile,
	}
)

type Metric struct {
	Name       string
	ClientCert string
	Keyfile    string
	Endpoint   string
}

//GetMetricsCommand generates the command needed to retrieve the metric on the provided nodes
func GetMetricsCommand(m Metric) (string, error) {
	if err := m.Validate(); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s --cert %s --key %s %s", curlCmd, m.ClientCert, m.Keyfile, m.Endpoint), nil
}

//Validate validates the metric object is complete
func (m Metric) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("metric name not found, list of metrics cannot be blank")

	}
	if _, found := endpoints[m.Name]; !found {
		return fmt.Errorf("metric [%s] is not yet supported", m.Name)
	}
	return nil
}

func NewMetric(name, cert, key, endpoint string) Metric {
	var m Metric
	m.Name = name
	if cert == "" {
		m.ClientCert = clientCerts[name]
	} else {
		m.ClientCert = cert
	}
	if key == "" {
		m.Keyfile = keyfiles[name]
	} else {
		m.Keyfile = key
	}
	if key == "" {
		m.Endpoint = endpoints[name]
	} else {
		m.Endpoint = endpoint
	}
	return m
}

//ParseMF uses the prometheus parser library to return a map of metrics from a returned curl from
//an OpenTelemetry metrics endpoint
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

//Plot chooses which chart type to used based on the parsed metric type
func Plot(metric *dto.MetricFamily, workdir, resource string) error {
	if metric.GetType() == dto.MetricType_HISTOGRAM {
		promHistogramBarChart(metric, workdir, resource)
	}
	return nil
}

//generateBarItems takes a parsed metric and generates a bar series
func generateBarItems(metric *dto.Metric) plotter.Values {
	var prevGroup float64
	var items plotter.Values
	for i, bucket := range metric.Histogram.Bucket {
		var newGroupValue float64
		if i > 0 {
			newGroupValue = float64(*bucket.CumulativeCount) - prevGroup
		} else {
			newGroupValue = float64(*bucket.CumulativeCount)
		}
		prevGroup = float64(*bucket.CumulativeCount)
		items = append(items, newGroupValue)
	}
	return items
}

//generateSeriesLabels  takes a parsed metric and generates list of X Axis labels
func generateSeriesLabels(metric *dto.Metric, labels []string) []string {
	for _, bucket := range metric.Histogram.Bucket {
		label := fmt.Sprintf("%v", *bucket.UpperBound)
		if !contains(labels, label) {
			labels = append(labels, label)
		}
	}
	return labels
}

func promHistogramBarChart(metric *dto.MetricFamily, workdir, resource string) {
	var xAxisLabels []string
	var chartHeight, chartWidth vg.Length
	logrus.Debugf("%s: creating barchart for [metric=%s] from [resource=%s]", metric.GetName(), resource)
	for _, m := range metric.GetMetric() {
		var chartName string
		p := plot.New()
		p.Title.Text = metric.GetHelp()
		p.X.Label.Text = metric.GetName()
		w := vg.Points(20)
		fileName := strings.Replace(resource, ".", "_", -1)
		switch metric.GetName() {
		case EtcdNetworkPeerRoundTripTimeSeconds:
			node := m.GetLabel()
			seriesName := fmt.Sprintf("%s_%s", node[0].GetName(), node[0].GetValue())
			chartName = fmt.Sprintf("%s/%s/%s_bucket_%s.png", workdir, fileName, metric.GetName(), seriesName)
			chartHeight = 8
			chartWidth = 8
			p.Title.Text = fmt.Sprintf("%s to %s", p.Title.Text, node[0].GetValue())
		default:
			chartName = fmt.Sprintf("%s/%s/%s_bucket.png", workdir, fileName, metric.GetName())
			chartHeight = 8
			chartWidth = 8
		}
		dataValues := generateBarItems(m)
		xAxisLabels = generateSeriesLabels(m, xAxisLabels)
		chart, err := plotter.NewBarChart(dataValues, w)
		if err != nil {
			panic(err)
		}
		p.Add(chart)
		chart.LineStyle.Width = vg.Length(0)
		chart.Color = plotutil.Color(2)
		p.NominalX(xAxisLabels...)
		if err := p.Save(chartWidth*vg.Inch, chartHeight*vg.Inch, chartName); err != nil {
			panic(err)
		}
	}
}
