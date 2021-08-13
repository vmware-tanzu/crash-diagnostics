// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
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

	It("throws an error when empty kube_config is used", func() {
		err = New().Exec("test.kube.config", strings.NewReader(`kube_config()`))
		Expect(err).To(HaveOccurred())
	})

	Context("With path", func() {
		Context("With kube_config set in the script", func() {

			BeforeEach(func() {
				crashdScript = `cfg = kube_config(path="/foo/bar/kube/config")`
				execSetup()
			})

			It("sets the path to the kubeconfig file", func() {
				kubeConfigData := executor.result["cfg"]
				Expect(kubeConfigData).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))

				cfg, _ := kubeConfigData.(*starlarkstruct.Struct)

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

				kubeConfigData := executor.result["cfg"]
				Expect(kubeConfigData).NotTo(BeNil())

				cfg, _ := kubeConfigData.(*starlarkstruct.Struct)

				val, err := cfg.Attr("path")
				Expect(err).To(BeNil())
				Expect(trimQuotes(val.String())).To(Equal("/foo/bar/kube/config"))
			})

			It("does not set the kube_config in the starlark thread", func() {
				kubeConfigData := executor.thread.Local(identifiers.kubeCfg)
				Expect(kubeConfigData).NotTo(BeNil())
			})
		})
	})

	Context("For default kube_config setup", func() {

		BeforeEach(func() {
			crashdScript = `foo = "bar"`
			execSetup()
		})

		It("does not set the default kube_config in the starlark thread", func() {
			kubeConfigData := executor.thread.Local(identifiers.kubeCfg)
			Expect(kubeConfigData).NotTo(BeNil())
		})
	})
})

var _ = Describe("KubeConfigFn", func() {

	Context("With capi_provider", func() {

		It("populates the path from the capi provider", func() {
			val, err := KubeConfigFn(&starlark.Thread{Name: "test.kube.config.fn"}, nil, nil,
				[]starlark.Tuple{
					[]starlark.Value{
						starlark.String("capi_provider"),
						starlarkstruct.FromStringDict(starlark.String(identifiers.capvProvider), starlark.StringDict{
							"kube_config": starlark.String("/foo/bar"),
						}),
					},
				})
			Expect(err).NotTo(HaveOccurred())

			cfg, _ := val.(*starlarkstruct.Struct)

			path, err := cfg.Attr("path")
			Expect(err).To(BeNil())
			Expect(trimQuotes(path.String())).To(Equal("/foo/bar"))
		})

		It("throws an error when an unknown provider is passed", func() {
			_, err := KubeConfigFn(&starlark.Thread{Name: "test.kube.config.fn"}, nil, nil,
				[]starlark.Tuple{
					[]starlark.Value{
						starlark.String("capi_provider"),
						starlarkstruct.FromStringDict(starlark.String("meh"), starlark.StringDict{
							"kube_config": starlark.String("/foo/bar"),
						}),
					},
				})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("unknown capi provider"))
		})
	})
})
