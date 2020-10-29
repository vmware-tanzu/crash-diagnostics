// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("kube_get", func() {

	var (
		executor *Executor
		err      error
	)

	execSetup := func(crashdScript string) {
		executor = New()
		err = executor.Exec("test.kube.get", strings.NewReader(crashdScript))
	}

	It("returns a list of k8s services as starlark objects", func() {
		crashdScript := fmt.Sprintf(`
set_defaults(kube_config(path="%s"))
kube_get_data = kube_get(groups=["core"], kinds=["services"], namespaces=["default", "kube-system"])
		`, k8sconfig)
		execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.result.Has("kube_get_data")).NotTo(BeNil())

		data := executor.result["kube_get_data"]
		Expect(data).NotTo(BeNil())

		dataStruct, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		objects, err := dataStruct.Attr("objs")
		Expect(err).NotTo(HaveOccurred())

		getDataList, _ := objects.(*starlark.List)
		Expect(getDataList.Len()).To(Equal(2))
	})

	It("returns a list of k8s nodes as starlark objects", func() {
		crashdScript := fmt.Sprintf(`
cfg = kube_config(path="%s")
kube_get_data = kube_get(groups=["core"], kinds=["nodes"], kube_config = cfg)
			`, k8sconfig)
		execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.result.Has("kube_get_data")).NotTo(BeNil())

		data := executor.result["kube_get_data"]
		Expect(data).NotTo(BeNil())

		dataStruct, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		objects, err := dataStruct.Attr("objs")
		Expect(err).NotTo(HaveOccurred())

		getDataList, _ := objects.(*starlark.List)
		Expect(getDataList.Len()).To(Equal(1))
	})

	It("returns a list of etcd containers as starlark objects", func() {
		crashdScript := fmt.Sprintf(`
kube_get_data = kube_get(namespaces=["kube-system"], containers=["etcd"], kube_config = kube_config(path="%s"))
			`, k8sconfig)
		execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.result.Has("kube_get_data")).NotTo(BeNil())

		data := executor.result["kube_get_data"]
		Expect(data).NotTo(BeNil())

		dataStruct, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		objects, err := dataStruct.Attr("objs")
		Expect(err).NotTo(HaveOccurred())

		getDataList, _ := objects.(*starlark.List)
		Expect(getDataList.Len()).To(BeNumerically(">=", 1))
	})

	It("returns a list of objects under different namespaces using categories as starlark objects", func() {
		crashdScript := fmt.Sprintf(`
kube_get_data = kube_get(categories=["all"], kube_config = kube_config(path="%s"))
			`, k8sconfig)
		execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.result.Has("kube_get_data")).NotTo(BeNil())

		data := executor.result["kube_get_data"]
		Expect(data).NotTo(BeNil())

		dataStruct, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		objects, err := dataStruct.Attr("objs")
		Expect(err).NotTo(HaveOccurred())

		getDataList, _ := objects.(*starlark.List)
		Expect(getDataList.Len()).To(BeNumerically(">=", 1))
	})

	DescribeTable("Incorrect kubeconfig", func(crashdScript string) {
		execSetup(crashdScript)
		Expect(err).To(HaveOccurred())
	},
		Entry("in global thread", fmt.Sprintf(`
set_defaults(kube_config(path="%s"))
kube_get(namespaces=["kube-system"], containers=["etcd"])`, "/foo/bar")),
		Entry("in function call", fmt.Sprintf(`
cfg = kube_config(path="%s")
kube_get(namespaces=["kube-system"], containers=["etcd"], kube_config=cfg)`, "/foo/bar")),
	)
})
