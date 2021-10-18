package main

import (
	"fmt"
	"os"
	"strconv"
)

type FuzzingMode int

const (
	Gofuzz FuzzingMode = iota
	Go_fuzz
	TestingFuzz
)

type GeneratingResult struct {
	GeneratingFilesAmount int
}

type ParsedArgs struct {
	Object FSObject
	Mode   FuzzingMode
}

func (result GeneratingResult) String() (resultStr string) {
	resultStr = strconv.Itoa(result.GeneratingFilesAmount)
	return
}

func ParseArgs() ParsedArgs {
	var (
		objStr string
		obj    FSObject
		fm     FuzzingMode
	)
	switch os.Args[1] {
	case "--gofuzz":
		fm = Gofuzz
		objStr = os.Args[2]
	case "--go-fuzz":
		fm = Go_fuzz
		objStr = os.Args[2]
	default:
		fm = TestingFuzz
		objStr = os.Args[1]
	}
	obj = NewFSObject(objStr)
	return ParsedArgs{obj, fm}
}

func main() {
	parsedArgs := ParseArgs()
	result := GenerateTestingObject(parsedArgs, true)
	fmt.Println(result)
}
