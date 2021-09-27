package main

import (
	"io/ioutil"
	"os"
)

func NewFSObject(fileName string) FSObject {
	fileInfo, _ := os.Stat(fileName)
	name := fileInfo.Name()
	mode := fileInfo.Mode()
	return FSObject{name, mode}
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
