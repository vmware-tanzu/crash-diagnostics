// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/vladimirvivien/gexe"
)

func main() {
	e := gexe.New()
	if e.Run("gofmt -s -l .") != "" {
		fmt.Println("Go code failed gofmt check:")
		e.Runout("gofmt -s -d .")
		os.Exit(1)
	}
}
