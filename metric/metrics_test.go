// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"reflect"
	"testing"
)

type getClientTestData struct {
	name           string
	clientCert     string
	serverKey      string
	endpoint       string
	expectedClient EtcdMetricClient
}

func TestGetClient(t *testing.T) {
	tests := []getClientTestData{
		{
			name:       "etcd",
			clientCert: "foo",
			serverKey:  "bar",
			endpoint:   "somewhere",
			expectedClient: EtcdMetricClient{
				Name:             "etcd",
				SupportedMetrics: etcdKnownMetrics(),
				ServerKey:        "bar",
				ClientCert:       "foo",
				Endpoint:         "somewhere",
				WorkDir:          "",
			},
		},
		{
			name:       "etcd",
			clientCert: "",
			serverKey:  "",
			endpoint:   "",
			expectedClient: EtcdMetricClient{
				Name:             "etcd",
				SupportedMetrics: etcdKnownMetrics(),
				ServerKey:        "",
				ClientCert:       "",
				Endpoint:         "",
				WorkDir:          "",
			},
		},
	}
	for _, test := range tests {
		client := GetClient(test.name, test.serverKey, test.clientCert, test.endpoint, "")
		if !reflect.DeepEqual(client, test.expectedClient) {
			t.Logf("Expected: %v, Got: %s", test.expectedClient, client)
			t.Fail()
		}
	}

}
