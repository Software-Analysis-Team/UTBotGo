package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GenerateTestingCode(parsedArgs ParsedArgs) (code string, err error) {
	fileName := parsedArgs.Object.Name
	file, subErr := ParseGoFile(fileName)
	if subErr != nil {
		return "", subErr
	}
	switch parsedArgs.Mode {
	case Gofuzz:
		code = GofuzzTestingCode(file)
	case Go_fuzz:
		code = Go_fuzzTestingCode(file)
	default:
		code = TestingFuzzTestingCode(file)
	}
	err = nil
	return
}

func GenerateTestingGoFile(parsedArgs ParsedArgs) bool {
	testedFileName := parsedArgs.Object.Name
	testingFileName :=
		strings.TrimSuffix(testedFileName, filepath.Ext(testedFileName)) + "_test.go"
	testingFile, err := os.OpenFile(
		testingFileName,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644, // read for all users, write only for owner
	)
	if err != nil {
		return false
	}
	testingCode, err := GenerateTestingCode(parsedArgs)
	if err != nil {
		return false
	}
	_, err = testingFile.WriteString(testingCode)
	if err != nil {
		return false
	}
	err = testingFile.Close()
	if err != nil {
		return false
	}
	return true
}

func GenerateTestingFile(parsedArgs ParsedArgs) (result GeneratingResult) {
	isGo, err := regexp.MatchString(`.*\.go$`, parsedArgs.Object.Name)
	if err != nil {
		return
	}
	isGoTest, err := regexp.MatchString(`.*_test\.go$`, parsedArgs.Object.Name)
	if err != nil {
		return
	}
	if isGo && !isGoTest {
		if GenerateTestingGoFile(parsedArgs) {
			result.GeneratingFilesAmount = 1
		}
	}
	return
}

func GenerateTestingDir(parsedArgs ParsedArgs, recursively bool) (result GeneratingResult) {
	if recursively {
		for _, testedObject := range parsedArgs.Object.ContainedFiles() {
			subArgs := parsedArgs
			subArgs.Object = testedObject
			subResult := GenerateTestingObject(subArgs, true)
			result.GeneratingFilesAmount += subResult.GeneratingFilesAmount
		}
	} else {
		for _, testedObject := range parsedArgs.Object.ContainedFiles() {
			if testedObject.IsFile() {
				subArgs := parsedArgs
				subArgs.Object = testedObject
				subResult := GenerateTestingFile(subArgs)
				result.GeneratingFilesAmount += subResult.GeneratingFilesAmount
			}
		}
	}
	return
}

func GenerateTestingObject(parsedArgs ParsedArgs, recursively bool) (result GeneratingResult) {
	switch {
	case parsedArgs.Object.IsDir():
		result = GenerateTestingDir(parsedArgs, recursively)
	case parsedArgs.Object.IsFile():
		result = GenerateTestingFile(parsedArgs)
	}
	return
}
