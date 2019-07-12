package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandProcessor_New(t *testing.T) {
	proc := NewCommandProcessor(
		[]string{"echo helloWorld"},
	)
	if len(proc.cmds) == 0 {
		t.Error("commands not set")
	}
	if len(proc.workDir) == 0 {
		t.Error("workdir not set")
	}
	if len(proc.outputPath) == 0 {
		t.Error("output dir not set")
	}
}

func TestCommandProcessor_RedirectToFile(t *testing.T) {
	tests := []struct {
		name           string
		executor       *CommandExecutor
		fileName       string
		expectedOutput string
		shouldFail     bool
	}{
		{
			name:           "redirect to file OK",
			executor:       NewCommandExecutor("echo", "HelloWorld"),
			expectedOutput: "HelloWorld",
			fileName:       fmt.Sprintf("/tmp/echo_HelloWorld%x.txt", rand.Int31()),
		},
		{
			name:       "bad command",
			executor:   NewCommandExecutor("foo"),
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &CommandProcessor{}
			reader, err := test.executor.Execute()
			if err != nil {
				if !test.shouldFail {
					t.Fatalf("unexpected error: %s", err)
				}
				t.Log(err)
				return
			}
			if err := p.redirectToFile(reader, test.fileName); err != nil {
				t.Fatal(err)
			}
			// is file created ok
			if _, err := os.Stat(test.fileName); err != nil {
				t.Fatal(err)
			}
			// test content
			file, err := os.Open(test.fileName)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				file.Close()
				if err := os.Remove(test.fileName); err != nil {
					t.Error(err)
				}
			}()

			result, err := ioutil.ReadFile(test.fileName)
			if err != nil {
				t.Fatal(err)
			}

			if strings.TrimSpace(string(result)) != test.expectedOutput {
				t.Errorf("unexpected content: %s", result)
			}

		})
	}
}

func TestCommandProcessor_GenFileName(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		fileName string
	}{
		{
			name:     "simple command",
			cmd:      "echo helloWorld",
			fileName: "echo_helloWorld",
		},

		{
			name:     "command with multiple named flags",
			cmd:      "iptables -L --n",
			fileName: "iptables_-L_--n",
		},
		{
			name:     "command with quoted flags",
			cmd:      `foo -A "hello world" --d 'bar'`,
			fileName: "foo_-A_hello_world_--d_bar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &CommandProcessor{}
			name := p.genFileName(test.cmd)
			if name != test.fileName {
				t.Errorf("expecting fileName %s, got %s", test.fileName, name)
			}
		})
	}
}

func TestCommandProcessor_Tar(t *testing.T) {
	tests := []struct {
		name     string
		workDir  string
		tarFile  string
		contents map[string]string
	}{
		{
			name:     "tar one file in same loc",
			workDir:  "/tmp/flaretest",
			tarFile:  "/tmp/out.tar.gz",
			contents: map[string]string{"flaretest/hello.txt": "helloWorld"},
		},
		// {
		// 	name:     "tar multiple files in same loc",
		// 	workDir:  "/tmp/flaretest",
		// 	tarFile:  "/tmp/out1.tar.gz",
		// 	contents: []string{"helloWorld", "hello universe!", "foo bar"},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &CommandProcessor{}

			workDir := test.workDir
			if err := os.MkdirAll(workDir, 0755); err != nil && !os.IsExist(err) {
				t.Fatal(err)
			}

			// create contents
			for name, content := range test.contents {
				reader := strings.NewReader(content)
				fileName := filepath.Join(workDir, name)
				if err := p.redirectToFile(reader, fileName); err != nil {
					t.Fatal(err)
				}
			}

			// create tar file
			if err := p.tar(test.tarFile, workDir); err != nil {
				t.Logf("failed to tar %s", workDir)
				t.Fatal(err)
			}

			if err := os.RemoveAll(workDir); err != nil {
				t.Logf("failed to clean working dir %s", workDir)
				t.Fatal(err)
			}

			// validate tar
			if _, err := os.Stat(test.tarFile); err != nil {
				t.Fatalf("unable to Stat %s", test.tarFile)
			}

			gzFile, err := os.Open(test.tarFile)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(test.tarFile)
			defer gzFile.Close()

			unzipper, err := gzip.NewReader(gzFile)
			if err != nil {
				t.Fatal(err)
			}
			defer unzipper.Close()
			untarrer := tar.NewReader(unzipper)

			for {
				file, err := untarrer.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}

				filePath := filepath.Join(filepath.Dir(test.tarFile), file.Name)
				mod := file.FileInfo().Mode()
				switch {
				case mod.IsDir():
					if err := os.MkdirAll(filePath, 0755); err != nil {
						t.Fatal(err)
					}
				case mod.IsRegular():
					var buf bytes.Buffer
					n, err := io.Copy(&buf, untarrer)
					if err != nil {
						t.Fatal(err)
					}
					if n != file.Size {
						t.Fatalf("unexpected file size extracted from tar: %d", file.Size)
					}
					if buf.String() != test.contents[file.Name] {
						t.Errorf("unexpected content from archiver for %s: %s", file.Name, buf.String())
					}
				default:
					t.Fatal("unsupported tar file type")
				}

			}
		})
	}
}

func TestCommandProcessor_Process(t *testing.T) {
	tests := []struct {
		name       string
		cmds       []string
		fileNames  []string
		outputs    []string
		shouldFail bool
	}{
		{
			name:      "process single command",
			cmds:      []string{"echo helloWorld"},
			fileNames: []string{"echo_helloWorld"},
			outputs:   []string{"helloWorld"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proc := NewCommandProcessor(test.cmds)
			err := proc.Process()
			if err != nil {
				if !test.shouldFail {
					t.Fatalf("unexpected error: %s", err)
				}
				t.Log(err)
				return
			}

			if _, err := os.Stat(proc.outputPath); err != nil {
				t.Fatalf("unable to Stat %s", proc.outputPath)
			}

			// validate content
			tarf, err := os.Open(proc.outputPath)
			if err != nil {
				t.Fatal(err)
			}

			var buf bytes.Buffer
			tarrer := tar.NewReader(tarf)
			i := 0
			for {
				hdr, err := tarrer.Next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						t.Fatal(err)
					}
				}
				if hdr.Name != test.fileNames[i] {
					t.Errorf("unexpected name %s", hdr.Name)
				}
				if _, err := io.Copy(&buf, tarrer); err != nil {
					t.Fatal(err)
				}
				if buf.String() != test.outputs[i] {
					t.Errorf("unexpected content from archiver for cmd %s", test.cmds[i])
				}
				i++
			}
		})
	}
}
