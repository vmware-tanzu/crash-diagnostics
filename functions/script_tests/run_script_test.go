package script_tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"github.com/vmware-tanzu/crash-diagnostics/functions/run"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

func TestRunScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(*testing.T, string)
	}{
		{
			name: "run command",
			script: fmt.Sprintf(`
result=run(
    cmd="""echo 'Hello World!'""",
    ssh_config=make_ssh_config(
       username="%s",
       port="%s",
       private_key_path="%s",
       max_retries=%d,
    ).config,
    resources=hostlist_provider(hosts=["127.0.0.1"]).resources,
)`, testSupport.CurrentUsername(), testSupport.PortValue(), testSupport.PrivateKeyPath(), testSupport.MaxConnectionRetries()),
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}
				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("run() should be assigned to a variable for test")
				}
				var result run.Result
				if err := typekit.Starlark(resultVal).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatalf("command failed: %s", result.Error)
				}
				if len(result.Procs) != 1 {
					t.Fatal("missing command result")
				}
				expected := "Hello World!"
				out := strings.TrimSpace(result.Procs[0].Output)
				if out != expected {
					t.Error("unexpected result:", output)
				}
			},
		},

		{
			name: "run command with aliases",
			script: fmt.Sprintf(`
result=run(
    cmd="""echo 'Hello Starlark!'""",
    ssh_config=ssh_config(
       username="%s",
       port="%s",
       private_key_path="%s",
       max_retries=%d,
    ),
    resources=resources(provider=hostlist_provider(hosts=["127.0.0.1"])),
)`, testSupport.CurrentUsername(), testSupport.PortValue(), testSupport.PrivateKeyPath(), testSupport.MaxConnectionRetries()),
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}
				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("run() should be assigned to a variable for test")
				}
				var result run.Result
				if err := typekit.Starlark(resultVal).Go(&result); err != nil {
					t.Fatal(err)
				}
				if result.Error != "" {
					t.Fatalf("command failed: %s", result.Error)
				}
				if len(result.Procs) != 1 {
					t.Fatal("missing command result")
				}
				expected := "Hello Starlark!"
				out := strings.TrimSpace(result.Procs[0].Output)
				if out != expected {
					t.Error("unexpected result:", output)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.script)
		})
	}
}
