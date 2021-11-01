package main

import (
	"bytes"
	"regexp"
)

func PutValuesFuncForGitignore(text []byte, values Values) []byte {
	useTemplatePattern := regexp.MustCompile(`#>\s+(.+)\n(.*)\n`)
	putValuePattern := regexp.MustCompile(`\{(\w+)\}`)
	return TemplatePutValuesFunc(text, values, useTemplatePattern, putValuePattern)
}

func PutValuesFuncForMarkdown(text []byte, values Values) []byte {
	useTemplatePattern := regexp.MustCompile(`(?s:<!--\s+\[\s*(.+?)\s+-->\s*(.*?)\s*<!--\s+\]\s+-->)`)
	putValuePattern := regexp.MustCompile(`\{(\w+)\}`)
	return TemplatePutValuesFunc(text, values, useTemplatePattern, putValuePattern)
}

func TemplatePutValuesFunc(text []byte, values Values, useTemplatePattern *regexp.Regexp, putValuePattern *regexp.Regexp) []byte {
	indices := useTemplatePattern.FindAllSubmatchIndex(text, -1)
	if indices == nil {
		return nil
	}
	var bb bytes.Buffer
	i := 0
	for _, index := range indices {
		// `index` has only 6 elements:
		//   0 and 1 - start and end of match
		//   2 and 3 - start and end of template
		//   4 and 5 - start and end of place for applying template
		bb.Write(text[i : index[4]])
		i = index[5]
		template := text[index[2] : index[3]]
		templatedText := putValuePattern.ReplaceAllFunc(template, func(match []byte) []byte {
			key := string(putValuePattern.FindSubmatch(match)[1])
			return values[key]
		})
		bb.Write(templatedText)
	}
	bb.Write(text[i:])
	return bb.Bytes()
}
