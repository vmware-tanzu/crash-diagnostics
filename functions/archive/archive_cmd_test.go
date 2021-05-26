package archive

import (
	"os"
	"testing"

	"go.starlark.net/starlark"
)

func TestCmd_Run(t *testing.T) {
	tests := []struct {
		name       string
		params     Args
		arc        Result
		shouldFail bool
	}{
		{
			name:       "empty param",
			shouldFail: true,
		},
		{
			name:   "default archive name",
			params: Args{SourcePaths: []string{"/tmp/crashd"}},
			arc:    Result{OutputFile: DefaultBundleName},
		},
		{
			name:   "archive name",
			params: Args{SourcePaths: []string{"/tmp/crashd"}, OutputFile: "test.tar.gz"},
			arc:    Result{OutputFile: "test.tar.gz"},
		},
		{
			name:   "multiple files",
			params: Args{SourcePaths: []string{"/tmp/crashd0", "/tmp/crashd1"}, OutputFile: "test.tar.gz"},
			arc:    Result{OutputFile: "test.tar.gz"},
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

			if result.Error != "" && !test.shouldFail {
				t.Errorf("unexpected error: %s", result.Error)
			}

			if result.Size == 0 && !test.shouldFail {
				t.Errorf("archive file has size 0")
			}

			// clean up
			if err := os.RemoveAll(result.OutputFile); err != nil {
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
