package metric

import (
	"fmt"
	"os"
	"testing"
)

type commandOutputTestData struct {
	name           string
	clientCert     string
	serverKey      string
	endpoint       string
	expectedOutput string
}

func TestEtcdMetricClient_GetCommandOutput(t *testing.T) {
	tests := []commandOutputTestData{
		{
			name:       "etcd",
			serverKey:  "",
			clientCert: "",
			endpoint:   "",
			expectedOutput: EtcdCurlCmd +
				" --cert " + EtcdDefaultClientCert +
				" --key " + EtcdDefaultKeyFile +
				" " + EtcdDefaultEndpoint,
		},
		{
			name:       "etcd",
			clientCert: "foo",
			serverKey:  "bar",
			endpoint:   "somewhere",
			expectedOutput: EtcdCurlCmd +
				" --cert foo" +
				" --key bar" +
				" somewhere",
		},
	}
	for _, test := range tests {
		client := GetClient(test.name, test.serverKey, test.clientCert, test.endpoint, "")
		command := client.GetCommandOutput()
		if command != test.expectedOutput {
			t.Logf("Expected %s Got: %s", test.expectedOutput, command)
			t.Fail()
		}
	}
}

func TestEtcdMetricClient_IsKnownMetric(t *testing.T) {
	tests := map[string]bool{
		"etcd_network_peer_round_trip_time_seconds": true,
		"etcd_disk_backend_commit_duration_seconds": true,
		"etcd_debugging_mvcc_pending_events_total":  false,
	}
	client := GetClient("etcd", "", "", "", "")
	for s, b := range tests {
		result := client.IsKnownMetric(s)
		if result != b {
			t.Logf("Metric supported %s is incorrect. Expected: %v Got: %v", s, b, result)
			t.Fail()
		}
	}
}

func TestEtcdMetricClient_Plot(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	testresults := fmt.Sprintf("%s/testdata/results", path)
	if _, err := os.Stat(testresults); os.IsNotExist(err) {
		t.Error(err)
	} else {
		os.RemoveAll(testresults)
	}
	_ = os.Mkdir(testresults, os.ModePerm)
	client := GetClient("etcd", "", "", "", testresults)
	testfile := fmt.Sprintf("%s/testdata/etcdresults.txt", path)
	results, err := client.Plot(testfile, nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(results)
	}
}
