// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/sirupsen/logrus"
)

var _ = Describe("ReadArgsFile", func() {
	var args map[string]string

	BeforeEach(func() {
		args = map[string]string{}
	})

	It("returns an error when an invalid file name is passed", func() {
		err := ReadArgsFile("/foo/blah", args)
		Expect(err).To(HaveOccurred())
	})

	Context("with valid file", func() {
		DescribeTable("length of args map", func(input string, size int, warnMsgPresent bool) {
			f := writeContentToFile(input)
			defer f.Close()

			warnBuffer := gbytes.NewBuffer()
			logrus.SetOutput(warnBuffer)

			err := ReadArgsFile(f.Name(), args)
			Expect(err).NotTo(HaveOccurred())
			Expect(args).To(HaveLen(size))

			if warnMsgPresent {
				Expect(warnBuffer).To(gbytes.Say("unknown entry in args file"))
			}
		},
			Entry("valid with no spaces", `
key=value
foo=bar
`, 2, false),
			Entry("valid with spaces", `
# key represents earth is round
key = value
foo= bar
bloop =blah
		`, 3, false),
			Entry("valid with empty values", `
key =
foo= bar
bloop=
		`, 3, false),
			Entry("invalid", `
key value
foo
bar
`, 0, true))
	})

	It("accepts comments in the args file", func() {
		f := writeContentToFile(`# key represents A
key = value
# foo represents B
foo= bar`)
		defer f.Close()

		err := ReadArgsFile(f.Name(), args)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(HaveLen(2))
	})

})

var writeContentToFile = func(content string) *os.File {
	f, err := ioutil.TempFile(os.TempDir(), "read_file_args")
	Expect(err).NotTo(HaveOccurred())

	err = ioutil.WriteFile(f.Name(), []byte(content), 0644)
	Expect(err).NotTo(HaveOccurred())

	return f
}
