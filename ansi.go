package tui

import (
	"regexp"
	"unicode/utf8"
)

var ansi = regexp.MustCompile(`[\x1B\x9B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\x07)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))`)

func cleanAnsi(str string) string {
	return ansi.ReplaceAllString(str, "")
}

type ansiMatch struct {
	Index int
	Seq   string
}

func findAnsiSequences(s string) ([]ansiMatch, string) {
	byteMatches := ansi.FindAllStringIndex(s, -1)
	if byteMatches == nil {
		return nil, s
	}

	var matches []ansiMatch
	for _, m := range byteMatches {
		startByte, endByte := m[0], m[1]
		startRune := utf8.RuneCountInString(s[:startByte])
		seq := s[startByte:endByte]
		matches = append(matches, ansiMatch{Index: startRune, Seq: seq})
	}
	return matches, cleanAnsi(s)
}
