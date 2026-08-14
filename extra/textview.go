package extra

import (
	"regexp"
	"strings"
)

type TextView struct {
	lines  []string
	offset int
	height int
	width  int
}

func (tv *TextView) ScrollUp() {
	if tv.offset > 0 {
		tv.offset--
	}
}

func (tv *TextView) ScrollDown() {
	if tv.offset < len(tv.lines)-tv.height {
		tv.offset++
	}
}

var replaceList = map[string]string{
	// Обычные цвета текста
	"fg-black":   "30",
	"fg-red":     "31",
	"fg-green":   "32",
	"fg-yellow":  "33",
	"fg-blue":    "34",
	"fg-purpure": "35",
	"fg-cyan":    "36",
	"fg-white":   "37",

	// Обычные цвета фона
	"bg-black":   "40",
	"bg-red":     "41",
	"bg-green":   "42",
	"bg-yellow":  "43",
	"bg-blue":    "44",
	"bg-purpure": "45",
	"bg-cyan":    "46",
	"bg-white":   "47",

	// Яркие цвета текста (коды 90-97)
	"fg-bright-black":   "90",
	"fg-bright-red":     "91",
	"fg-bright-green":   "92",
	"fg-bright-yellow":  "93",
	"fg-bright-blue":    "94",
	"fg-bright-purpure": "95",
	"fg-bright-cyan":    "96",
	"fg-bright-white":   "97",

	// Яркие цвета фона (коды 100-107)
	"bg-bright-black":   "100",
	"bg-bright-red":     "101",
	"bg-bright-green":   "102",
	"bg-bright-yellow":  "103",
	"bg-bright-blue":    "104",
	"bg-bright-purpure": "105",
	"bg-bright-cyan":    "106",
	"bg-bright-white":   "107",

	// Стили
	"bold":      "1",
	"cursive":   "3",
	"italic":    "3",
	"underline": "4",
	"blink":     "5",
	"reverse":   "7",

	// Сброс
	"reset": "0",
	"-":     "0",
}

func NewTextView(h int) *TextView {
	return &TextView{
		height: h,
	}
}

func (tv *TextView) InnerText() string {
	linesANSI := []string{}
	for _, line := range tv.lines {
		for k, v := range replaceList {
			line = strings.ReplaceAll(line, "["+k+"]", "\033["+v+"m")
		}
		linesANSI = append(linesANSI, line)
	}

	end := tv.offset + tv.height
	if end > len(tv.lines) {
		end = len(tv.lines)
	}
	return strings.Join(linesANSI[tv.offset:end], "\n")
}

func (tv *TextView) MaxHeight() int {
	return tv.height
}

func (tv *TextView) WithLines(s []string) *TextView {
	tv.lines = s
	return tv
}

func (tv *TextView) Append(s string) *TextView {
	tv.lines = append(tv.lines, s)
	return tv
}

var tagRegex = regexp.MustCompile(`\[([a-z_-]+)\]`)

func (tv *TextView) MaxWidth() int {
	m := 0
	for _, line := range tv.lines {
		l := len(tagRegex.ReplaceAllString(line, ""))
		if m < l {
			m = l
		}
	}

	return m
}
