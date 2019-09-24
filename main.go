// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vivienv/crash-diagnostics/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
