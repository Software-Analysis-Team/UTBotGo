package main

import (
	"os"
	"regexp"
)

func CreateTestingGoFile(testedFile FSObject) {
	testedFileName := testedFile.Name
	testingFileName := testedFileName[:len(testedFileName)-3] + "_test.go"
	testingFile, _ := os.OpenFile(
		testingFileName,
		os.O_WRONLY|os.O_CREATE,
		0644,
	)
	testedGoFile, _ := NewGoFile(testedFileName)
	testSet := testedGoFile.GenerateTests()
	testingCode := PrintTestingFile(testedGoFile, testSet)
	_, _ = testingFile.WriteString(testingCode)
	_ = testingFile.Close()
}

func CreateTestingFile(testedFile FSObject) {
	isGo, _ := regexp.MatchString(`.*\.go$`, testedFile.Name)
	isGoTest, _ := regexp.MatchString(`.*_test\.go$`, testedFile.Name)
	if isGo && !isGoTest {
		CreateTestingGoFile(testedFile)
	}
}

func CreateTestingDir(testedDir FSObject, recursively bool) {
	for _, testedObject := range testedDir.ContainedFiles() {
		if recursively {
			CreateTestingObject(testedObject, true)
		} else {
			if testedObject.IsFile() {
				CreateTestingFile(testedObject)
			}
		}
	}
}

func CreateTestingObject(testedObject FSObject, recursively bool) {
	switch {
	case testedObject.IsDir():
		CreateTestingDir(testedObject, recursively)
	case testedObject.IsFile():
		CreateTestingFile(testedObject)
	}
}

func main() {
	inputFile := NewFSObject(os.Args[1])
	CreateTestingObject(inputFile, false)
}
