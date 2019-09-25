// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"

	"gitlab.eng.vmware.com/vivienv/crash-diagnostics/script"
)

//exeAuthConfig retrieves a viable AuthConfig command from script
func exeAuthConfig(src *script.Script) (*script.AuthConfigCommand, error) {
	authCmds, ok := src.Preambles[script.CmdAuthConfig]
	if !ok {
		return nil, fmt.Errorf("Script missing valid %s", script.CmdAuthConfig)
	}
	authCmd := authCmds[0].(*script.AuthConfigCommand)
	return authCmd, nil
}
