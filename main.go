package main

import (
	"io/ioutil"
	"os"
	"regexp"
)

func CreateTestingCode(testedFileName string) string {
	testedFile, _ := NewFile(testedFileName)
	testSet := testedFile.GenerateTests()
	return PrintTestingFile(testedFile, testSet)
}

func CreateTestingFile(testedFileName string) {
	isGo, _ := regexp.MatchString(`.*\.go$`, testedFileName)
	if !isGo {
		return
	}
	testingFileName := testedFileName[:len(testedFileName)-3] + "_test.go"
	testingFile, _ := os.OpenFile(
		testingFileName,
		os.O_WRONLY|os.O_CREATE,
		0644)
	testingCode := CreateTestingCode(testedFileName)
	_, _ = testingFile.WriteString(testingCode)
	_ = testingFile.Close()
}

func CreateTestingDir(testedDirName string) {
	files, _ := ioutil.ReadDir(testedDirName)
	for _, file := range files {
		fileName := file.Name()
		isGoTest, _ := regexp.MatchString(`.*_test\.go$`, fileName)
		if !isGoTest && !file.IsDir() {
			CreateTestingFile(testedDirName + fileName)
		}
	}
}

func CreateTestingObject(testedObject string) {
	info, _ := os.Stat(testedObject)
	switch mode := info.Mode(); {
	case mode.IsDir():
		CreateTestingDir(testedObject)
	case mode.IsRegular():
		CreateTestingFile(testedObject)
	}
}

func main() {
	CreateTestingObject(os.Args[1])
}
