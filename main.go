package main

import (
    "os"
    "io/ioutil"
    "regexp"
)

func CreateTestingCode (testedFileName string) string {
    testedFile, _ := NewFile(testedFileName)
    testSet := testedFile.GenerateTests()
    return PrintTestingFile(testedFile, testSet)
}

func CreateTestingFile (testedFileName string) {
    testingFileName := testedFileName[:len(testedFileName)-3] + "_test.go"
    testingFile, _ := os.OpenFile(
        testingFileName,
        os.O_WRONLY|os.O_CREATE,
        0644)
    defer testingFile.Close()
    testingCode := CreateTestingCode(testedFileName)
    testingFile.WriteString(testingCode)
}

func CreateTestingDir (testedDirName string) {
    files, _ := ioutil.ReadDir(testedDirName)
    for _, file := range files {
        fileName := file.Name()
        isGo, _ := regexp.MatchString(`.*\.go$`, fileName)
        isGoTest, _ := regexp.MatchString(`.*_test\.go$`, fileName)
        if isGo && ! isGoTest && ! file.IsDir() {
            CreateTestingFile(testedDirName + fileName)
        }
    }
}

func main () {
    CreateTestingDir("test-files/")
}
