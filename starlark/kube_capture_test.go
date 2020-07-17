// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("kube_capture", func() {

	var (
		workdir  string
		executor *Executor
		err      error
	)

	execSetup := func(crashdScript string) {
		executor = New()
		err = executor.Exec("test.kube.capture", strings.NewReader(crashdScript))
	}

	BeforeEach(func() {
		workdir, err = ioutil.TempDir(os.TempDir(), "test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(workdir)
	})

	It("creates a directory and files for namespaced objects", func() {
		crashdScript := fmt.Sprintf(`
crashd_config(workdir="%s")
set_as_default(kube_config = kube_config(path="%s"))
kube_data = kube_capture(what="objects", groups=["core"], kinds=["services"], namespaces=["default", "kube-system"])
		`, workdir, k8sconfig)
		execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.result.Has("kube_data")).NotTo(BeNil())

		data := executor.result["kube_data"]
		Expect(data).NotTo(BeNil())

		dataStruct, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		fileVal, err := dataStruct.Attr("file")
		Expect(err).NotTo(HaveOccurred())

		fileValStr, ok := fileVal.(starlark.String)
		Expect(ok).To(BeTrue())

		kubeCaptureDir := fileValStr.GoString()
		Expect(kubeCaptureDir).To(BeADirectory())
		Expect(filepath.Join(kubeCaptureDir, "kube-system")).To(BeADirectory())

		Expect(filepath.Join(kubeCaptureDir, "default", "services.json")).To(BeARegularFile())
		Expect(filepath.Join(kubeCaptureDir, "kube-system", "services.json")).To(BeARegularFile())
	})

	It("creates a directory and files for non-namespaced objects", func() {
		crashdScript := fmt.Sprintf(`
crashd_config(workdir="%s")
cfg = kube_config(path="%s")
kube_data = kube_capture(what="objects", groups=["core"], kinds=["nodes"], kube_config = cfg)
		`, workdir, k8sconfig)
		execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.result.Has("kube_data")).NotTo(BeNil())

		data := executor.result["kube_data"]
		Expect(data).NotTo(BeNil())

		dataStruct, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		fileVal, err := dataStruct.Attr("file")
		Expect(err).NotTo(HaveOccurred())

		fileValStr, ok := fileVal.(starlark.String)
		Expect(ok).To(BeTrue())

		kubeCaptureDir := fileValStr.GoString()
		Expect(kubeCaptureDir).To(BeADirectory())
		Expect(filepath.Join(kubeCaptureDir, "nodes.json")).To(BeARegularFile())
	})

	It("creates a directory and log files for all objects in a namespace", func() {
		crashdScript := fmt.Sprintf(`
crashd_config(workdir="%s")
kube_data = kube_capture(what="logs", namespaces=["kube-system"], kube_config = kube_config(path="%s"))
		`, workdir, k8sconfig)
		execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.result.Has("kube_data")).NotTo(BeNil())

		data := executor.result["kube_data"]
		Expect(data).NotTo(BeNil())

		dataStruct, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		fileVal, err := dataStruct.Attr("file")
		Expect(err).NotTo(HaveOccurred())

		fileValStr, ok := fileVal.(starlark.String)
		Expect(ok).To(BeTrue())

		kubeCaptureDir := fileValStr.GoString()
		Expect(kubeCaptureDir).To(BeADirectory())
		Expect(filepath.Join(kubeCaptureDir, "kube-system")).To(BeADirectory())

		files, err := ioutil.ReadDir(filepath.Join(kubeCaptureDir, "kube-system"))
		Expect(err).NotTo(HaveOccurred())
		Expect(len(files)).NotTo(BeNumerically("<", 3))
	})

	It("creates a log file for specific container in a namespace", func() {
		crashdScript := fmt.Sprintf(`
crashd_config(workdir="%s")
cfg = kube_config(path="%s")
kube_data = kube_capture(what="logs", namespaces=["kube-system"], containers=["etcd"], kube_config = cfg)
		`, workdir, k8sconfig)
		execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.result.Has("kube_data")).NotTo(BeNil())

		data := executor.result["kube_data"]
		Expect(data).NotTo(BeNil())

		dataStruct, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		fileVal, err := dataStruct.Attr("file")
		Expect(err).NotTo(HaveOccurred())

		fileValStr, ok := fileVal.(starlark.String)
		Expect(ok).To(BeTrue())

		kubeCaptureDir := fileValStr.GoString()
		Expect(kubeCaptureDir).To(BeADirectory())
		Expect(filepath.Join(kubeCaptureDir, "kube-system")).To(BeADirectory())

		files, err := ioutil.ReadDir(filepath.Join(kubeCaptureDir, "kube-system"))
		Expect(err).NotTo(HaveOccurred())
		Expect(files).NotTo(HaveLen(0))
	})

	DescribeTable("Incorrect kubeconfig", func(crashdScript string) {
		execSetup(crashdScript)
		Expect(err).To(HaveOccurred())
	},
		Entry("in global thread", fmt.Sprintf(`
cfg = kube_config(path="%s")
kube_capture(what="logs", namespaces=["kube-system"], containers=["etcd"], kube_config = cfg)`, "/foo/bar")),
		Entry("in function call", fmt.Sprintf(`
cfg = kube_config(path="%s")
kube_capture(what="logs", namespaces=["kube-system"], containers=["etcd"], kube_config=cfg)`, "/foo/bar")),
	)
})
