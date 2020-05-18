package script

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	spaceSep       = regexp.MustCompile(`\s`)
	paramSep       = regexp.MustCompile(`:`)
	quoteSet       = regexp.MustCompile(`[\"\']`)
	cmdSep         = regexp.MustCompile(`\s`)
	namedParamRegx = regexp.MustCompile(`^([a-z0-9_\-]+)(:)(["']{0,1}.+["']{0,1})$`)
)

// mapArgs takes the rawArgs in the form of
//
//    param0:"val0" param1:"val1" ... paramN:"valN"
//
// The param name must be followed by a colon and the value
// may be quoted or unquoted. It is an error if
// split(rawArgs[n], ":") yields to a len(slice) < 2.
func mapArgs(rawArgs string) (map[string]string, error) {
	argMap := make(map[string]string)

	// split params: param0:<param-val0> paramN:<param-valN> badparam
	params, err := commandSplit(rawArgs)
	if err != nil {
		return nil, err
	}

	// for each, split pram:<pram-value> into {param, <param-val>}
	for _, param := range params {
		cmdName, cmdStr, err := namedParamSplit(param)
		if err != nil {
			return nil, fmt.Errorf("map args: %s", err)
		}
		argMap[cmdName] = cmdStr
	}

	return argMap, nil
}

// isNamedParam returs true if str has the form
//
//    name:value
//
func isNamedParam(str string) bool {
	return namedParamRegx.MatchString(str)
}

// makeParam
func makeNamedPram(name, value string) string {
	value = strings.TrimSpace(value)
	// possibly already quoted
	if value[0] == '"' || value[0] == '\'' {
		return fmt.Sprintf("%s:%s", name, value)
	}
	// return as quoted
	return fmt.Sprintf(`%s:'%s'`, name, value)
}

func validateRawArgs(cmdName, rawArgs string) error {
	cmd, ok := Cmds[cmdName]
	if !ok {
		return fmt.Errorf("%s is unknown", cmdName)
	}
	if len(rawArgs) == 0 && cmd.MinArgs > 0 {
		return fmt.Errorf("%s must have at least %d argument(s)", cmdName, cmd.MinArgs)
	}
	return nil
}

func validateCmdArgs(cmdName string, args map[string]string) error {
	cmd, ok := Cmds[cmdName]
	if !ok {
		return fmt.Errorf("%s is unknown", cmdName)
	}

	minArgs := cmd.MinArgs
	maxArgs := cmd.MaxArgs

	if len(args) < minArgs {
		return fmt.Errorf("%s must have at least %d argument(s)", cmdName, minArgs)
	}

	if maxArgs > -1 && len(args) > maxArgs {
		return fmt.Errorf("%s can only have up to %d argument(s)", cmdName, maxArgs)
	}

	return nil
}