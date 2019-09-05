package exec

import (
	"fmt"

	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func exeAs(src *script.Script) (*script.AsCommand, error) {
	asCmds, ok := src.Preambles[script.CmdAs]
	if !ok {
		return nil, fmt.Errorf("Script missing valid %s", script.CmdAs)
	}
	asCmd := asCmds[0].(*script.AsCommand)
	return asCmd, nil
}
