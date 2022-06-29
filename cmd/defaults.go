package cmd

import (
	"os"
	"path/filepath"
)

var (
	// CrashdDir is the directory path created at crashd runtime
	CrashdDir = filepath.Join(os.Getenv("HOME"), ".crashd")

	// ArgsFile is the path of the defaults args file.
	ArgsFile = filepath.Join(CrashdDir, "args")
)

// This creates a crashd directory which can be used as a default workdir
// for script execution. It will also house the default args file.
func CreateCrashdDir() error {
	if _, err := os.Stat(CrashdDir); os.IsNotExist(err) {
		return os.Mkdir(CrashdDir, 0755)
	}
	return nil
}
