// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("set_defaults", func() {

	DescribeTable("sets the inputs as default", func(crashdScript string) {
		e := New()
		err := e.Exec("test.set_defaults", strings.NewReader(crashdScript))
		Expect(err).NotTo(HaveOccurred())

		kubeConfig := e.thread.Local(identifiers.kubeCfg)
		Expect(kubeConfig).NotTo(BeNil())
		Expect(kubeConfig).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))

		sshConfig := e.thread.Local(identifiers.sshCfg)
		Expect(sshConfig).NotTo(BeNil())
		Expect(sshConfig).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))

		resources := e.thread.Local(identifiers.resources)
		Expect(resources).NotTo(BeNil())
		Expect(resources).To(BeAssignableToTypeOf(&starlark.List{}))
	},
		Entry("single inputs", `
kube_cfg = kube_config(path="/foo/bar")
set_defaults(kube_cfg)

ssh_cfg = ssh_config(username="baz")
set_defaults(ssh_cfg)

res = resources(hosts=["127.0.0.1","localhost"])
set_defaults(res)
`),
		Entry("single inputs with inline declarations", `
set_defaults(kube_config(path="/foo/bar"))
set_defaults(ssh_config(username="baz"))
set_defaults(resources(hosts=["127.0.0.1","localhost"]))
`),
		Entry("multiple inputs with inline declarations", `
set_defaults(kube_config(path="/foo/bar"), ssh_config(username="baz"))
set_defaults(resources(hosts=["127.0.0.1","localhost"]))
`),
		Entry("multiple inputs with inline declarations", `
set_defaults(ssh_config(username="baz"))
set_defaults(kube_config(path="/foo/bar"), resources(hosts=["127.0.0.1","localhost"]))
`),
	)

	Context("When a default ssh_config is not declared", func() {

		It("fails to evaluate resources as a set_defaults option", func() {
			e := New()
			err := e.Exec("test.set_defaults", strings.NewReader(`
ssh_cfg = ssh_config(username="baz")
set_defaults(resources = resources(hosts=["127.0.0.1","localhost"]), ssh_config = ssh_cfg)
`))
			Expect(err).To(HaveOccurred())
		})
	})

	DescribeTable("throws an error", func(crashdScript string) {
		e := New()
		err := e.Exec("test.set_defaults", strings.NewReader(crashdScript))
		Expect(err).To(HaveOccurred())

		kubeConfig := e.thread.Local(identifiers.kubeCfg)
		Expect(kubeConfig).NotTo(BeNil())

		sshConfig := e.thread.Local(identifiers.sshCfg)
		Expect(sshConfig).NotTo(BeNil())
	}, Entry("no input", `
kube_cfg = kube_config(path="/foo/bar")
ssh_cfg = ssh_config(username="baz")
set_defaults()
`),
		Entry("incorrect input", `
set_defaults("/foo")
`),
		Entry("keyword inputs", `
set_defaults(kube_config = kube_config(path="/foo/bar"))
`),
	)
})
