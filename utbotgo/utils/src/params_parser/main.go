package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetPackageFiles(configPath string) (filePaths []string, err error) {
	configText, err := os.ReadFile(configPath)
	if err != nil {
		return
	}
	var config map[string]interface{}
	err = json.Unmarshal(configText, &config)
	if err != nil {
		return
	}
	goFilesAsInterface, ok := config["GoFiles"]
	if ! ok {
		err = fmt.Errorf("key \"GoFiles\" not found")
		return
	}
	goFilesAsArray, ok := goFilesAsInterface.([]interface{})
	if ! ok {
		err = fmt.Errorf("value of key \"GoFiles\" must be array of strings")
		return
	}
	filePaths = make([]string, len(goFilesAsArray))
	for i, goFile := range goFilesAsArray {
		filePaths[i], ok = goFile.(string)
		if ! ok {
			err = fmt.Errorf("value of key \"GoFiles\" must be array of strings")
			return
		}
	}
	return
}

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
	filePaths, err := GetPackageFiles(config.PackageConfigPath)
	throw(err)
	functions, err := GetFunctions(filePaths, config.FunctionNames)
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
