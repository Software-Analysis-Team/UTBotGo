package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func throw(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	wd, err := os.Getwd()
	throw(err)
	wd, err = filepath.EvalSymlinks(wd)
	throw(err)
	session, err := NewSession(wd)
	throw(err)
	args := os.Args
	var data []byte
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-i":
			i++
			fileName := args[i]
			data, err = os.ReadFile(fileName)
			throw(err)
			err = session.FromJSON(data)
			throw(err)
		case "-o":
			i++
			fileName := args[i]
			data, err = session.ToJSON()
			throw(err)
			err = os.WriteFile(fileName, data, 0666)
			throw(err)
		case "-c":
			i++
			inputFileName := args[i]
			i++
			outputFileName := args[i]
			err = session.ExecuteFile(inputFileName, outputFileName)
			throw(err)
		default:
			fileName := args[i]
			err = session.ExecuteFile(fileName, fileName)
			throw(err)
		}
	}
}
