package main

import (
	"io/fs"
	"io/ioutil"
	"os"
)

type FSObject struct {
	Name string
	Mode fs.FileMode
}

func NewFSObject(fileName string) FSObject {
	fileInfo, _ := os.Stat(fileName)
	mode := fileInfo.Mode()
	return FSObject{fileName, mode}
}

func (fso FSObject) IsDir() bool {
	return fso.Mode.IsDir()
}

func (fso FSObject) IsFile() bool {
	return fso.Mode.IsRegular()
}

func (fso FSObject) ContainedFiles() []FSObject {
	name := fso.Name
	if name == "" || name[len(name)-1] != '/' {
		name += "/"
	}
	subFilesInfo, err := ioutil.ReadDir(fso.Name)
	if err != nil {
		return nil
	}
	res := make([]FSObject, len(subFilesInfo))
	for i, subFileInfo := range subFilesInfo {
		subName := subFileInfo.Name()
		subMode := subFileInfo.Mode()
		res[i] = FSObject{name + subName, subMode}
	}
	return res
}
