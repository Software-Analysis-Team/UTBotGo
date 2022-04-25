package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SubstringIndex struct {
	Begin int
	End   int
}

func (index SubstringIndex) IsSubIndex(supindex SubstringIndex) bool {
	return supindex.Begin <= index.Begin && supindex.End >= index.End
}

type Selector struct {
	Regexps map[Name]*regexp.Regexp
	Names []Name
	Subnames map[Name]Name
	Alternatives map[Name][]*regexp.Regexp
	GeneralIsAlt map[Name]struct{}
	Aliases map[Name][]Name
}

func NewSelector(syntax Syntax) Selector {
	var selector Selector
	selector.Regexps = make(map[Name]*regexp.Regexp)
	selector.Subnames = make(map[Name]Name)
	selector.Alternatives = make(map[Name][]*regexp.Regexp)
	selector.GeneralIsAlt = make(map[Name]struct{})
	selector.Aliases = make(map[Name][]Name)
	for name := range syntax {
		selector.createRegexpByName(name, syntax)
	}
	return selector
}

var openP, closeP, recS, altS, multiS, multiAltS, aliaS = '<', '>', '*', '|', ':', '+', '='

var subnameRegexp = regexp.MustCompile(`(\w+)(?:\[(\d+)\])?`)

func (selector *Selector) createRegexpByName(name Name, syntax Syntax) {
	if selector.thereIsName(name) {
		return
	}
	src := syntax[name]
	selector.createRegexpByNameStart(name, &src)
	var openIdxs []int
	for i := range src {
		if src[i] == uint8(openP) && (i == 0 || src[i-1] != '\\') {
			openIdxs = append(openIdxs, i)
		}
	}
	if len(openIdxs) == 0 {
		selector.Regexps[name] = regexp.MustCompile(src)
		return
	}
	var closeIdxs []int
	for i := range src {
		if src[i] == uint8(closeP) && (i == 0 || src[i-1] != '\\') {
			closeIdxs = append(closeIdxs, i)
		}
	}
	if len(closeIdxs) == 0 {
		selector.Regexps[name] = regexp.MustCompile(src)
		return
	}
	oi, ci := 0, 0
	var builder strings.Builder
	builder.WriteString(src[: openIdxs[0]])
	for oi < len(openIdxs) && ci < len(closeIdxs) {
		if openIdxs[oi] < closeIdxs[ci] && (oi == len(openIdxs)-1 || openIdxs[oi+1] > closeIdxs[ci]) {
			selector.createRegexpByNameForSubstring(name, src[openIdxs[oi]+1 : closeIdxs[ci]], syntax, &builder)
			if oi == len(openIdxs)-1 {
				builder.WriteString(src[closeIdxs[ci]+1 :])
			} else {
				builder.WriteString(src[closeIdxs[ci]+1 : openIdxs[oi+1]])
			}
			oi++
			ci++
		} else if openIdxs[oi] < closeIdxs[ci] && openIdxs[oi+1] < closeIdxs[ci] {
			builder.WriteString(src[openIdxs[oi] : openIdxs[oi+1]])
			oi++
		} else {
			ci++
		}
	}
	if oi < len(openIdxs) {
		builder.WriteString(src[openIdxs[oi] :])
	}
	selector.createRegexpByNameFinal(name, syntax, &builder)
}

func (selector *Selector) thereIsName(name Name) bool {
	_, thereIsRegexp := selector.Regexps[name]
	_, thereAreAlternatives := selector.Alternatives[name]
	return thereIsRegexp || thereAreAlternatives
}

func (selector *Selector) createRegexpByNameStart(name Name, src *string) {
	for true {
		switch {
		case len(*src) > 3 && (*src)[0] == uint8(openP) && (*src)[1] == uint8(recS):
			var i int
			for i = 2; i < len(*src) && (*src)[i] != uint8(closeP); i++ {}
			if i == len(*src) {
				return
			}
			selector.Subnames[name] = Name((*src)[2 : i])
			*src = (*src)[i + 1 :]
		case len(*src) > 3 && (*src)[0] == uint8(openP) && (*src)[1] == uint8(aliaS):
			var i int
			for i = 2; i < len(*src) && (*src)[i] != uint8(closeP); i++ {}
			if i == len(*src) {
				return
			}
			if _, ok := selector.Aliases[name]; !ok {
				selector.Aliases[name] = []Name{}
			}
			selector.Aliases[name] = append(selector.Aliases[name], Name((*src)[2 : i]))
			*src = (*src)[i + 1 :]
		default:
			return
		}
	}
}

