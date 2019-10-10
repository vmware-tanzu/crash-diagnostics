// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/vladimirvivien/echo"
	ci "github.com/vmware-tanzu/crash-diagnostics/.ci/common"
)

func main() {
	arches := []string{"amd64"}
	oses := []string{"darwin", "linux"}

	e := echo.New()
	e.SetEnv("PKG_ROOT", ci.PkgRoot)
	e.SetEnv("VERSION", ci.Version)
	e.SetEnv("GIT_SHA", ci.GitSHA)
	e.SetEnv("LDFLAGS", `"-X ${PKG_ROOT}/buildinfo.Version=${VERSION} -X ${PKG_ROOT}/buildinfo.GitSHA=${GIT_SHA}"`)

	for _, arch := range arches {
		for _, os := range oses {
			binary := fmt.Sprintf(".build/%s/%s/crash-diagnostics", arch, os)
			gobuild(arch, os, e.Val("LDFLAGS"), binary)
		}
	}
}

func gobuild(arch, os, ldflags, binary string) {
	b := echo.New()
	b.Conf.SetPanicOnErr(true)
	b.SetVar("arch", arch)
	b.SetVar("os", os)
	b.SetVar("ldflags", ldflags)
	b.SetVar("binary", binary)
	result := b.Env("CGO_ENABLED=0 GOOS=$os GOARCH=$arch").Run("go build -o $binary -ldflags $ldflags .")
	if !b.Empty(result) {
		fmt.Printf("Build for %s/%s failed: %s\n", arch, os, result)
		return
	}
	fmt.Printf("Build %s/%s OK: %s\n", arch, os, binary)
}
