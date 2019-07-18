package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vivienv/flare/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