func (selector *Selector) createRegexpByNameForSubstring(name Name, substring string, syntax Syntax, builder *strings.Builder) {
	switch {
	case substring == string(altS):
		if _, ok := selector.Alternatives[name]; !ok {
			selector.Alternatives[name] = []*regexp.Regexp{}
		}
		altRegexp := regexp.MustCompile(builder.String())
		selector.Alternatives[name] = append(selector.Alternatives[name], altRegexp)
		*builder = strings.Builder{}
	case substring == string(multiS):
		selector.Regexps[name] = regexp.MustCompile(builder.String())
		selector.Alternatives[name] = []*regexp.Regexp{}
		*builder = strings.Builder{}
	case substring == string(multiAltS):
		selector.Regexps[name] = regexp.MustCompile(builder.String())
		selector.Alternatives[name] = []*regexp.Regexp{}
		selector.GeneralIsAlt[name] = struct{}{}
		*builder = strings.Builder{}
	case subnameRegexp.MatchString(substring):
		submatch := subnameRegexp.FindStringSubmatch(substring)
		subname := Name(submatch[1])
		selector.createRegexpByName(subname, syntax)
		var re *regexp.Regexp
		if submatch[2] == "" {
			re = selector.Regexps[subname]
			if _, ok := selector.Alternatives[subname]; ok {
				subname = nameEncoder.Encode(subname)
			}
		} else {
			index, _ := strconv.Atoi(submatch[2])
			index--
			re = selector.Alternatives[subname][index]
		}
		builder.WriteString(fmt.Sprintf("(?P<_%d>%s)", len(selector.Names), re.String()))
		selector.Names = append(selector.Names, subname)
	default:
		builder.WriteString(fmt.Sprintf("%s%s%s", string(openP), substring, string(closeP)))
	}
}

func (selector *Selector) createRegexpByNameFinal(name Name, syntax Syntax, builder *strings.Builder) {
	re := regexp.MustCompile(builder.String())
	if _, ok := selector.Alternatives[name]; ok {
		selector.Alternatives[name] = append(selector.Alternatives[name], re)
	} else {
		selector.Regexps[name] = re
	}
}

var subexpNamePat = regexp.MustCompile(`_(\d+)`)

func (selector *Selector) selectIndicesByNames(name Name, subexpNames []string, submatchIndices []int) map[Name]SubstringIndex {
	result := make(map[Name]SubstringIndex)
	for i, subexpName := range subexpNames {
		subexpIndex := SubstringIndex{
			Begin: submatchIndices[i*2],
			End: submatchIndices[i*2+1],
		}
		if subexpIndex.Begin != -1 {
			if i == 0 {
				names := []Name{name}
				if aliases, ok := selector.Aliases[name]; ok {
					names = append(names, aliases...)
				}
				for _, name = range names {
					if _, ok := result[name]; !ok || subexpIndex.Begin != subexpIndex.End {
						result[name] = subexpIndex
					}
				}
			} else if subexpNamePat.MatchString(subexpName) {
				idx, _ := strconv.Atoi(subexpNamePat.FindStringSubmatch(subexpName)[1])
				subname := selector.Names[idx]
				subnames := []Name{subname}
				if aliases, ok := selector.Aliases[subname]; ok {
					subnames = append(subnames, aliases...)
				}
				for _, subname = range subnames {
					if _, ok := result[subname]; !ok || subexpIndex.Begin != subexpIndex.End {
						result[subname] = subexpIndex
					}
				}
			}
		}
	}
	return result
}

func (selector *Selector) selectAllMatches(name Name, re *regexp.Regexp, src String) []map[Name]SubstringIndex {
	subexpNames := re.SubexpNames()
	submatchIndices := re.FindAllSubmatchIndex(src, -1)
	result := make([]map[Name]SubstringIndex, len(submatchIndices))
	for j, submatchIndex := range submatchIndices {
		result[j] = selector.selectIndicesByNames(name, subexpNames, submatchIndex)
	}
	return result
}

func (selector *Selector) Select(name Name, src String) []map[Name]SubstringIndex {
	if alternatives, ok := selector.Alternatives[name]; !ok {
		return selector.selectAllMatches(name, selector.Regexps[name], src)
	} else {
		var result []map[Name]SubstringIndex
		for _, alternative := range alternatives {
			result = append(result, selector.selectAllMatches(name, alternative, src)...)
		}
		return result
	}
}
