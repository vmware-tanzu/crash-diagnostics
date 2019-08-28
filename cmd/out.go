package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gitlab.eng.vmware.com/vivienv/flare/archiver"

	"github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/vivienv/flare/exec"

	"gitlab.eng.vmware.com/vivienv/flare/script"

	"github.com/spf13/cobra"
)

type outFlags struct {
	file   string
	output string
}

func newOutCommand() *cobra.Command {
	flags := &outFlags{
		file:   "flare.file",
		output: "out.tar.gz",
	}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "out",
		Short: "outputs an archive from collected data",
		Long:  "outputs an archive from data collected from the specified machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOut(flags, args)
		},
	}
	cmd.Flags().StringVar(&flags.file, "file", flags.file, "the path to the flare.file (default ./flare.file)")
	cmd.Flags().StringVar(&flags.output, "output", flags.output, "the path to the generated archive file (default out.tar.gz)")
	return cmd
}

func runOut(flag *outFlags, args []string) error {
	var ff io.Reader

	file, err := os.Open(flag.file)
	if err != nil {
		logrus.Warnf("Unable find %s, using sensible defaults", flag.file)
		ff = bytes.NewReader([]byte(flarefile))
	} else {
		ff = file
	}
	defer file.Close()

	flare, err := script.Parse(ff)
	if err != nil {
		return err
	}

	exe := exec.New(flare)
	if err := exe.Execute(); err != nil {
		return err
	}

	workdirs, ok := flare.Preambles[script.CmdWorkDir]
	if !ok {
		return fmt.Errorf("Script missing %s", script.CmdWorkDir)
	}
	workdir := workdirs[0].(*script.WorkdirCommand)
	if err := archiver.Tar(flag.output, workdir.Dir()); err != nil {
		return err
	}
	logrus.Infof("Created archive %s", flag.output)
	logrus.Info("Output done")

	return nil
}
