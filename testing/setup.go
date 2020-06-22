// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	InfraSetupWait = time.Second * 11

	rnd              = rand.New(rand.NewSource(time.Now().Unix()))
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

//NextSSHPort returns a pseudo-rando test [2200 .. 2230]
func NextSSHPort() string {
	port := 2200 + rnd.Intn(30)
	return fmt.Sprintf("%d", port)
}

func NextSSHContainerName() string {
	return fmt.Sprintf("crashd-test-%x", rnd.Uint64())
}
