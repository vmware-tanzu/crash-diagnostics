// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"fmt"
	"os"
	"testing"
)

func TestPlot(t *testing.T) {
	names := []string{"etcd"}
	for _, name := range names {

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
		client := GetClient(name, "", "", "", testresults)
		testfile := fmt.Sprintf("%s/testdata/etcdresults.txt", path)
		results, err := client.Plot(testfile, nil)
		if err != nil {
			t.Log(err)
			t.Fail()
		} else {
			t.Log(results)
		}
	}
}
