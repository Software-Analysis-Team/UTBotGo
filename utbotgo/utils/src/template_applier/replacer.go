package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type ReplaceInfo struct {
	Index  SubstringIndex
	String String
}

type ReplacesInfo []ReplaceInfo

func (replacesInfo ReplacesInfo) Len() int {
	return len(replacesInfo)
}

func (replacesInfo ReplacesInfo) Less(i, j int) bool {
	return replacesInfo[i].Index.Begin < replacesInfo[j].Index.Begin
}

func (replacesInfo ReplacesInfo) Swap(i, j int) {
	replacesInfo[i], replacesInfo[j] = replacesInfo[j], replacesInfo[i]
}

type Replacer struct {
	Source   String
	Replaces ReplacesInfo
	Name     Name
	Indices  map[Name]SubstringIndex
	Selector *Selector
}

func shortcutString(str String) string {
	if len(str) > 64 {
		return "\"" + string(str[:64]) + "\"..."
	}
	return "\"" + string(str) + "\""
}

func (replacer Replacer) String() string {
	var b strings.Builder
	b.WriteString("{\n")
	b.WriteString(fmt.Sprintf("  Name: %s\n", replacer.Name))
	b.WriteString(fmt.Sprintf("  Source: %s\n", shortcutString(replacer.Source)))
	b.WriteString("  Indices:\n")
	for name, index := range replacer.Indices {
		b.WriteString(fmt.Sprintf("    %s: %s\n", name, shortcutString(replacer.Source[index.Begin:index.End])))
	}
	b.WriteString("}")
	return b.String()
}

func NewReplacer(src String, syntax Syntax) Replacer {
	selector := NewSelector(syntax)
	return Replacer{
		Source: src,
		Selector: &selector,
	}
}

func (replacer Replacer) getAllIndices(name Name) (String, int, []map[Name]SubstringIndex, bool) {
	var (
		selectableName Name
		source String
		offset int
	)
	currIndex, currIndexExists := replacer.Indices[replacer.Name]
	newIndex,  newIndexExists  := replacer.Indices[name]
	subname,   subnameExists   := replacer.Selector.Subnames[replacer.Name]
	multiName := nameEncoder.Encode(name)
	multiNameIndex, multiNameExists := replacer.Indices[multiName]
	switch {
	case replacer.Indices == nil:
		selectableName = name
		source = replacer.Source
		offset = 0
	case ! currIndexExists:
		return nil, 0, nil, false
	case newIndexExists && newIndex.IsSubIndex(currIndex):
		return replacer.Source, 0, []map[Name]SubstringIndex{replacer.Indices}, true
	case multiNameExists:
		selectableName = name
		source = replacer.Source[multiNameIndex.Begin : multiNameIndex.End]
		offset = multiNameIndex.Begin
	case subnameExists:
		selectableName = subname
		source = replacer.Source[currIndex.Begin : currIndex.End]
		offset = currIndex.Begin
	default:
		return nil, 0, nil, false
	}
	selectResult := replacer.Selector.Select(selectableName, source)
	var allIndices []map[Name]SubstringIndex
	if _, ok := replacer.Selector.GeneralIsAlt[selectableName]; ok {
		indices := make(map[Name]SubstringIndex)
		for name_, index_ := range replacer.Indices {
			index_.Begin -= offset
			index_.End   -= offset
			if index_.Begin >= 0 {
				indices[name_] = index_
			}
		}
		allIndices = []map[Name]SubstringIndex{indices}
		allIndices = append(allIndices, selectResult...)
	} else {
		allIndices = selectResult
	}
	return source, offset, allIndices, true
}

func (replacer Replacer) Get(name Name) Replacer {
	newReplacer := replacer
	newReplacer.Name = name
	return newReplacer
}

func (replacer Replacer) GetAll(name Name) []Replacer {
	source, _, allIndices, ok := replacer.getAllIndices(name)
	if !ok {
		_, currIndexExists := replacer.Indices[replacer.Name]
		_, newIndexExists  := replacer.Indices[name]
		if (! currIndexExists) || (! newIndexExists) {
			return []Replacer{replacer.Get(name)}
		} else {
			return nil
		}
	}
	replacers := make([]Replacer, len(allIndices))
	for i, indices := range allIndices {
		replacers[i] = Replacer{
			Source:   source,
			Name:     name,
			Indices:  indices,
			Selector: replacer.Selector,
		}
	}
	return replacers
}

func (replacer Replacer) GetString() String {
	sort.Sort(replacer.Replaces)
	var index SubstringIndex
	if replacer.Name == "" {
		index = SubstringIndex{0, len(replacer.Source)}
	} else {
		index = replacer.Indices[replacer.Name]
	}
	var buf bytes.Buffer
	prevEnd := index.Begin
	var i int
	for i = 0; i < len(replacer.Replaces) && replacer.Replaces[i].Index.Begin < index.Begin; i++ {}
	if i == len(replacer.Replaces) {
		return replacer.Source[index.Begin : index.End]
	}
	for ; i < len(replacer.Replaces) && replacer.Replaces[i].Index.End <= index.End; i++ {
		replace := replacer.Replaces[i]
		buf.Write(replacer.Source[prevEnd : replace.Index.Begin])
		buf.Write(replace.String)
		prevEnd = replace.Index.End
	}
	buf.Write(replacer.Source[prevEnd : index.End])
	return buf.Bytes()
}

