package script

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// wordSplit splits space-separted strings into words including quoted words:
//
//     aaa "bbb" "ccc ddd" eee
//
//  In case of aaa"abcd", the whole thing is returned as aaa"abcd" including qoutes.
//  In case of "aaa"bbb will be returned as two words "aaa" and "bbb"
func wordSplit(val string) ([]string, error) {
	rdr := bufio.NewReader(strings.NewReader(val))
	var startQuote rune
	var word strings.Builder
	words := make([]string, 0)
	inWord := false
	inQuote := false
	squashed := false

	for {
		token, _, err := rdr.ReadRune()
		if err != nil {
			if err == io.EOF {
				remainder := word.String()
				if len(remainder) > 0 {
					words = append(words, remainder)
				}
				return words, nil
			}
			return nil, err
		}

		switch {
		case isChar(token):
			if !inWord {
				inWord = true
			}
			word.WriteRune(token)

		case isQuote(token):
			if !inWord {
				inWord, inQuote = true, true
				startQuote = token
				continue
			}

			// handles case when unquoted runs into quoted: abc"defg"
			// start the quote here
			if inWord && !inQuote {
				inQuote, squashed = true, true
				startQuote = token
				word.WriteRune(token)
				continue
			}

			// handle embedded quote (i.e "'aa'")
			if inWord && inQuote && token != startQuote {
				word.WriteRune(token)
				continue
			}

			// capture closing quote when in abc"defg"
			if squashed {
				word.WriteRune(token)
			}

			inWord = false
			inQuote = false
			squashed = false
			//store
			words = append(words, word.String())
			word.Reset()

		case unicode.IsSpace(token):
			if !inWord {
				inWord = false
				continue
			}

			// capture quoted space
			if inWord && inQuote {
				word.WriteRune(token)
				continue
			}

			// end of word
			inWord = false
			words = append(words, word.String())
			word.Reset()
		}
	}
}

func isQuote(r rune) bool {
	switch r {
	case '"', '\'':
		return true
	}
	return false
}

func isChar(r rune) bool {
	return !isQuote(r) && !unicode.IsSpace(r)
}

func quote(str string) string {
	if strings.Index(str, `'`) > -1 {
		return doubleQuote(str)
	}
	if strings.Index(str, `"`) > -1 {
		return singleQuote(str)
	}
	return doubleQuote(str)
}

func doubleQuote(val string) string {
	return fmt.Sprintf(`"%s"`, val)
}

func singleQuote(val string) string {
	return fmt.Sprintf(`'%s'`, val)
}

func isQuoted(val string) bool {
	single := `'`
	dbl := `"`
	if strings.HasPrefix(val, single) && strings.HasSuffix(val, single) {
		return true
	}
	if strings.HasPrefix(val, dbl) && strings.HasSuffix(val, dbl) {
		return true
	}
	return false
}

func trimQuotes(val string) string {
	single := `'`
	dbl := `"`

	if strings.HasPrefix(val, single) || strings.HasPrefix(val, dbl) {
		val = strings.TrimPrefix(val, val[0:1])
	}
	if strings.HasSuffix(val, single) || strings.HasSuffix(val, dbl) {
		val = strings.TrimSuffix(val, val[len(val)-1:len(val)])
	}

	return val
}
