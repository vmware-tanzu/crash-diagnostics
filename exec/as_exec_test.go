package exec

import (
	"fmt"
	"os"
	"testing"

	"gitlab.eng.vmware.com/vivienv/flare/script"
)

func TestExecAS(t *testing.T) {
	tests := []execTest{
		{
			name: "Exec AS with userid and groupid",
			source: func() string {
				uid := os.Getuid()
				gid := os.Getgid()
				return fmt.Sprintf("AS %d:%d", uid, gid)
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "Exec AS with userid only",
			source: func() string {
				uid := os.Getuid()
				return fmt.Sprintf("AS %d", uid)
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "Exec AS with unknown uid gid",
			source: func() string {
				return "AS foo:bar"
			},
			exec: func(s *script.Script) error {
				e := New(s)
				if err := e.Execute(); err != nil {
					return err
				}
				return nil
			},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runExecutorTest(t, test)
		})
	}
}