func (replacer Replacer) ThereIs(name Name) bool {
	_, ok := replacer.Indices[name]
	return ok
}

func (replacer *Replacer) addReplace(index SubstringIndex, newString String) {
	replacer.Replaces = append(replacer.Replaces, ReplaceInfo{
		Index:  index,
		String: newString,
	})
}

func (replacer Replacer) getSubIndices(name Name) []SubstringIndex {
	_, offset, allIndices, ok := replacer.getAllIndices(name)
	if !ok {
		return nil
	}
	indices := make([]SubstringIndex, len(allIndices))
	for i := range allIndices {
		indices[i] = allIndices[i][name]
		indices[i].Begin += offset
		indices[i].End   += offset
	}
	return indices
}

func (replacer *Replacer) ReplaceString(newString String) {
	replacer.addReplace(replacer.Indices[replacer.Name], newString)
}

func (replacer *Replacer) ReplaceManyStringsAtOnce(newStrings []String, separator String) {
	replacer.addReplace(replacer.Indices[replacer.Name], separator.JoinStrings(newStrings))
}

func (replacer Replacer) ReplaceManyStrings(newStrings []String) []Replacer {
	replacers := make([]Replacer, len(newStrings))
	for i := range replacers {
		replacers[i] = replacer
		replacers[i].ReplaceString(newStrings[i])
	}
	return replacers
}

func (replacer *Replacer) SubReplaceString(subname Name, newString []String) {
	for i, index := range replacer.getSubIndices(subname) {
		replacer.addReplace(index, newString[i])
	}
}

func (replacer *Replacer) SubReplaceManyStringsAtOnce(subname Name, newStrings [][]String, separator String) {
	for i, index := range replacer.getSubIndices(subname) {
		replacer.addReplace(index, separator.JoinStrings(newStrings[i]))
	}
}

func (replacer Replacer) SubReplaceManyStrings(subname Name, newStrings [][]String) []Replacer {
	if len(newStrings) == 0 {
		return []Replacer{replacer}
	}
	indices := replacer.getSubIndices(subname)
	replacers := make([]Replacer, len(newStrings[0]))
	for j := range replacers {
		replacers[j] = replacer
		for i, index := range indices {
			replacers[j].addReplace(index, newStrings[i][j])
		}
	}
	return replacers
}

func (replacer *Replacer) Replace(srcReplacer Replacer) {
	replacer.ReplaceString(srcReplacer.GetString())
}

func (replacer *Replacer) ReplaceManyAtOnce(srcReplacers []Replacer, separator String) {
	newStrings := make([]String, len(srcReplacers))
	for i := range newStrings {
		newStrings[i] = srcReplacers[i].GetString()
	}
	replacer.ReplaceManyStringsAtOnce(newStrings, separator)
}

func (replacer Replacer) ReplaceMany(srcReplacers []Replacer) []Replacer {
	newStrings := make([]String, len(srcReplacers))
	for i := range newStrings {
		newStrings[i] = srcReplacers[i].GetString()
	}
	return replacer.ReplaceManyStrings(newStrings)
}

func (replacer *Replacer) SubReplace(subname Name, srcReplacer []Replacer) {
	newString := make([]String, len(srcReplacer))
	for i := range newString {
		newString[i] = srcReplacer[i].GetString()
	}
	replacer.SubReplaceString(subname, newString)
}

func (replacer *Replacer) SubReplaceManyAtOnce(subname Name, srcReplacers [][]Replacer, separator String) {
	newStrings := make([][]String, len(srcReplacers))
	for i := range newStrings {
		newStrings[i] = make([]String, len(srcReplacers[i]))
		for j := range newStrings[i] {
			newStrings[i][j] = srcReplacers[i][j].GetString()
		}
	}
	replacer.SubReplaceManyStringsAtOnce(subname, newStrings, separator)
}

func (replacer Replacer) SubReplaceMany(subname Name, srcReplacers [][]Replacer) []Replacer {
	newStrings := make([][]String, len(srcReplacers))
	for i := range newStrings {
		newStrings[i] = make([]String, len(srcReplacers[i]))
		for j := range newStrings[i] {
			newStrings[i][j] = srcReplacers[i][j].GetString()
		}
	}
	return replacer.SubReplaceManyStrings(subname, newStrings)
}

func (replacer *Replacer) UpdateReplaces(srcReplacer *Replacer) {
	if len(replacer.Source) == len(srcReplacer.Source) {
		replacer.Replaces = append(replacer.Replaces, srcReplacer.Replaces...)
	} else {
		offset := replacer.Indices[replacer.Name].Begin
		for _, replace := range srcReplacer.Replaces {
			index := replace.Index
			index.Begin += offset
			index.End += offset
			newString := replace.String
			replacer.addReplace(index, newString)
		}
	}
	srcReplacer.Replaces = nil
}
