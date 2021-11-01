package main

import (
	"os"
	"path/filepath"
)

type PutValuesFunc func([]byte, Values) []byte

type PutValuesFuncForSpecFile struct {
	GlobPattern   string
	PutValuesFunc PutValuesFunc
}

var patterns = []PutValuesFuncForSpecFile{
	{
		`.gitignore`,
		PutValuesFuncForGitignore,
	},
	{
		`*.md`,
		PutValuesFuncForMarkdown,
	},
}

func PutValuesFuncForFile(patterns []PutValuesFuncForSpecFile, fileName string) (PutValuesFunc, error) {
	var (
		m bool
		e error
	)
	for _, pattern := range patterns {
		m, e = filepath.Match(pattern.GlobPattern, fileName)
		if e != nil || m {
			return pattern.PutValuesFunc, e
		}
	}
	return nil, nil
}

func PutValuesIntoFile(filePath string, values Values) error {
	fileName := filepath.Base(filePath)
	putValuesFunc, err := PutValuesFuncForFile(patterns, fileName)
	if err != nil || putValuesFunc == nil {
		return err
	}
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fileContent = putValuesFunc(fileContent, values)
	if fileContent == nil {
		return nil
	}
	err = os.WriteFile(filePath, fileContent, 0)
	return err
}
