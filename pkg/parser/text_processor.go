package parser

import (
	"regexp"
	"strings"
)

type TextProcessor interface {
	removeCursive(text string) string
	removeHTML(text string) string
	removeLists(text string) string
	removeStrong(text string) string
	removeMultipleLinesRefs(text string) string

	processFigureBrackets(text string) string
	processRefs(text string) (string, []string)

	ProcessText(text string) (string, []string)
}

type WikiTextProcessor struct {
	titlesRe            *regexp.Regexp
	refsRe              *regexp.Regexp
	listsRe             *regexp.Regexp
	figureBracketsRe    *regexp.Regexp
	htmlRe              *regexp.Regexp
	cursiveRe           *regexp.Regexp
	strongRe            *regexp.Regexp
	multipleLinesRefsRe *regexp.Regexp
}

func NewWikiParser() *WikiTextProcessor {
	return &WikiTextProcessor{
		titlesRe:            regexp.MustCompile(`== [\w ]+ ==`),
		refsRe:              regexp.MustCompile(`\[\[(.*?)]]`),
		listsRe:             regexp.MustCompile(`{\|[\S\n ]+\|}`),
		figureBracketsRe:    regexp.MustCompile(`{{(.*?)}}`),
		htmlRe:              regexp.MustCompile(`<(.*?)>`),
		cursiveRe:           regexp.MustCompile(`''(.*?)''`),
		strongRe:            regexp.MustCompile(`'''(.*?)'''`),
		multipleLinesRefsRe: regexp.MustCompile(`{{[^}]+}}`),
	}	
}

func (w *WikiTextProcessor) GetTitles(text string) []string {
	return append([]string{"Заголовок"}, w.titlesRe.FindAllString(text, -1)...)
}

func (w *WikiTextProcessor) SplitText(text string) []string {
	return w.titlesRe.Split(text, -1)
}

func (w *WikiTextProcessor) processRefs(text string) (string, []string) {
	var refsSlice []string
	var processedText, textRef, externalRef string
	for _, matchStr := range w.refsRe.FindAllString(text, -1) {
		bufSlice := strings.Split(matchStr, "|")
		if len(bufSlice) == 2 {
			externalRef, textRef = bufSlice[0], bufSlice[1]
			refsSlice = append(refsSlice, externalRef)
		} else if len(bufSlice) == 1 {
			externalRef, textRef = bufSlice[0], bufSlice[0]
			refsSlice = append(refsSlice, externalRef)
		} else if strings.Contains(matchStr, "Файл:") {
			textRef = ""
		}
		processedText = strings.Replace(processedText, "[[" + matchStr + "]]", textRef, 1)
	}
	return processedText, refsSlice
}

func (w *WikiTextProcessor) removeLists(text string) string {
	var processedText string
	for _, matchStr := range w.listsRe.FindAllString(text, -1) {
		processedText = strings.Replace(processedText, matchStr, "", 1)
	}
	return processedText
}

func (w *WikiTextProcessor) processFigureBrackets(text string) string {
	var processedText string
	for _, matchStr := range w.figureBracketsRe.FindAllString(text, -1) {
		bufSlice := strings.Split(matchStr, "|")
		if len(bufSlice) == 2 {
			processedText = strings.Replace(processedText, "{{" + matchStr + "}}", bufSlice[1], 1)
		} else {
			processedText = strings.Replace(processedText, "{{" + matchStr + "}}", "", 1)
		}
	}
	return processedText
}

func (w *WikiTextProcessor) removeCursive(text string) string {
	var processedText string
	for _, matchStr := range w.cursiveRe.FindAllString(text, -1) {
		processedText = strings.Replace(processedText, `''` + matchStr + `''`, matchStr, 1)
	}
	return processedText
}

func (w *WikiTextProcessor) removeStrong(text string) string {
	var processedText string
	for _, matchStr := range w.strongRe.FindAllString(text, -1) {
		processedText = strings.Replace(processedText, `'''` + matchStr + `'''`, matchStr, 1)
	}
	return processedText
}

func (w *WikiTextProcessor) removeHTML(text string) string {
	var processedText string
	for _, matchStr := range w.htmlRe.FindAllString(text, -1) {
		processedText = strings.Replace(processedText, `<` + matchStr + `>`, "", 1)
	}
	return processedText
}

func (w *WikiTextProcessor) removeMultipleLinesRefs(text string) string {
	var processedText string
	for _, matchStr := range w.multipleLinesRefsRe.FindAllString(text, -1) {
		processedText = strings.Replace(processedText, matchStr, "", 1)
	}
	return processedText
}

func (w *WikiTextProcessor) ProcessText(text string) (string, []string) {
	processedText := w.removeCursive(text)
	processedText = w.removeHTML(text)
	processedText = w.removeStrong(text)
	processedText = w.processFigureBrackets(text)

	processedText = strings.Replace(processedText, "\n", "", -1)
	processedText = strings.Replace(processedText, "\xa0", " ", -1)

	processedText = w.removeLists(text)
	processedText = w.removeMultipleLinesRefs(text)
	return w.processRefs(processedText)
}
