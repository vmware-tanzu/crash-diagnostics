// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"flag"

	"github.com/sirupsen/logrus"
)

var (
	sshContainerName = "test-sshd"
	sshPort          = "2222"
)

// Init initializes testing
func Init() {
	debug := false
	flag.BoolVar(&debug, "debug", debug, "Enables debug level")
	flag.Parse()

	logLevel := logrus.InfoLevel
	if debug {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)
}

// DefaultSSHPort is the default SSH port
func DefaultSSHPort() string {
	return sshPort
}
