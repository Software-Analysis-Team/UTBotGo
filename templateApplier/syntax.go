package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Syntax map[Name]string

type FileSyntaxes struct{
	ForBase map[string]Syntax
	ForExt  map[string]Syntax
}

var TASDsyntax = Syntax{
	"syntaxBlock": `for (<file>|<extension>):\n<syntaxDescr>`,
	"file":        `<name>`,
	"extension":   `\*<name>`,
	"name":        `(\w|\.)+`,
	"syntaxDescr": `<*unitDescr>(.|\n\s)+\n`,
	"unitDescr":   `\s+<unitName>:\s+<unitBody>\n`,
	"unitName":    `\w+`,
	"unitBody":    `.*`,
}

func SearchTASDfile(wd string) (string, error) {
	var (
		files []os.DirEntry
		err   error
	)
	for wd != "." {
		files, err = os.ReadDir(wd)
		if err != nil {
			return "", err
		}
		for _, file := range files {
			if file.Name() == ".tasd" {
				return filepath.Join(wd, ".tasd"), nil
			}
		}
	}
	return "", fmt.Errorf("file .tasd was not found")
}

func ReadTASDfile(fileName string) (FileSyntaxes, error) {
	src, err := os.ReadFile(fileName)
	if err != nil {
		return FileSyntaxes{}, err
	}
	return ParseTASDfile(src), nil
}

func ParseTASDfile(src []byte) (result FileSyntaxes) {
	result.ForBase = make(map[string]Syntax)
	result.ForExt  = make(map[string]Syntax)
	r := NewReplacer(src, TASDsyntax)
	for _, rSyntaxBlock := range r.GetAll("syntaxBlock") {
		rSyntaxDescr := rSyntaxBlock.Get("syntaxDescr")
		newSyntax := make(Syntax)
		for _, rUnitDescr := range rSyntaxDescr.GetAll("unitDescr") {
			unitName :=   Name(rUnitDescr.Get("unitName").GetString())
			unitBody := string(rUnitDescr.Get("unitBody").GetString())
			newSyntax[unitName] = unitBody
		}
		name := string(rSyntaxBlock.Get("name").GetString())
		if rSyntaxBlock.ThereIs("file") {
			result.ForBase[name] = newSyntax
		} else {
			result.ForExt[name]  = newSyntax
		}
	}
	return
}
