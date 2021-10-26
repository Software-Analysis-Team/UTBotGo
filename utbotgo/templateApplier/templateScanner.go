package main

import "strings"

type TemplateData struct {
	FillValue bool
	Value strings.Builder
	Values map[string]string
	Output strings.Builder
}

func NewTemplateScanner(values map[string]string) RuneScanner {
	return RuneScanner{
		TemplateData{Values: values},
		func(r rune, res interface{}) interface{} {
			switch data := res.(type) {
			case TemplateData:
				switch r {
				case '{':
					data.FillValue = true
				case '}':
					newValue := data.Values[data.Value.String()]
					data.Output.WriteString(newValue)
					data.FillValue = false
					data.Value.Reset()
				default:
					if data.FillValue {
						data.Value.WriteRune(r)
					} else {
						data.Output.WriteRune(r)
					}
				}
				return data
			default:
				return nil
			}
		},
	}
}
