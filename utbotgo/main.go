package main

import (
	"fmt"
	"os"
	"strconv"
)

type GeneratingResult struct {
	GeneratingFilesAmount int
}

func (result GeneratingResult) String() (resultStr string) {
	resultStr = strconv.Itoa(result.GeneratingFilesAmount)
	return
}

func main() {
	parsedArgs, err := ParseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result, err := GenerateTestingObject(parsedArgs, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result)
}
