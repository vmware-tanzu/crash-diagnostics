// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"testing"
)

func TestCreateKindCluster(t *testing.T) {
	k := NewKindCluster("./kind-cluster-docker.yaml", "testing-test-cluster", "")
	if err := k.Create(); err != nil {
		t.Error(err)
	}

	if err := k.Create(); err != nil {
		t.Error(err)
	}

	if k.GetKubeCtlContext() != "kind-testing-test-cluster" {
		t.Errorf("Unexpected kubectl context name %s", k.GetKubeCtlContext())
	}

	if err := k.Destroy(); err != nil {
		t.Error(err)
	}
}
