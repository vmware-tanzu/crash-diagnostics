// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"flag"
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	InfraSetupWait = time.Second * 11

	rnd              = rand.New(rand.NewSource(time.Now().Unix()))
	sshContainerName = "test-sshd"
	sshPort          = NextPortValue()
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

//NextPortValue returns a pseudo-rando test [2200 .. 2230]
func NextPortValue() string {
	port := 2200 + rnd.Intn(90)
	return fmt.Sprintf("%d", port)
}

func NextResourceName() string {
	return fmt.Sprintf("crashd-test-%x", rnd.Uint64())
}

func GetSSHKeyDirectory() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return path.Join(filepath.Dir(d), "testing", "keys")
}

func GetSSHPrivateKey() string {
	return filepath.Join(GetSSHKeyDirectory(), "id_rsa")
}

func GetSSHUsername() string {
	return "vivienv"
}
