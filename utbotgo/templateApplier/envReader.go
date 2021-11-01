package main

import (
	"fmt"
	"regexp"
)

type Values map[string][]byte

func ReadEnv(text []byte) Values {
	values := make(map[string][]byte)
	keyPattern := regexp.MustCompile(`(\w+)`)
	valuePattern := regexp.MustCompile(`(.*)`)
	putValuePattern := regexp.MustCompile(fmt.Sprintf(`\$\(%s\)`, keyPattern))
	envPattern := regexp.MustCompile(fmt.Sprintf(`%s\s*=\s*%s\n`, keyPattern, valuePattern))
	for _, match := range envPattern.FindAllSubmatch(text, -1) {
		key := string(match[1])
		value := putValuePattern.ReplaceAllFunc(match[2], func(submatch []byte) []byte {
			subkey := string(keyPattern.FindSubmatch(submatch)[1])
			return values[subkey]
		})
		values[key] = value
	}
	return values
}
