package script_tests

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	testcrashd "github.com/vmware-tanzu/crash-diagnostics/testing"
)

var (
	testSupport *testcrashd.TestSupport
)

func TestMain(m *testing.M) {
	test, err := testcrashd.Init()
	if err != nil {
		logrus.Fatal(err)
	}
	testSupport = test
	// precaution
	if testSupport == nil {
		logrus.Fatal("failed to setup test support")
	}

	if err := testSupport.SetupSSHServer(); err != nil {
		logrus.Fatal(err)
	}

	result := m.Run()

	if err := testSupport.TearDown(); err != nil {
		logrus.Fatal(err)
	}

	os.Exit(result)
}
