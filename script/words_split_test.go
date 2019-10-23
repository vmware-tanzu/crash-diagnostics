package script

import "testing"

func TestWordSplit(t *testing.T) {
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
			words, err := wordSplit(test.str)
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

func TestWordSplitTrimQuotes(t *testing.T) {
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
