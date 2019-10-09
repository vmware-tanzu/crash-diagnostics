package common

import (
	"fmt"

	"github.com/vladimirvivien/echo"
)

var (
	// PkgRoot project package root
	PkgRoot string
	// Version default build version
	Version string
	// GitSHA last commit sha
	GitSHA string
)

func init() {
	e := echo.New()
	PkgRoot = "github.com/vmware-tanzu/crash-diagnostics"
	Version = fmt.Sprintf("%s-unreleased", e.Run("git rev-parse --abbrev-ref HEAD"))
	GitSHA = e.Run("git rev-parse HEAD")
}
