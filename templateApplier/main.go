package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	args := append(os.Args, "-c")
	var (
		changedFiles []string
		data []byte
	)
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
			err = os.WriteFile(fileName, data, 0644)
			throw(err)
		case "-r":
			i++
			fileName := args[i]
			if data, err = os.ReadFile(fileName); err == nil {
				err = os.Remove(fileName)
				throw(err)
				changedFiles = append(changedFiles, strings.Fields(string(data))...)
			}
		case "-w":
			i++
			if changedFiles != nil {
				fileName := args[i]
				if _, err = os.Stat(fileName); err == nil {
					data, err = os.ReadFile(fileName)
					throw(err)
					data = append(data, ' ')
				} else {
					data = nil
				}
				data = append(data, strings.Join(changedFiles, " ")...)
				err = os.WriteFile(fileName, data, 0644)
				throw(err)
				changedFiles = nil
			}
		case "-c":
			newTime := time.Now().Local().Add(time.Second)
			for _, fileName := range changedFiles {
				err = os.Chtimes(fileName, newTime, newTime)
				throw(err)
			}
			changedFiles = nil
		default:
			fileName := args[i]
			var wasChanged bool
			wasChanged, err = session.ExecuteFile(fileName)
			throw(err)
			if wasChanged {
				changedFiles = append(changedFiles, fileName)
			}
		}
	}
}
