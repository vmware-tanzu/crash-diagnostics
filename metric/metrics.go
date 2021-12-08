////// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
////// SPDX-License-Identifier: Apache-2.0
////
package metric

type MetricGraph interface {
	Plot(string, []string) ([]string, error)
	GetCommandOutput() string
	IsKnownMetric(string) bool
}

func GetClient(name, serverKey, clientCert, endpoint, workDir string) MetricGraph {
	switch name {
	case "etcd":
		return EtcdMetricClient{
			Name:             "etcd",
			SupportedMetrics: etcdKnownMetrics(),
			ServerKey:        serverKey,
			ClientCert:       clientCert,
			Endpoint:         endpoint,
			WorkDir:          workDir,
		}
	}
	return nil
}
