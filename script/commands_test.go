package script

import (
	"strings"
	"testing"
)

func TestCommands_Parse(t *testing.T) {
	tests := []struct {
		name       string
		commands   string
		preambles  map[string][]Command
		actions    []Command
		shouldFail bool
	}{
		{
			name:      "single preamble",
			commands:  "FROM default",
			preambles: map[string][]Command{"FROM": {{Index: 1, Name: "FROM", Args: []string{"default"}}}},
		},
		// {
		// 	name:     "single action",
		// 	commands: "COPY a",
		// 	actions:  []Command{{Index: 1, Name: "COPY", Args: []string{"a"}}},
		// },
		// {
		// 	name:      "multiple commands",
		// 	commands:  "FROM default\nCOPY a",
		// 	preambles: map[string]Command{"FROM": Command{Index: 1, Name: "FROM", Args: []string{"default"}}},
		// 	actions:   []Command{{Index: 2, Name: "COPY", Args: []string{"a"}}},
		// },
		// {
		// 	name:       "single unsupported command",
		// 	commands:   "FOO default",
		// 	shouldFail: true,
		// },
		// {
		// 	name:       "multiple with unsupported command",
		// 	commands:   "FOO default\nCOPY /abc /edf",
		// 	shouldFail: true,
		// },
		// {
		// 	name:       "single low case command",
		// 	commands:   "foo default",
		// 	shouldFail: true,
		// },
		// {
		// 	name:       "multiple with low case command",
		// 	commands:   "From default\nCOPY /abc /edf",
		// 	shouldFail: true,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script, err := Parse(strings.NewReader(test.commands))
			if err != nil {
				if !test.shouldFail {
					t.Fatal(err)
				}
				t.Log(err)
				return
			}
			preambles := script.Preambles
			actions := script.Actions
			if len(test.preambles) != len(preambles) {
				t.Fatalf("expecting %d preambles, got %d", len(test.preambles), len(preambles))
			}
			for name, cmds := range test.preambles {
				if preambles[name] == nil {
					t.Errorf("missing expected preamble %s", name)
				}
				if len(cmds) != len(preambles[name]) {
					t.Errorf("unexpected number directives for preamble %s: expecting %d got %d", name, len(preambles[name]), len(cmds))
				}
				for i, cmd := range cmds {
					if cmd.Index != preambles[name][i].Index {
						t.Errorf("%s preamble index mismatched: %d != %d", name, cmd.Index, preambles[name][i].Index)
					}
					if len(cmd.Args) != len(preambles[name][i].Args) {
						t.Errorf("%s preamble args mismatched: %d != %d", name, len(cmd.Args), len(preambles[name][i].Args))
					}
				}
			}

			if len(test.actions) != len(actions) {
				t.Fatalf("expecting %d actions, got %d", len(test.actions), len(actions))
			}

			for i := range test.actions {
				if test.actions[i].Index != actions[i].Index {
					t.Errorf("expecting command index: %v, got: %v", test.actions[i].Index, actions[i].Index)
				}
				if test.actions[i].Name != actions[i].Name {
					t.Errorf("expecting name: %v, got: %v", test.actions[i].Name, actions[i].Name)
				}
				if len(test.actions[i].Args) != len(actions[i].Args) {
					t.Errorf("expecting args count: %d, got: %d", len(test.actions[i].Args), len(actions[i].Args))
				}
			}
		})
	}
}
