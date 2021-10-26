package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	makefileText, _ := os.ReadFile(os.Args[1])
	var values map[string]string
	switch res := ScanBytes(makefileText, NewMakefileScanner()).(type) {
	case MakefileData:
		values = res.Values
	}
	inputText, _ := io.ReadAll(os.Stdin)
	switch res := ScanBytes(inputText, NewTemplateScanner(values)).(type) {
	case TemplateData:
		fmt.Print(res.Output.String())
	}
}
