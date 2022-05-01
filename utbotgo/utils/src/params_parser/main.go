package main

import (
	"fmt"
	"os"
)

func GetFunctions(filePaths []string, functionNames map[string]struct{}) (functions []GoFunction, err error) {
	var goFile GoFile
	for _, filePath := range filePaths {
		goFile, err = ParseGoFile(filePath)
		if err != nil {
			return
		}
		for _, function := range goFile.Functions() {
			if _, ok := functionNames[function.Name()]; ok {
				functions = append(functions, function)
			}
		}
	}
	return
}

func throw(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	config, helpMessage, err := ParseArgs(os.Args)
	throw(err)
	if helpMessage != ""{
		fmt.Print(helpMessage)
		return
	}
	functions, err := GetFunctions(config.GoFilePaths, config.FunctionNames)
	throw(err)
	if config.TaArgsOutputTemplate != "" {
		err = PrintFunctionTaInfos(functions, config.TaArgsOutputTemplate)
		throw(err)
	}
	if config.KleeArgsOutputTemplate != "" {
		err = PrintFunctionKleeInfos(functions, config.KleeArgsOutputTemplate)
		throw(err)
	}
}
