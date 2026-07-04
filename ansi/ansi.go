package ansi

import (
	"regexp"
	"unicode/utf8"
)

var ansi = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

func Strip(str string) string {
	return ansi.ReplaceAllString(str, "")
}

type AnsiMatch struct {
	Index int
	Seq   string
}

func Find(s string) ([]AnsiMatch, string) {
	byteMatches := ansi.FindAllStringIndex(s, -1)
	if byteMatches == nil {
		return nil, s
	}

	var matches []AnsiMatch
	for _, m := range byteMatches {
		startByte, endByte := m[0], m[1]
		startRune := utf8.RuneCountInString(s[:startByte])
		seq := s[startByte:endByte]
		matches = append(matches, AnsiMatch{Index: startRune, Seq: seq})
	}
	return matches, Strip(s)
}
