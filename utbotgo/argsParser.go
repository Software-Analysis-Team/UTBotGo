package main

import (
	"fmt"
	"os"
)

type FuzzingMode int

const (
	Gofuzz FuzzingMode = iota
	Go_fuzz
)

type ParsedArgs struct {
	Object FSObject
	Mode   FuzzingMode
}

func ParseArgs() (ParsedArgs, error) {
	var (
		objStr string
		obj    FSObject
		fm     FuzzingMode
	)
	if len(os.Args) < 2 {
		return ParsedArgs{}, fmt.Errorf("not enough args")
	}
	switch os.Args[1] {
	case "--go-fuzz":
		fm = Go_fuzz
		objStr = os.Args[2]
	default:
		fm = Gofuzz
		objStr = os.Args[1]
	}
	obj = NewFSObject(objStr)
	return ParsedArgs{obj, fm}, nil
}
