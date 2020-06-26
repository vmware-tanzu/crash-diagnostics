// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"

	"go.starlark.net/starlarkstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kube_config", func() {

	var (
		crashdScript string
		executor     *Executor
		err          error
	)

	execSetup := func() {
		executor = New()
		err = executor.Exec("test.kube.config", strings.NewReader(crashdScript))
		Expect(err).To(BeNil())
	}

	Context("With kube_config set in the script", func() {

		BeforeEach(func() {
			crashdScript = `kube_config(path="/foo/bar/kube/config")`
			execSetup()
		})

		It("sets the kube_config in the starlark thread", func() {
			kubeConfigData := executor.thread.Local(identifiers.kubeCfg)
			Expect(kubeConfigData).NotTo(BeNil())
		})

		It("sets the path to the kubeconfig file", func() {
			kubeConfigData := executor.thread.Local(identifiers.kubeCfg)
			Expect(kubeConfigData).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))

			cfg, _ := kubeConfigData.(*starlarkstruct.Struct)
			Expect(cfg.AttrNames()).To(HaveLen(1))

			val, err := cfg.Attr("path")
			Expect(err).To(BeNil())
			Expect(trimQuotes(val.String())).To(Equal("/foo/bar/kube/config"))
		})
	})

	Context("With kube_config returned as a value", func() {

		BeforeEach(func() {
			crashdScript = `cfg = kube_config(path="/foo/bar/kube/config")`
			execSetup()
		})

		It("returns the kube config as a result", func() {
			Expect(executor.result.Has("cfg")).NotTo(BeNil())
		})

		It("also sets the kube_config in the starlark thread", func() {
			kubeConfigData := executor.thread.Local(identifiers.kubeCfg)
			Expect(kubeConfigData).NotTo(BeNil())

			cfg, _ := kubeConfigData.(*starlarkstruct.Struct)
			Expect(cfg.AttrNames()).To(HaveLen(1))

			val, err := cfg.Attr("path")
			Expect(err).To(BeNil())
			Expect(trimQuotes(val.String())).To(Equal("/foo/bar/kube/config"))
		})
	})

	Context("With default kube_config setup", func() {

		BeforeEach(func() {
			crashdScript = `foo = "bar"`
			execSetup()
		})

		It("sets the default kube_config in the starlark thread", func() {
			kubeConfigData := executor.thread.Local(identifiers.kubeCfg)
			Expect(kubeConfigData).NotTo(BeNil())

			cfg, _ := kubeConfigData.(*starlarkstruct.Struct)
			Expect(cfg.AttrNames()).To(HaveLen(1))

			val, err := cfg.Attr("path")
			Expect(err).To(BeNil())
			Expect(trimQuotes(val.String())).To(ContainSubstring("/.kube/config"))
		})
	})
})
