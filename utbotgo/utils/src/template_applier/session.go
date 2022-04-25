package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Session struct {
	Variables Struct
	Syntaxes  FileSyntaxes
}

func NewSession(wd string) (Session, error) {
	fileName, err := SearchTASDfile(wd)
	if err != nil {
		return Session{}, err
	}
	syntaxes, err := ReadTASDfile(fileName)
	if err != nil {
		return Session{}, err
	}
	return Session{
		Variables: make(Struct),
		Syntaxes:  syntaxes,
	}, nil
}

func (session *Session) GetSyntax(fileName string) Syntax {
	fileBase := filepath.Base(fileName)
	if syntax, ok := session.Syntaxes.ForBase[fileBase]; ok {
		return syntax
	}
	fileExt := filepath.Ext(fileName)
	if syntax, ok := session.Syntaxes.ForExt[fileExt]; ok {
		return syntax
	}
	return nil
}

var fieldNameNames = [5]Name{
	"oneFieldName",
	"allFieldsName",
	"strictNewFieldName",
	"maybeNewFieldName",
	"fieldName",
}

func (session *Session) ApplyTemplate(rTemplate Replacer) []String {
	var varNames []VarName
	for _, rGetterVar := range rTemplate.GetAll("getterVar") {
		if rGetterVar.ThereIs("template") && (! rTemplate.ThereIs("getterVar")) {
			rSubTemplate := rGetterVar.Get("template")
			appliedSubTemplate := session.ApplyTemplate(rSubTemplate)
			var separator String
			if rGetterVar.ThereIs("separator") {
				separator = rGetterVar.Get("separator").GetString()
			}
			rGetterVar.ReplaceString(separator.JoinStrings(appliedSubTemplate))
			rTemplate.UpdateReplaces(&rGetterVar)
		} else {
			rVarName := rGetterVar.Get("varName")
			var varName VarName
			for _, rField := range rVarName.GetAll("field") {
				var fieldName FieldName
				for _, fieldNameName := range fieldNameNames {
					if rField.ThereIs(fieldNameName) {
						fieldName = FieldName(rField.Get(fieldNameName).GetString())
						break
					}
				}
				varName = append(varName, fieldName)
			}
			varNames = append(varNames, varName)
		}
	}
	strings := session.Variables.GetStrings(varNames)
	rsAppliedTemplate := rTemplate.SubReplaceManyStrings("getterVar", strings)
	appliedTemplate := make([]String, len(rsAppliedTemplate))
	for i := range appliedTemplate {
		appliedTemplate[i] = rsAppliedTemplate[i].GetString()
	}
	return appliedTemplate
}

func (session *Session) ExecuteScript(src []byte, syntax Syntax) String {
	r := NewReplacer(src, syntax)
	for _, rCodeBlock := range r.GetAll("codeBlock") {
		for _, rStatement := range rCodeBlock.GetAll("statement") {
			rTemplate := rStatement.Get("template")
			appliedTemplate := session.ApplyTemplate(rTemplate)
			if rStatement.ThereIs("result") {
				rResult := rStatement.Get("result")
				i := len(appliedTemplate) - 1
				if i == -1 {
					rResult.ReplaceString(nil)
				} else {
					rResult.ReplaceString(appliedTemplate[i])
				}
				rStatement.UpdateReplaces(&rResult)
				rCodeBlock.UpdateReplaces(&rStatement)
				r.UpdateReplaces(&rCodeBlock)
			}
			if rStatement.ThereIs("setterVar") || rStatement.ThereIs("adderVar") {
				var rExtVarName Replacer
				if rStatement.ThereIs("extVarName") {
					rExtVarName = rStatement.Get("extVarName")
				} else {
					rExtVarName = rStatement.Get("varName")
				}
				var extVarName ExtVarName
				for _, rField := range rExtVarName.GetAll("field") {
					var extFieldName ExtFieldName
					switch {
					case rField.ThereIs("oneFieldName"):
						extFieldName = ExtFieldName{
							Name: FieldName(rField.Get("oneFieldName").GetString()),
							Type: OneField,
						}
					case rField.ThereIs("allFieldsName"):
						extFieldName = ExtFieldName{
							Name: FieldName(rField.Get("allFieldsName").GetString()),
							Type: AllFields,
						}
					case rField.ThereIs("strictNewFieldName"):
						extFieldName = ExtFieldName{
							Name: FieldName(rField.Get("strictNewFieldName").GetString()),
							Type: StrictNewField,
						}
					case rField.ThereIs("maybeNewFieldName"):
						extFieldName = ExtFieldName{
							Name: FieldName(rField.Get("maybeNewFieldName").GetString()),
							Type: MaybeNewField,
						}
					default:
						extFieldName = ExtFieldName{
							Name: FieldName(rField.Get("fieldName").GetString()),
							Type: OneField,
						}
					}
					extVarName = append(extVarName, extFieldName)
				}
				if rStatement.ThereIs("setterVar") {
					session.Variables.SetStrings(extVarName, appliedTemplate)
				} else {
					session.Variables.AddStrings(extVarName, appliedTemplate)
				}
			}
		}
	}
	return r.GetString()
}

func (session *Session) ExecuteFile(inputFileName, outputFileName string) (err error) {
	source, err := os.ReadFile(inputFileName)
	if err != nil {
		return
	}
	syntax := session.GetSyntax(inputFileName)
	if syntax == nil {
		return
	}
	result := session.ExecuteScript(source, syntax)
	wasChanged := false
	if inputFileName != outputFileName {
		wasChanged = true
	} else if len(source) != len(result) {
		wasChanged = true
	} else {
		for i := range source {
			if source[i] != result[i] {
				wasChanged = true
				break
			}
		}
	}
	if wasChanged {
		err = os.WriteFile(outputFileName, result, 0666)
	}
	return
}

func (session *Session) ToJSON() ([]byte, error) {
	return json.Marshal(ToGoStruct(session.Variables))
}

func (session *Session) FromJSON(data []byte) error {
	var result interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	session.Variables = FromGoStruct(result).(Struct)
	return nil
}
