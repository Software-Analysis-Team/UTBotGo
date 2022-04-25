package main

import (
	"bytes"
	"os"
	"regexp"
)

func main() {
	kleeVars := make(map[string][]byte)
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	text, _ := os.ReadFile(inputFile)
	kleeVarDefRe := regexp.MustCompile(`(?m)^@klee_modify\.\.((?:\w|\.)+) = internal constant i8 (\d+)$`)
	for _, submatch := range kleeVarDefRe.FindAllSubmatch(text, -1) {
		kleeVars[string(submatch[1])] = submatch[2]
	}
	kleeVarUsingRe := regexp.MustCompile(`i8\* @klee_modify\.\.((?:\w|\.)+)`)
	text = kleeVarUsingRe.ReplaceAllFunc(text, func(submatch []byte) []byte {
		kleeVar := string(kleeVarUsingRe.FindSubmatch(submatch)[1])
		var buf bytes.Buffer
		buf.WriteString("i8* getelementptr inbounds ([")
		buf.Write(kleeVars[kleeVar])
		buf.WriteString(" x i8], [")
		buf.Write(kleeVars[kleeVar])
		buf.WriteString(" x i8]* @")
		buf.WriteString(kleeVar)
		buf.WriteString(", i64 0, i64 0)")
		return buf.Bytes()
	})
	_ = os.WriteFile(outputFile, text, 0666)
}
