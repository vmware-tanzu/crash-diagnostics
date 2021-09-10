package script_tests

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/crash-diagnostics/exec"
	"github.com/vmware-tanzu/crash-diagnostics/functions/providers"
	"github.com/vmware-tanzu/crash-diagnostics/typekit"
)

func TestHostlistProviderScript(t *testing.T) {
	tests := []struct {
		name   string
		script string
		eval   func(*testing.T, string)
	}{
		{
			name:   "simple script",
			script: `result=hostlist_provider(hosts=["127.0.0.1", "localhost"])`,
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}

				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("hostlist_provider() should be assigned to a variable for testing")
				}
				var result providers.Result
				if err := typekit.Starlark(resultVal).Go(&result); err != nil {
					t.Fatal(err)
				}
				if len(result.Resources.Hosts) != 2 {
					t.Errorf("unexpected host count %d", len(result.Resources.Hosts))
				}
				for i := range result.Resources.Hosts {
					if result.Resources.Hosts[i] != "127.0.0.1" && result.Resources.Hosts[i] != "localhost" {
						t.Errorf("unexpected resource hosts values %s", result.Resources.Hosts[i])
					}
				}
			},
		},
		{
			name:   "resources with hostlist_provider",
			script: `result=resources(provider=hostlist_provider(hosts=["127.0.0.1", "localhost"]))`,
			eval: func(t *testing.T, script string) {
				output, err := exec.Run("test.star", strings.NewReader(script), nil)
				if err != nil {
					t.Fatal(err)
				}

				resultVal := output["result"]
				if resultVal == nil {
					t.Fatal("resources() should be assigned to a variable for testing")
				}
				var resources providers.Resources
				if err := typekit.Starlark(resultVal).Go(&resources); err != nil {
					t.Fatal(err)
				}
				if len(resources.Hosts) != 2 {
					t.Errorf("unexpected host count %d", len(resources.Hosts))
				}
				for i := range resources.Hosts {
					if resources.Hosts[i] != "127.0.0.1" && resources.Hosts[i] != "localhost" {
						t.Errorf("unexpected resource hosts values %s", resources.Hosts[i])
					}
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
