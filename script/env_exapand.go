// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"
)

// runeStack is a super simple stack (with slice backing) used to keep track
// of special boundary characters.
type runeStack struct {
	store []rune
	top   int
}

func newRuneStack() *runeStack {
	return &runeStack{store: []rune{}, top: -1}
}

func (r *runeStack) push(val rune) {
	r.top++
	if r.top > len(r.store)-1 {
		r.store = append(r.store, val)
	} else {
		r.store[r.top] = val
	}
}

func (r *runeStack) pop() rune {
	if r.isEmpty() {
		return 0
	}
	val := r.store[r.top]
	r.top--
	return val
}

func (r *runeStack) peek() rune {
	if r.isEmpty() {
		return 0
	}
	return r.store[r.top]
}

func (r *runeStack) isEmpty() bool {
	return (r.top < 0)
}

func (r *runeStack) depth() int {
	return r.top + 1
}

// ExpandEnv searches str for $value or ${value} which is then evaluated
// using os.ExpandEnv. expandVar supports escaping expansion using \$. For instance,
// when \$value or \${value} is encountered, it is not expanded, leaving the original
// values in the string as $value or ${value}.
func ExpandEnv(str string) string {
	stack := newRuneStack()
	rdr := bufio.NewReader(strings.NewReader(str))
	var result strings.Builder
	var variable strings.Builder

	inVar := false
	//inEscape := false

	// The algorithm is simple:
	// a) when boundary char \ or $ is encountered: push onto stack
	// b) (escape) if stack.top = \ and $ is encountered, skip slash, pop all items and $ unto result string
	// c) if in scape write all subsequent chars in result (except \ prefix) until/including space char or end of string
	// d) if inVar ($ followed by nonspace), save all subsequent char in variable until a space char or end of string
	for {
		token, _, err := rdr.ReadRune()
		if err != nil {
			// resolve outstanding vars and save dangling slashes/dollar signs at EOF
			if err == io.EOF {
				popAll(&result, stack)
				if inVar {
					result.WriteString(resloveVar(&variable))
				}
			}
			return result.String()
		}

		switch {
		// save '\' on stack for later
		case isBackSlash(token):
			stack.push(token)

		// if '$':
		// 1) if stack.top = '\', escape
		// 2) else save on stack
		case isDollarSign(token):
			if isBackSlash(stack.peek()) {
				stack.pop()
				popAll(&result, stack)
				result.WriteRune(token)
				continue
			}
			stack.push(token)

		// if '{':
		// if stack.top = '$', start of ${variable}
		// else write token unto result
		case isOpenCurly(token):
			if isDollarSign(stack.peek()) {
				inVar = true
				variable.WriteRune(stack.pop())
				popAll(&result, stack)
				variable.WriteRune(token)
				continue
			}
			result.WriteRune(token)

		// handle all other chars
		default:
			switch {
			// if '}':
			// if in var, assume var boundary, resolve/save var in result str
			// else, save token in result srt
			case isCloseCurly(token):
				if inVar {
					inVar = false
					variable.WriteRune(token)
					result.WriteString(resloveVar(&variable))
					continue
				}
				result.WriteRune(token)

			// if token is boundary (space, punctuations, symbols, etc):
			// 1) if in var, assume var boundary resolve/save var in result str
			// 2) or, write tokens to result string
			case isBoundary(token):
				if inVar {
					inVar = false
					result.WriteString(resloveVar(&variable))
					result.WriteRune(token)
					continue
				}
				popAll(&result, stack)
				result.WriteRune(token)

			// if letter:
			// 1) if in var, save var name
			// 2) if stack.top = '$', assume start of var
			// 3) otherwise write token in result
			default:
				if inVar {
					variable.WriteRune(token)
					continue
				}

				if isDollarSign(stack.peek()) {
					inVar = true
					variable.WriteRune(stack.pop())
					variable.WriteRune(token)
					continue
				}

				popAll(&result, stack)
				result.WriteRune(token)
			}
		}
	}
}

func isDollarSign(r rune) bool {
	if r == '$' {
		return true
	}
	return false
}

func isBackSlash(r rune) bool {
	if r == '\\' {
		return true
	}
	return false
}

func isOpenCurly(r rune) bool {
	if r == '{' {
		return true
	}
	return false
}
func isCloseCurly(r rune) bool {
	if r == '}' {
		return true
	}
	return false
}
func popAll(target *strings.Builder, stack *runeStack) {
	for !stack.isEmpty() {
		target.WriteRune(stack.pop())
	}
}

func resloveVar(variable *strings.Builder) string {
	val := variable.String()
	variable.Reset()
	return os.ExpandEnv(val)
}

func isBoundary(token rune) bool {
	switch {
	case unicode.IsSpace(token), token == ':', token == '#', token == '%':
		return true
	}
	return false
}
