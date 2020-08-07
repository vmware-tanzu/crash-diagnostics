// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	k8sconfig string
	workdir   string
)

func TestStarlarkSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Starlark Suite")
}

var _ = BeforeSuite(func() {
	// setup (if necessary) and retrieve kind's kubecfg
	k8sCfg, err := testSupport.SetupKindKubeConfig()
	Expect(err).NotTo(HaveOccurred())
	k8sconfig = k8sCfg
	workdir = testSupport.TmpDirRoot()
})

var _ = AfterSuite(func() {
	// clean up is done in main_test.go
})
