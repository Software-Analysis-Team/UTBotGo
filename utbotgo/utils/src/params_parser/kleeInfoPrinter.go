package main
import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type ArgKleeInfos map[string]string

type FunctionKleeInfo struct {
	Name string `json:"name"`
	Args ArgKleeInfos `json:"args"`
}

func GetArgKleeInfo(function GoFunction) (ArgKleeInfos, bool) {
	params, err := function.Params()
	if err != nil {
		return nil, false
	}
	args := make(ArgKleeInfos)
	for i := range params {
		args[params[i].Name] = params[i].Type.StringAsKleeType()
	}
	return args, true
}

func PrintFunctionKleeInfos(functions []GoFunction, taArgsOutputTemplate string) (err error) {
	var argInfos []ArgKleeInfos
	for _, function := range functions {
		argInfo, ok := GetArgKleeInfo(function)
		if ! ok {
			continue
		}
		argInfos = append(argInfos, argInfo)
	}
	if strings.Count(taArgsOutputTemplate, "%") == 1 {
		for i, argInfo := range argInfos {
			functionName := functions[i].Name()
			fileName := strings.Replace(taArgsOutputTemplate, "%", functionName, 1)
			var jsonText []byte
			jsonText, err = json.Marshal(argInfo)
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
		functionInfos := make([]FunctionKleeInfo, len(functions))
		for i, function := range functions {
			functionInfos[i] = FunctionKleeInfo{
				Name: function.Name(),
				Args: argInfos[i],
			}
		}
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
