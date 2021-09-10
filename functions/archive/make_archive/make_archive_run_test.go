package make_archive

import (
	"os"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/functions/archive"
	"go.starlark.net/starlark"
)

func TestArchiveRun(t *testing.T) {
	tests := []struct {
		name       string
		params     archive.Args
		arc        archive.Result
		shouldFail bool
	}{
		{
			name:       "empty param",
			shouldFail: true,
		},
		{
			name:   "default archive name",
			params: archive.Args{SourcePaths: []string{"/tmp/crashd"}},
			arc:    archive.Result{Archive: archive.Archive{OutputFile: DefaultBundleName}},
		},
		{
			name:   "archive name",
			params: archive.Args{SourcePaths: []string{"/tmp/crashd"}, OutputFile: "test.tar.gz"},
			arc:    archive.Result{Archive: archive.Archive{OutputFile: "test.tar.gz"}},
		},
		{
			name:   "multiple files",
			params: archive.Args{SourcePaths: []string{"/tmp/crashd0", "/tmp/crashd1"}, OutputFile: "test.tar.gz"},
			arc:    archive.Result{Archive: archive.Archive{OutputFile: "test.tar.gz"}},
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

			result := Run(&starlark.Thread{}, test.params)
			if result.Error != "" && !test.shouldFail {
				t.Errorf("unexpected error: %s", result.Error)
			}

			if result.Archive.Size == 0 && !test.shouldFail {
				t.Errorf("archive file has size 0")
			}

			// clean up
			if err := os.RemoveAll(result.Archive.OutputFile); err != nil {
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
