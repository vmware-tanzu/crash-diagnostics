// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("set_as_default", func() {

	It("sets the inputs as default", func() {
		e := New()
		err := e.Exec("test.set_as_default", strings.NewReader(`
kube_cfg = kube_config(path="/foo/bar")
ssh_cfg = ssh_config(username="baz")
set_as_default(ssh_config = ssh_cfg, kube_config = kube_cfg)
set_as_default(resources = resources(hosts=["127.0.0.1","localhost"]))
`))
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
	})

	Context("When a default ssh_config is not declared", func() {

		It("fails to evaluate resources as a set_as_default option", func() {
			e := New()
			err := e.Exec("test.set_as_default", strings.NewReader(`
ssh_cfg = ssh_config(username="baz")
set_as_default(resources = resources(hosts=["127.0.0.1","localhost"]), ssh_config = ssh_cfg)
`))
			Expect(err).To(HaveOccurred())
		})
	})

	It("throws an error", func() {
		e := New()
		err := e.Exec("test.set_as_default", strings.NewReader(`
kube_cfg = kube_config(path="/foo/bar")
ssh_cfg = ssh_config(username="baz")
set_as_default()
`))
		Expect(err).To(HaveOccurred())

		kubeConfig := e.thread.Local(identifiers.kubeCfg)
		Expect(kubeConfig).To(BeNil())

		sshConfig := e.thread.Local(identifiers.sshCfg)
		Expect(sshConfig).To(BeNil())
	})
})
