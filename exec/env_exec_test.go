package exec

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func createTestShellScript(t *testing.T, fname string, content string) error {
	execFile, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer execFile.Close()
	t.Logf("Creating shell script file %s", fname)
	_, err = io.Copy(execFile, strings.NewReader(content))
	return err
}
func TestExecENV(t *testing.T) {
	tests := []execTest{
		{
			name: "ENV with multiple key/values",
			source: func() string {
				return "ENV MSG1=HELLO\nENV MSG2=WORLD MSG3=!\nCAPTURE ./foo.sh"
			},
			exec: func(s *script.Script) error {
				// create an executable script to apply ENV
				scriptName := "foo.sh"
				sh := "#!/bin/sh\necho $MSG1 $MSG2 $MSG3"
				msgExpected := "HELLO WORLD !"
				if err := createTestShellScript(t, scriptName, sh); err != nil {
					return err
				}
				defer os.RemoveAll(scriptName)

				machine := s.Preambles[script.CmdFrom][0].(*script.FromCommand).Machines()[0].Address
				workdir := s.Preambles[script.CmdWorkDir][0].(*script.WorkdirCommand)
				cmd := s.Actions[0].(*script.CaptureCommand)

				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}

				fileName := filepath.Join(workdir.Dir(), machine, fmt.Sprintf("%s.txt", flatCmd(cmd.GetCliString())))
				if _, err := os.Stat(fileName); err != nil {
					return err
				}

				file, err := ioutil.ReadFile(fileName)
				if err != nil {
					return err
				}
				if strings.TrimSpace(string(file)) != msgExpected {
					return fmt.Errorf("ENV value not applied during CAPATURE")
				}

				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}
