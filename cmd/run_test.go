// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run", func() {

	Context("For args-file and args", func() {

		var argsBackupFile string

		JustBeforeEach(func() {
			if _, err := os.Stat(ArgsFile); err == nil {
				argsBackupFile = fmt.Sprintf("%s.BKP.%s", ArgsFile, time.Now().String())
				Expect(os.Rename(ArgsFile, argsBackupFile)).NotTo(HaveOccurred())
			}
		})

		JustAfterEach(func() {
			if argsBackupFile != "" {
				Expect(os.Rename(argsBackupFile, ArgsFile)).NotTo(HaveOccurred())
			}
		})

		DescribeTable("processScriptArguments", func(argsFileContent string, args map[string]string, size int) {
			f, err := os.CreateTemp(os.TempDir(), "")
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(f.Name(), []byte(argsFileContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			defer f.Close()

			flags := &runFlags{
				args:     args,
				argsFile: f.Name(),
			}
			scriptArgs, err := processScriptArguments(flags)
			Expect(err).NotTo(HaveOccurred())
			Expect(scriptArgs).To(HaveLen(size))
		},
			Entry("no overlapping keys", "key=value", map[string]string{"a": "b"}, 2),
			Entry("overlapping keys", "key=value", map[string]string{"key": "b"}, 1),
			Entry("file with no keys", "", map[string]string{"key": "b"}, 1),
			Entry("with file and without args", "key=value", map[string]string{}, 1),
		)

		Context("With no default args file", func() {
			DescribeTable("does not throw an error", func(args map[string]string, size int) {
				scriptArgs, err := processScriptArguments(&runFlags{
					args:     args,
					argsFile: ArgsFile,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(scriptArgs).To(HaveLen(size))
			},
				Entry("with args", map[string]string{"a": "b"}, 1),
				Entry("without args", map[string]string{}, 0),
			)
		})
	})
})
