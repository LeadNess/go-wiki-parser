package parser

import (
	"regexp"
)

type TextProcessor interface {
	removeCursive(text string) error
	removeHTML(text string) error
	removeComments(text string) error
	removeLists(text string) error
	removeStrong(text string) error
	removeMultipleLinesRefs(text string) error
	processFigureBrackets(text string) error
	processRefs(text string) ([]string, error)
}

type WikiParser struct {
	titlesRe            *regexp.Regexp
	refsRe              *regexp.Regexp
	listsRe             *regexp.Regexp
	figureBracketsRe    *regexp.Regexp
	htmlRe              *regexp.Regexp
	cursiveRe           *regexp.Regexp
	strongRe            *regexp.Regexp
	multipleLinesRefsRe *regexp.Regexp
}

func NewWikiParser() *WikiParser {
	return &WikiParser{
		titlesRe:            regexp.MustCompile(`== [\w ]+ ==`),
		refsRe:              regexp.MustCompile(`\[\[(.*?)\]\]`),
		listsRe:             regexp.MustCompile(`{\|[\S\n ]+\|}`),
		figureBracketsRe:    regexp.MustCompile(`{{(.*?)}}`),
		htmlRe:              regexp.MustCompile(`<(.*?)>`),
		cursiveRe:           regexp.MustCompile(`\'\'(.*?)\'\'`),
		strongRe:            regexp.MustCompile(`\'\'\'(.*?)\'\'\'`),
		multipleLinesRefsRe: regexp.MustCompile(`{{[^}]+}}`),
	}	
}
