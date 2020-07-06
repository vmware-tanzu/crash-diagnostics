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

var _ = Describe("resources with kube_nodes_provider()", func() {

	It("populates the resources with the cluster nodes as hosts", func() {
		crashdScript := fmt.Sprintf(`
cfg = kube_config(path="%s")
ssh_config(username="uname", private_key_path="path")
res = resources(provider=kube_nodes_provider(kube_config=cfg))`, k8sconfig)

		executor := New()
		err := executor.Exec("test.resources.kube.nodes.provider", strings.NewReader(crashdScript))
		Expect(err).NotTo(HaveOccurred())

		data := executor.result["res"]
		Expect(data).NotTo(BeNil())

		resources, ok := data.(*starlark.List)
		Expect(ok).To(BeTrue())
		Expect(resources.Len()).To(Equal(1))

		resStruct, ok := resources.Index(0).(*starlarkstruct.Struct)
		Expect(ok).To(BeTrue())

		val, err := resStruct.Attr("kind")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(val.String())).To(Equal(identifiers.hostResource))

		transport, err := resStruct.Attr("transport")
		Expect(err).NotTo(HaveOccurred())
		Expect(trimQuotes(transport.String())).To(Equal("ssh"))

		sshCfg, err := resStruct.Attr(identifiers.sshCfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(sshCfg).NotTo(BeNil())

		host, err := resStruct.Attr("host")
		Expect(err).NotTo(HaveOccurred())
		// Regex to match IP address of the host
		Expect(trimQuotes(host.String())).To(MatchRegexp("^([1-9]?[0-9]{2}\\.)([0-9]{1,3}\\.){2}[0-9]{1,3}$"))
	})
})
