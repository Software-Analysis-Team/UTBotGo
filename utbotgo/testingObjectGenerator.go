package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GenerateTestingCode(parsedArgs ParsedArgs) (code string, err error) {
	fileName := parsedArgs.Object.Name
	file, err := ParseGoFile(fileName)
	if err != nil {
		return
	}
	switch parsedArgs.Mode {
	case Gofuzz:
		code = GofuzzTestingCode(file)
	case Go_fuzz:
		code = Go_fuzzTestingCode(file)
	default:
		code = TestingFuzzTestingCode(file)
	}
	return
}

func GenerateTestingGoFile(parsedArgs ParsedArgs) (err error) {
	testedFileName := parsedArgs.Object.Name
	testingFileName :=
		strings.TrimSuffix(testedFileName, filepath.Ext(testedFileName)) + "_test.go"
	testingCode, err := GenerateTestingCode(parsedArgs)
	if err != nil {
		return
	}
	err = os.WriteFile(
		testingFileName,
		[]byte(testingCode),
		0644, // read for all users, write only for owner
	)
	return
}

func GenerateTestingFile(parsedArgs ParsedArgs) (result GeneratingResult, err error) {
	isGo, err := regexp.MatchString(`.*\.go$`, parsedArgs.Object.Name)
	if err != nil {
		return
	}
	isGoTest, err := regexp.MatchString(`.*_test\.go$`, parsedArgs.Object.Name)
	if err != nil {
		return
	}
	if isGo && !isGoTest {
		err = GenerateTestingGoFile(parsedArgs)
		if err != nil {
			return
		}
		result.GeneratingFilesAmount = 1
	}
	return
}

func GenerateTestingDir(parsedArgs ParsedArgs, recursively bool) (result GeneratingResult, err error) {
	var subResult GeneratingResult
	if recursively {
		for _, testedObject := range parsedArgs.Object.ContainedFiles() {
			subArgs := parsedArgs
			subArgs.Object = testedObject
			subResult, err = GenerateTestingObject(subArgs, true)
			if err != nil {
				return
			}
			result.GeneratingFilesAmount += subResult.GeneratingFilesAmount
		}
	} else {
		for _, testedObject := range parsedArgs.Object.ContainedFiles() {
			if testedObject.IsFile() {
				subArgs := parsedArgs
				subArgs.Object = testedObject
				subResult, err = GenerateTestingFile(subArgs)
				if err != nil {
					return
				}
				result.GeneratingFilesAmount += subResult.GeneratingFilesAmount
			}
		}
	}
	return
}

func GenerateTestingObject(parsedArgs ParsedArgs, recursively bool) (result GeneratingResult, err error) {
	switch {
	case parsedArgs.Object.IsDir():
		result, err = GenerateTestingDir(parsedArgs, recursively)
	case parsedArgs.Object.IsFile():
		result, err = GenerateTestingFile(parsedArgs)
	default:
		err = fmt.Errorf("unknown object: %s", parsedArgs.Object.Name)
	}
	return
}
