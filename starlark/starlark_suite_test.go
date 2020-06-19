package starlark

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStarlark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Starlark Suite")
}
