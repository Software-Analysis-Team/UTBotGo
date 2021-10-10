package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GenerateTestingGoFile(testedFile FSObject) bool {
	testedFileName := testedFile.Name
	testingFileName :=
		strings.TrimSuffix(testedFileName, filepath.Ext(testedFileName)) + "_test.go"
	testingFile, err := os.OpenFile(
		testingFileName,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644, // read for all users, write only for owner
	)
	if err != nil { return false }
	testedGoFile, err := NewGoFile(testedFileName)
	if err != nil { return false }
	testSet := testedGoFile.GenerateTests()
	testingCode := PrintTestingFile(testedGoFile, testSet)
	_, err = testingFile.WriteString(testingCode)
	if err != nil { return false }
	err = testingFile.Close()
	if err != nil { return false }
	return true
}

func GenerateTestingFile(testedFile FSObject) (genFilesAmt int) {
	isGo, _ := regexp.MatchString(`.*\.go$`, testedFile.Name)
	isGoTest, _ := regexp.MatchString(`.*_test\.go$`, testedFile.Name)
	if isGo && !isGoTest {
		if GenerateTestingGoFile(testedFile) { genFilesAmt = 1 }
	}
	return
}

func GenerateTestingDir(testedDir FSObject, recursively bool) (genFilesAmt int) {
	if recursively {
		for _, testedObject := range testedDir.ContainedFiles() {
			genFilesAmt += GenerateTestingObject(testedObject, true)
		}
	} else {
		for _, testedObject := range testedDir.ContainedFiles() {
			if testedObject.IsFile() {
				genFilesAmt += GenerateTestingFile(testedObject)
			}
		}
	}
	return
}

func GenerateTestingObject(testedObject FSObject, recursively bool) (genFilesAmt int) {
	switch {
	case testedObject.IsDir():
		genFilesAmt += GenerateTestingDir(testedObject, recursively)
	case testedObject.IsFile():
		genFilesAmt += GenerateTestingFile(testedObject)
	}
	return
}

func main() {
	inputFile := NewFSObject(os.Args[1])
	genFilesAmt := GenerateTestingObject(inputFile, true)
	fmt.Println(genFilesAmt)
}
