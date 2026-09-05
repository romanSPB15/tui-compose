package extra

import (
	"regexp"
	"strings"

	"github.com/romanSPB15/tui-compose/v4/cell"
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
	"fg-black":   "30",
	"fg-red":     "31",
	"fg-green":   "32",
	"fg-yellow":  "33",
	"fg-blue":    "34",
	"fg-purpure": "35",
	"fg-magenta": "35",
	"fg-cyan":    "36",
	"fg-white":   "37",

	"bg-black":   "40",
	"bg-red":     "41",
	"bg-green":   "42",
	"bg-yellow":  "43",
	"bg-blue":    "44",
	"bg-purpure": "45",
	"bg-magenta": "45",
	"bg-cyan":    "46",
	"bg-white":   "47",

	"fg-bright-black":   "90",
	"fg-bright-red":     "91",
	"fg-bright-green":   "92",
	"fg-bright-yellow":  "93",
	"fg-bright-blue":    "94",
	"fg-bright-purpure": "95",
	"fg-bright-magenta": "95",
	"fg-bright-cyan":    "96",
	"fg-bright-white":   "97",

	"bg-bright-black":   "100",
	"bg-bright-red":     "101",
	"bg-bright-green":   "102",
	"bg-bright-yellow":  "103",
	"bg-bright-blue":    "104",
	"bg-bright-purpure": "105",
	"bg-bright-magenta": "105",
	"bg-bright-cyan":    "106",
	"bg-bright-white":   "107",

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

func (tv *TextView) Render(buf [][]cell.Cell) {
	start := tv.offset
	end := start + tv.height
	if end > len(tv.lines) {
		end = len(tv.lines)
	}
	w := tv.width
	if w == 0 {
		w = tv.Width()
	}

	var cellBuf []cell.Cell

	for y, lineIdx := start, 0; y < end && y < len(buf); y, lineIdx = y+1, lineIdx+1 {
		line := tv.lines[lineIdx]

		ansiLine := line
		for k, v := range replaceList {
			ansiLine = strings.ReplaceAll(ansiLine, "["+k+"]", "\033["+v+"m")
		}

		cells, _ := cell.ParseFromTo(ansiLine, &cellBuf, cell.Style{})

		row := buf[y]
		for i, c := range cells {
			if i >= w {
				break
			}
			row[i] = c
		}

		for i := len(cells); i < w && i < len(row); i++ {
			row[i] = cell.Cell{Char: ' ', Style: cell.Style{}}
		}
	}
}

func (tv *TextView) Height() int {
	return tv.height
}

func (tv *TextView) WithLines(s []string) *TextView {
	tv.lines = s
	return tv
}

func (tv *TextView) WithFixedWidth(w int) *TextView {
	tv.width = w
	return tv
}

func (tv *TextView) Append(s string) *TextView {
	tv.lines = append(tv.lines, s)
	return tv
}

var tagRegex = regexp.MustCompile(`\[([a-z_-]+)\]`)

func (tv *TextView) Width() int {
	if tv.width != 0 {
		return tv.width
	}
	m := 0
	for _, line := range tv.lines {
		cleaned := tagRegex.ReplaceAllStringFunc(line, func(match string) string {
			key := match[1 : len(match)-1]
			if _, ok := replaceList[key]; ok {
				return ""
			}
			return match
		})
		l := len([]rune(cleaned))
		if l > m {
			m = l
		}
	}
	return m
}
