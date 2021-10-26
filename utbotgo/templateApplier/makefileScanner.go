package main

import "strings"

type MakefileData struct {
	FillValue bool
	FillSubValue bool
	Key strings.Builder
	Value strings.Builder
	SubValue strings.Builder
	Values map[string]string
}

func NewMakefileScanner() RuneScanner {
	return RuneScanner{
		MakefileData{Values: make(map[string]string)},
		func(r rune, res interface{}) interface{} {
			switch data := res.(type) {
			case MakefileData:
				switch r {
				case '\n':
					if data.FillValue {
						data.Values[data.Key.String()] = data.Value.String()
					}
					data.FillValue = false
					data.Key.Reset()
					data.Value.Reset()
				case ' ':
					if data.FillValue && data.Value.Len() != 0 {
						data.Value.WriteRune(r)
					}
				case '$':
					if data.FillValue {
						data.FillSubValue = true
					}
				case '=':
					data.FillValue = true
				default:
					if data.FillValue {
						if data.FillSubValue {
							switch r {
							case '(':
							case ')':
								subValue := data.Values[data.SubValue.String()]
								data.Value.WriteString(subValue)
								data.FillSubValue = false
								data.SubValue.Reset()
							default:
								data.SubValue.WriteRune(r)
							}
						} else {
							data.Value.WriteRune(r)
						}
					} else {
						data.Key.WriteRune(r)
					}
				}
				return data
			default:
				return nil
			}
		},
	}
}
