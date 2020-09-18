// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExpandPath", func() {

	It("returns the same path when input does not contain ~", func() {
		input := "/foo/bar"
		path, err := ExpandPath(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(path).To(Equal(input))
	})

	It("replaces the ~ with home directory path", func() {
		input := "~/foo/bar"
		path, err := ExpandPath(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(path).NotTo(Equal(input))
		Expect(path).NotTo(ContainSubstring("~"))
	})
})
