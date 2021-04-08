package archive

import (
	"os"
	"testing"

	"go.starlark.net/starlark"
)

func TestCmd_Run(t *testing.T) {
	tests := []struct {
		name       string
		params     Params
		arc        Archive
		shouldFail bool
	}{
		{
			name:       "empty param",
			shouldFail: true,
		},
		{
			name:   "default archive name",
			params: Params{SourcePaths: []string{"/tmp/crashd"}},
			arc:    Archive{SourcePaths: []string{"/tmp/crashd"}, OutputFile: DefaultBundleName},
		},
		{
			name:   "archive name",
			params: Params{SourcePaths: []string{"/tmp/crashd"}, OutputFile: "test.tar.gz"},
			arc:    Archive{SourcePaths: []string{"/tmp/crashd"}, OutputFile: "test.tar.gz"},
		},
		{
			name:   "multiple files",
			params: Params{SourcePaths: []string{"/tmp/crashd0", "/tmp/crashd1"}, OutputFile: "test.tar.gz"},
			arc:    Archive{SourcePaths: []string{"/tmp/crashd0", "/tmp/crashd1"}, OutputFile: "test.tar.gz"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create source dirs
			for _, p := range test.params.SourcePaths {
				if err := os.MkdirAll(p, 0744); err != nil {
					t.Fatal(err)
				}
			}

			result, err := newCmd().Run(&starlark.Thread{}, test.params)
			if err != nil {
				t.Fatal(err)
			}

			arc, ok := result.Value().(Archive)
			if !ok {
				t.Fatalf("unexpected type %T returned by function %s", result.Value(), FuncName)
			}

			if result.Err() != "" && !test.shouldFail {
				t.Errorf("unexpected error: %s", result.Err())
			}

			if len(arc.SourcePaths) != len(test.arc.SourcePaths) {
				t.Errorf("unexpected source paths length: %d", len(arc.SourcePaths))
			}

			if arc.Size == 0 && !test.shouldFail {
				t.Errorf("archive file has size 0")
			}

			// clean up
			if err := os.RemoveAll(arc.OutputFile); err != nil {
				t.Log(err)
			}

			for _, p := range test.params.SourcePaths {
				if err := os.RemoveAll(p); err != nil {
					t.Log(err)
				}
			}

		})
	}
}
