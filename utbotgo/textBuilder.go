package main

import "strings"

type TextBuilder struct {
	Content []interface{}
}

type TextMetaData int

const (
	StartSection TextMetaData = iota
	EndSection
	StartSubsection
	EndSubsection
)

func (tb *TextBuilder) WriteLine(text string) {
	tb.Content = append(tb.Content, text)
}

func (tb *TextBuilder) Section(cont func()) {
	tb.Content = append(tb.Content, StartSection)
	cont()
	tb.Content = append(tb.Content, EndSection)
}

func (tb *TextBuilder) Subsection(begin string, end string, cont func()) {
	if begin != "" {
		tb.Content = append(tb.Content, begin)
	}
	tb.Content = append(tb.Content, StartSubsection)
	cont()
	tb.Content = append(tb.Content, EndSubsection)
	if end != "" {
		tb.Content = append(tb.Content, end)
	}
}

func (tb *TextBuilder) String() string {
	sb := &strings.Builder{}
	tabsAmount := 0
	for i, u := range tb.Content {
		switch unit := u.(type) {
		case string:
			LineWithTabs(sb, tabsAmount, unit)
		case TextMetaData:
			switch unit {
			case StartSection:
			case EndSection:
				if i != len(tb.Content) - 1 && tb.Content[i+1] == StartSection {
					LineWithTabs(sb, 0, "")
				}
			case StartSubsection:
				tabsAmount++
			case EndSubsection:
				tabsAmount--
			}
		}
	}
	return sb.String()
}

func LineWithTabs(sb *strings.Builder, tabsAmount int, text string) {
	for j := 0; j < tabsAmount; j++ {
		sb.WriteString("\t")
	}
	sb.WriteString(text)
	sb.WriteString("\n")
}
