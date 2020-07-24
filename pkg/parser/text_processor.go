package parser

import (
	"regexp"
	"strings"
)

type TextProcessor interface {
	removeCursive(text string) string
	removeComments(text string) string
	removeHTML(text string) string
	removeLists(text string) string
	removeStrong(text string) string
	removeInternetRefs(text string) string

	processFigureBrackets(text string) string
	processRefs(text string) (string, []string)

	ProcessText(text string) (string, []string)
}

type WikiTextProcessor struct {
	titlesRe            *regexp.Regexp
	refsRe              *regexp.Regexp
	listsRe             *regexp.Regexp
	figureBracketsRe    *regexp.Regexp
	commentsRe          *regexp.Regexp
	htmlRe              *regexp.Regexp
	cursiveRe           *regexp.Regexp
	strongRe            *regexp.Regexp
	internetRefsRe      *regexp.Regexp
	multipleLinesRefsRe *regexp.Regexp
}

func NewWikiTextProcessor() *WikiTextProcessor {
	return &WikiTextProcessor{
		titlesRe:            regexp.MustCompile(`== (.*?) ==`),
		refsRe:              regexp.MustCompile(`\[\[(.*?)]]`),
		internetRefsRe:      regexp.MustCompile(`\[(.*?)]`),
		listsRe:             regexp.MustCompile(`{\|(.*?)\|}`),
		figureBracketsRe:    regexp.MustCompile(`{{(.*?)}}`),
		htmlRe:              regexp.MustCompile(`<(.*?)>`),
		commentsRe:          regexp.MustCompile(`<!--(.*?)-->`),
		cursiveRe:           regexp.MustCompile(`''(.*?)''`),
		strongRe:            regexp.MustCompile(`'''(.*?)'''`),
		multipleLinesRefsRe: regexp.MustCompile(`\|(.*?)}}`),
	}	
}

func (w *WikiTextProcessor) GetTitles(text string) []string {
	titlesSlice := []string{"Заголовок"}
	for _, title := range w.titlesRe.FindAllString(text, -1) {
		title := strings.Replace(title, ".", " ", -1)
		title = strings.Trim(title, "=")
		titlesSlice = append(titlesSlice, strings.TrimSpace(title))
	}
	return titlesSlice
}

func (w *WikiTextProcessor) SplitText(text string) []string {
	return w.titlesRe.Split(text, -1)
}

func (w *WikiTextProcessor) processRefs(text string) (string, []string) {
	processedText := text
	var refsSlice []string
	var textRef, externalRef string
	for _, matchStr := range w.refsRe.FindAllString(text, -1) {
		bufSlice := strings.Split(matchStr[2:len(matchStr) - 2], "|")
		if len(bufSlice) == 2 {
			externalRef, textRef = bufSlice[0], bufSlice[1]
			refsSlice = append(refsSlice, externalRef)
		} else if len(bufSlice) == 1 {
			externalRef, textRef = bufSlice[0], bufSlice[0]
			refsSlice = append(refsSlice, externalRef)
		} else if strings.Contains(matchStr, "Файл:") {
			textRef = ""
		}
		processedText = strings.Replace(processedText, matchStr, textRef, 1)
	}
	return w.multipleLinesRefsRe.ReplaceAllString(processedText, ""), refsSlice
}

func (w *WikiTextProcessor) removeLists(text string) string {
	return w.listsRe.ReplaceAllString(text, "")
}

func (w *WikiTextProcessor) processFigureBrackets(text string) string {
	processedText := text
	for _, matchStr := range w.figureBracketsRe.FindAllString(text, -1) {
		bufSlice := strings.Split(matchStr[2:len(matchStr) - 2], "|")
		if len(bufSlice) == 2 {
			processedText = strings.Replace(processedText, matchStr, bufSlice[1], 1)
		} else {
			processedText = strings.Replace(processedText, matchStr, "", 1)
		}
	}
	return processedText
}

func (w *WikiTextProcessor) removeCursive(text string) string {
	processedText := text
	for _, matchStr := range w.cursiveRe.FindAllString(text, -1) {
		processedText = strings.Replace(processedText, matchStr, matchStr[2:len(matchStr) - 2], 1)
	}
	return processedText
}

func (w *WikiTextProcessor) removeStrong(text string) string {
	processedText := text
	for _, matchStr := range w.strongRe.FindAllString(text, -1) {
		processedText = strings.Replace(processedText, matchStr, matchStr[3:len(matchStr) - 3], 1)
	}
	return processedText
}

func (w *WikiTextProcessor) removeComments(text string) string {
	return w.commentsRe.ReplaceAllString(text, "")
}

func (w *WikiTextProcessor) removeHTML(text string) string {
	return w.htmlRe.ReplaceAllString(text, "")
}

func (w *WikiTextProcessor) removeMultipleLinesRefs(text string) string {
	return w.multipleLinesRefsRe.ReplaceAllString(text, "")
}

func (w *WikiTextProcessor) removeInternetRefs(text string) string {
	return w.internetRefsRe.ReplaceAllString(text, "")
}

func (w *WikiTextProcessor) ProcessText(text string) (string, []string) {
	processedText := strings.Replace(text, "\n", "", -1)
	processedText = strings.Replace(processedText, "\a0", " ", -1)
	processedText = strings.Replace(processedText, "\u0301", "", -1)
	processedText = strings.Trim(processedText, "=")

	processedText = w.removeStrong(processedText)
	processedText = w.removeCursive(processedText)
	processedText = w.removeComments(processedText)
	processedText = w.removeHTML(processedText)
	processedText = w.processFigureBrackets(processedText)

	processedText, refsSlice := w.processRefs(processedText)

	processedText = w.removeLists(processedText)
	processedText = w.removeInternetRefs(processedText)

	return processedText, refsSlice
}
