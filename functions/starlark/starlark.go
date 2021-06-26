package starlark

import (
	"errors"
	"fmt"

	"github.com/vmware-tanzu/crash-diagnostics/functions/scriptconf"
	"github.com/vmware-tanzu/crash-diagnostics/functions/sshconf"
	"go.starlark.net/starlark"
)

// SetupThreadDefaults setups Starlark default thread values
func SetupThreadDefaults(thread *starlark.Thread) error {
	if thread == nil {
		return errors.New("thread defaults failed: nil thread")
	}

	if _, err := scriptconf.MakeConfigForThread(thread); err != nil {
		return fmt.Errorf("default script config: failed: %w", err)
	}
	if _, err := sshconf.MakeConfigForThread(thread); err != nil {
		return fmt.Errorf("default ssh config: failed: %w", err)
	}
	return nil
}
