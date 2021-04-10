package helper

import (
	"regexp"
	"strings"
)

var (
	regexpLineBreak = regexp.MustCompile(`\r?\n`)
	regexpIndent = regexp.MustCompile(`^ +`)
	regexpYAMLCommentWithIndent = regexp.MustCompile(`^( *)# `)
)

func splitLine(s string) []string {
	return regexpLineBreak.Split(s, -1)
}

func lineCount(s string) int {
	return len(splitLine(s))
}

func indent(s string) int {
	return len(regexpIndent.FindString(s))
}

func uncommentYAML(s string) string {
	ss := ""
	for _, line := range splitLine(s) {
		ss = ss + regexpYAMLCommentWithIndent.ReplaceAllString(line, "$1") + "\n"
	}
	return strings.TrimSuffix(ss, "\n")
}