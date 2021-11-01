package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	envText, _ := io.ReadAll(os.Stdin)
	values := ReadEnv(envText)
	var err error
	for _, filePath := range os.Args[1:] {
		err = PutValuesIntoFile(filePath, values)
		if err != nil {
			fmt.Println(err)
		}
	}
}
