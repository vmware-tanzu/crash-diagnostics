package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func TestExecutorNew(t *testing.T) {
	e := NewCommandExecutor("foo", []string{"bar", "baz"}...)
	if &e.command == nil {
		t.Error("command not set")
	}
}

func TestExecutorExec(t *testing.T) {
	tests := []struct {
		name           string
		cmd            string
		args           []string
		expectedOutput string
		shouldFail     bool
	}{
		{
			name:           "normal exec",
			cmd:            "echo",
			args:           []string{"HelloWorld"},
			expectedOutput: "HelloWorld",
		},
		{
			name:       "missing cmd",
			cmd:        "foobar",
			shouldFail: true,
		},
	}

	for _, test := range tests {
		exec := NewCommandExecutor(test.cmd, test.args...)
		result, err := exec.Execute()
		if err != nil {
			if !test.shouldFail {
				t.Fatalf("unexpected error: %s", err)
			}
			t.Log(err)
			continue
		}

		output := new(bytes.Buffer)
		io.Copy(output, result)
		if strings.TrimSpace(output.String()) != test.expectedOutput {
			t.Errorf("unexpected Executor.Exec() result: %s", result)
		}
	}
}

func TestExecutorExecToFile(t *testing.T) {
	tests := []struct {
		name           string
		cmd            string
		args           []string
		fileName       string
		expectedOutput string
		shouldFail     bool
	}{
		{
			name:           "normal exec to file",
			cmd:            "echo",
			args:           []string{"HelloWorld"},
			fileName:       fmt.Sprintf("/tmp/%x.txt", rand.Int31()),
			expectedOutput: "HelloWorld",
		},
		{
			name:       "bad exec to file",
			cmd:        "foo",
			fileName:   fmt.Sprintf("/tmp/%x.txt", rand.Int31()),
			shouldFail: true,
		},
	}

	for _, test := range tests {
		exec := NewCommandExecutor(test.cmd, test.args...)
		err := exec.ExecToFile(test.fileName)
		if err != nil {
			if !test.shouldFail {
				t.Fatalf("unexpected error: %s", err)
			}
			t.Log(err)
			continue
		}

		if _, err := os.Stat(test.fileName); err != nil {
			t.Fatal(err)
		}
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
			t.Errorf("unexpected Executor.Exec() result: %s", result)
		}
	}
}
