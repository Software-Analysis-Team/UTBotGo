package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type ArgTaInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ResultTaInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FunctionTaInfo struct {
	Name string `json:"name"`
	Args []ArgTaInfo `json:"args"`
	Results []ResultTaInfo `json:"results"`
}

func GetFunctionTaInfo(function GoFunction) (FunctionTaInfo, bool) {
	name := function.Name()
	params, err := function.Params()
	if err != nil {
		return FunctionTaInfo{}, false
	}
	args := make([]ArgTaInfo, len(params))
	for i := range params {
		args[i] = ArgTaInfo{
			Name: params[i].Name,
			Type: params[i].Type.StringAsGoType(),
		}
	}
	results_, err := function.Results()
	if err != nil {
		return FunctionTaInfo{}, false
	}
	results := make([]ResultTaInfo, len(results_))
	for i := range results_ {
		results[i] = ResultTaInfo{
			Name: results_[i].Name,
			Type: results_[i].Type.StringAsGoType(),
		}
	}
	return FunctionTaInfo{
		Name: name,
		Args: args,
		Results: results,
	}, true
}

func PrintFunctionTaInfos(functions []GoFunction, taArgsOutputTemplate string) (err error) {
	var functionInfos []FunctionTaInfo
	for _, function := range functions {
		functionInfo, ok := GetFunctionTaInfo(function)
		if ! ok {
			continue
		}
		functionInfos = append(functionInfos, functionInfo)
	}
	if strings.Count(taArgsOutputTemplate, "%") == 1 {
		for i, functionInfo := range functionInfos {
			functionName := functions[i].Name()
			fileName := strings.Replace(taArgsOutputTemplate, "%", functionName, 1)
			var jsonText []byte
			jsonText, err = json.Marshal(functionInfo)
			if err != nil {
				return
			}
			err = os.MkdirAll(path.Dir(fileName), 0777)
			if err != nil {
				return
			}
			err = os.WriteFile(fileName, jsonText, 0666)
			if err != nil {
				return
			}
		}
	} else {
		var jsonText []byte
		jsonText, err = json.Marshal(functionInfos)
		if err != nil {
			return
		}
		switch taArgsOutputTemplate {
		case "-":
			_, err = fmt.Println(string(jsonText))
		case "!":
			_, err = fmt.Fprintln(os.Stderr, string(jsonText))
		default:
			err = os.MkdirAll(path.Dir(taArgsOutputTemplate), 0777)
			if err != nil {
				return
			}
			err = os.WriteFile(taArgsOutputTemplate, jsonText, 0666)
		}
	}
	return
}
