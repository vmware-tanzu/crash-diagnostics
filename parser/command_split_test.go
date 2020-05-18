package parser

import "testing"

func TestCommandSplit(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		words []string
	}{
		{
			name:  "no quotes",
			str:   `aaa bbb ccc ddd`,
			words: []string{"aaa", "bbb", "ccc", "ddd"},
		},
		{
			name:  "all quotes",
			str:   `"aaa" "bbb" "ccc" "ddd"`,
			words: []string{"aaa", "bbb", "ccc", "ddd"},
		},
		{
			name:  "mix unquoted quoted",
			str:   `aaa "bbb" "ccc ddd"`,
			words: []string{"aaa", "bbb", "ccc ddd"},
		},
		{
			name:  "mix quoted unquoted",
			str:   `"aaa" "bbb ccc" ddd`,
			words: []string{"aaa", "bbb ccc", "ddd"},
		},
		{
			name:  "front quote runin",
			str:   `aaa"bbb ccc" ddd`,
			words: []string{"aaa\"bbb ccc\"", "ddd"},
		},
		{
			name:  "back quote runin",
			str:   `aaa "bbb ccc"ddd`,
			words: []string{"aaa", "bbb ccc", "ddd"},
		},
		{
			name:  "embedded single quotes",
			str:   `aaa "'bbb' ccc" ddd`,
			words: []string{"aaa", "'bbb' ccc", "ddd"},
		},
		{
			name:  "embedded double quotes",
			str:   `'aaa' '"bbb ccc"' ddd`,
			words: []string{"aaa", `"bbb ccc"`, "ddd"},
		},
		{
			name:  "embedded double quotes runins",
			str:   `aaa'"bbb ccc"' ddd`,
			words: []string{`aaa'"bbb ccc"'`, "ddd"},
		},
		{
			name:  "embedded single quotes runins",
			str:   `aaa"bbb 'ccc'" ddd`,
			words: []string{`aaa"bbb 'ccc'"`, "ddd"},
		},
		{
			name:  "actual exec command",
			str:   `/bin/bash -c 'echo "Hello World"'`,
			words: []string{`/bin/bash`, `-c`, `echo "Hello World"`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			words, err := commandSplit(test.str)
			if err != nil {
				t.Error(err)
			}
			if len(words) != len(test.words) {
				t.Fatalf("unexpected length: want %#v, got %#v", test.words, words)
			}
			for i := range words {
				if words[i] != test.words[i] {
					t.Errorf("word mistached:\ngot %#v\nwant %#v", words, test.words)
				}
			}
		})
	}
}

func TestCommandSplitTrimQuotes(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		result string
	}{
		{
			name:   "balanced double quote",
			str:    `"aa bb cc dd"`,
			result: "aa bb cc dd",
		},
		{
			name:   "balanced single quote",
			str:    `'aa bb cc dd'`,
			result: "aa bb cc dd",
		},
		{
			name:   "balanced double with embedded single",
			str:    `"aa 'bb cc' dd"`,
			result: "aa 'bb cc' dd",
		},
		{
			name:   "balanced single with embedded double",
			str:    `'"aa bb" cc dd'`,
			result: `"aa bb" cc dd`,
		},
		{
			name:   "balanced single with embedded singles",
			str:    `''aa bb cc' dd'`,
			result: `'aa bb cc' dd`,
		},
		{
			name:   "unbalanced singles with embedded singles",
			str:    `aa bb cc' dd'`,
			result: `aa bb cc' dd`,
		},
		{
			name:   "unbalanced singles with embedded doubles",
			str:    `'aa "bb cc" dd`,
			result: `aa "bb cc" dd`,
		},
		{
			name:   "unbalanced double with embedded singles",
			str:    `aa 'bb cc' dd"`,
			result: `aa 'bb cc' dd`,
		},
		{
			name:   "unbalanced double with embedded doubles",
			str:    `"aa "bb cc" dd`,
			result: `aa "bb cc" dd`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := trimQuotes(test.str)
			if result != test.result {
				t.Fatalf("unexpected result: want %v, got %v", test.result, result)
			}
		})
	}
}

func TestNamedParamSplit(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		parts []string
	}{
		{
			name:  "no quotes",
			str:   `cmd:name:value`,
			parts: []string{"cmd", "name:value"},
		},
		{
			name:  "single quotes",
			str:   `cmd:'name:single-quote-value'`,
			parts: []string{"cmd", "name:single-quote-value"},
		},
		{
			name:  "double quotes",
			str:   `cmd:"name: double-quote-value"`,
			parts: []string{"cmd", "name: double-quote-value"},
		},
		{
			name:  "mismatch quotes",
			str:   `cmd:'name:mismatch-quote-value"`,
			parts: []string{"cmd", "name:mismatch-quote-value"},
		},
		{
			name:  "unbalanced quotes",
			str:   `cmd:'unbalanced-quote:value`,
			parts: []string{"cmd", "unbalanced-quote:value"},
		},
		{
			name:  "malformed param",
			str:   `cmd:'malformed-param' cmd:abc`,
			parts: []string{"cmd", "malformed-param' cmd:abc"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, val, err := namedParamSplit(test.str)
			if err != nil {
				t.Error(err)
			}
			if test.parts[0] != name {
				t.Fatalf("expecting param name %s, got %s", test.parts[0], name)
			}
			if test.parts[1] != val {
				t.Fatalf("expecting param value [%s], got [%s]", test.parts[1], val)
			}
		})
	}
}
