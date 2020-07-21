// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kube_nodes_provider", func() {
	var (
		executor *Executor
		err      error
	)

	execSetup := func(crashdScript string) error {
		executor = New()
		return executor.Exec("test.kube.nodes.provider", strings.NewReader(crashdScript))
	}

	It("returns a struct with the list of k8s nodes", func() {
		crashdScript := fmt.Sprintf(`
cfg = kube_config(path="%s")
provider = kube_nodes_provider(kube_config = cfg, ssh_config = ssh_config(username="uname", private_key_path="path"))`, k8sconfig)
		err = execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())

		data := executor.result["provider"]
		Expect(data).NotTo(BeNil())

		provider, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		val, err := provider.Attr("hosts")
		Expect(err).NotTo(HaveOccurred())

		list := val.(*starlark.List)
		Expect(list.Len()).To(Equal(1))
	})

	It("returns a struct with ssh config", func() {
		crashdScript := fmt.Sprintf(`
cfg = kube_config(path="%s")
ssh_cfg = ssh_config(username="uname", private_key_path="path")
provider = kube_nodes_provider(kube_config=cfg, ssh_config = ssh_cfg)`, k8sconfig)
		err = execSetup(crashdScript)
		Expect(err).NotTo(HaveOccurred())

		data := executor.result["provider"]
		Expect(data).NotTo(BeNil())

		provider, ok := data.(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		sshCfg, err := provider.Attr(identifiers.sshCfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(sshCfg).NotTo(BeNil())
	})
})
