// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	kind      *testcrashd.KindCluster
	waitTime  = time.Second * 11
	k8sconfig string
)

func TestStarlark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Starlark Suite")
}

var _ = BeforeSuite(func() {
	clusterName := "crashd-test-cluster"
	tmpFile, err := ioutil.TempFile(os.TempDir(), clusterName)
	Expect(err).NotTo(HaveOccurred())
	k8sconfig = tmpFile.Name()

	// create kind cluster
	kind = testcrashd.NewKindCluster("../testing/kind-cluster-docker.yaml", clusterName)
	err = kind.Create()
	Expect(err).NotTo(HaveOccurred())

	err = kind.MakeKubeConfigFile(k8sconfig)
	Expect(err).NotTo(HaveOccurred())

	logrus.Infof("Sleeping %v ... waiting for pods", waitTime)
	time.Sleep(waitTime)
})

var _ = AfterSuite(func() {
	kind.Destroy()
	os.RemoveAll(k8sconfig)
})
