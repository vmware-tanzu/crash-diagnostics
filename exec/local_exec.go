package exec

import (
	"github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func exeLocally(src *script.Script, workdir string) error {
	envPairs := exeEnvs(src)
	asCmd, err := exeAs(src)
	if err != nil {
		return err
	}

	for _, action := range src.Actions {
		switch cmd := action.(type) {
		case *script.CopyCommand:
			if err := exeCopy(asCmd, cmd, workdir); err != nil {
				return err
			}
		case *script.CaptureCommand:
			// capture command output
			if err := exeCapture(asCmd, cmd, envPairs, workdir); err != nil {
				return err
			}
		default:
			logrus.Errorf("Unsupported command %T", cmd)
		}
	}

	return nil
}
