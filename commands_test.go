package main

import (
	"strings"
	"testing"
)

func TestCommands_Parse(t *testing.T) {
	tests := []struct {
		name       string
		commands   string
		preambles  map[string]command
		actions    []command
		shouldFail bool
	}{
		{
			name:      "single preamble",
			commands:  "FROM default",
			preambles: map[string]command{"FROM": command{index: 1, name: "FROM", args: []string{"default"}}},
		},
		{
			name:     "single action",
			commands: "COPY a",
			actions:  []command{{index: 1, name: "COPY", args: []string{"a"}}},
		},
		{
			name:      "multiple commands",
			commands:  "FROM default\nCOPY a",
			preambles: map[string]command{"FROM": command{index: 1, name: "FROM", args: []string{"default"}}},
			actions:   []command{{index: 2, name: "COPY", args: []string{"a"}}},
		},
		{
			name:       "single unsupported command",
			commands:   "FOO default",
			shouldFail: true,
		},
		{
			name:       "multiple with unsupported command",
			commands:   "FOO default\nCOPY /abc /edf",
			shouldFail: true,
		},
		{
			name:       "single low case command",
			commands:   "foo default",
			shouldFail: true,
		},
		{
			name:       "multiple with low case command",
			commands:   "From default\nCOPY /abc /edf",
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script, err := parse(strings.NewReader(test.commands))
			if err != nil {
				if !test.shouldFail {
					t.Fatal(err)
				}
				t.Log(err)
				return
			}
			preambles := script.preambles
			actions := script.actions
			if len(test.preambles) != len(preambles) {
				t.Fatalf("expecting %d preambles, got %d", len(test.preambles), len(preambles))
			}
			for name, cmd := range test.preambles {
				if preambles[name] == nil {
					t.Errorf("missing expected preamble %s", name)
				}
				if cmd.index != preambles[name].index {
					t.Errorf("%s preamble index mismatched: %d != %d", name, cmd.index, preambles[name].index)
				}
				if len(cmd.args) != len(preambles[name].args) {
					t.Errorf("%s preamble args mismatched: %d != %d", name, len(cmd.args), len(preambles[name].args))
				}
			}

			if len(test.actions) != len(actions) {
				t.Fatalf("expecting %d actions, got %d", len(test.actions), len(actions))
			}

			for i := range test.actions {
				if test.actions[i].index != actions[i].index {
					t.Errorf("expecting command index: %v, got: %v", test.actions[i].index, actions[i].index)
				}
				if test.actions[i].name != actions[i].name {
					t.Errorf("expecting name: %v, got: %v", test.actions[i].name, actions[i].name)
				}
				if len(test.actions[i].args) != len(actions[i].args) {
					t.Errorf("expecting args count: %d, got: %d", len(test.actions[i].args), len(actions[i].args))
				}
			}
		})
	}
}
