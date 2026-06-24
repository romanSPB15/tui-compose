package tui

import (
	"regexp"
)

const ansi = "[\x1B\x9B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\x07)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

func cleanAnsi(str string) string {
	return regexp.MustCompile(ansi).ReplaceAllString(str, "")
}
