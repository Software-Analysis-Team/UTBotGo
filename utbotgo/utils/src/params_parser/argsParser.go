package main

import (
	"fmt"
	"strings"
)

type Config struct {
	GoFilePaths            []string
	TaArgsOutputTemplate   string
	KleeArgsOutputTemplate string
	FunctionNames          map[string]struct{}
}

func ParseArgs(args []string) (config Config, helpMessage string, err error) {
	config.FunctionNames = make(map[string]struct{})
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--go-files":
			i++
			if i == len(args) {
				err = fmt.Errorf("file after flag '--go-files' is requred")
				return
			}
			config.GoFilePaths = strings.Split(args[i], " ")
		case "--klee":
			i++
			if i == len(args) {
				err = fmt.Errorf("template after flag '--klee' is requred")
				return
			}
			config.KleeArgsOutputTemplate = args[i]
		case "--ta":
			i++
			if i == len(args) {
				err = fmt.Errorf("template after flag '--ta' is requred")
				return
			}
			config.TaArgsOutputTemplate = args[i]
		default:
			if args[i][0] == '-' {
				helpMessage = HelpMessage
				return
			}
			config.FunctionNames[args[i]] = struct{}{}
		}
	}
	if config.GoFilePaths == nil {
		err = fmt.Errorf("option '--go-files' is required")
	}
	return
}

var HelpMessage = `
Usage: params_parser [options] functionNames...

Description:
    'params_parser' print information about specified functions in current Go
    package in JSON format.

Arguments:
    --go-files <file names> (required option)
        Go-files in current package
    --klee <template>
        template for file, in which will print information about specified
        functions in JSON format for command 'ktest-tool'
    --ta <template>
        as option '--klee', but for command 'ta'
    <functionName>
        name of Go function, about which information is needed

Template:
    Template is string, specified file(s) for printing information. There are
    some cases:
      - if template is "-", information about all specified functions will be
        printed on standard output stream
      - if template is "!", information about all specified functions will be
        printed on standard error stream
      - if template contain one char '%', for all specified functions'
        information about it will be printed in separated file with name
        obtained by substitution function name to the place of char '%'
      - otherwise information about all specified functions will be printed to
        file with template name
`[1:]
