package tui

import (
	"strconv"
	"strings"
)

// Style — битовая маска стиля виджета.
type Style uint16

// Цвета (0–16)
const (
	FrDefault = iota
	FrBlack
	FrRed
	FrGreen
	FrYellow
	FrBlue
	FrMagenta
	FrCyan
	FrWhite
	FrBrightBlack
	FrBrightRed
	FrBrightGreen
	FrBrightYellow
	FrBrightBlue
	FrBrightMagenta
	FrBrightCyan
	FrBrightWhite
)

// Атрибуты
const (
	Bold      Style = 1 << 10
	Italic    Style = 1 << 11
	Underline Style = 1 << 12
	Blink     Style = 1 << 13
	Reverse   Style = 1 << 14
	Reset     Style = 1 << 15
)

const (
	BgDefault       Style = FrDefault << 5
	BgBlack         Style = FrBlack << 5
	BgRed           Style = FrRed << 5
	BgGreen         Style = FrGreen << 5
	BgYellow        Style = FrYellow << 5
	BgBlue          Style = FrBlue << 5
	BgMagenta       Style = FrMagenta << 5
	BgCyan          Style = FrCyan << 5
	BgWhite         Style = FrWhite << 5
	BgBrightBlack   Style = FrBrightBlack << 5
	BgBrightRed     Style = FrBrightRed << 5
	BgBrightGreen   Style = FrBrightGreen << 5
	BgBrightYellow  Style = FrBrightYellow << 5
	BgBrightBlue    Style = FrBrightBlue << 5
	BgBrightMagenta Style = FrBrightMagenta << 5
	BgBrightCyan    Style = FrBrightCyan << 5
	BgBrightWhite   Style = FrBrightWhite << 5
)

func (s Style) String() string {
	if s == 0 {
		return ""
	}

	var codes []int

	fg := int(s & 0x1F)
	if fg != 0 {
		if fg <= 8 {
			codes = append(codes, fg+29)
		} else {
			codes = append(codes, fg+81)
		}
	}

	bg := int((s >> 5) & 0x1F)
	if bg != 0 {
		if bg <= 8 {
			codes = append(codes, bg+39)
		} else {
			codes = append(codes, bg+91)
		}
	}

	if s&Bold != 0 {
		codes = append(codes, 1)
	}
	if s&Italic != 0 {
		codes = append(codes, 3)
	}
	if s&Underline != 0 {
		codes = append(codes, 4)
	}
	if s&Blink != 0 {
		codes = append(codes, 5)
	}
	if s&Reverse != 0 {
		codes = append(codes, 7)
	}
	if s&Reset != 0 {
		codes = []int{0}
	}

	if len(codes) == 0 {
		return ""
	}

	codesString := []string{}

	for _, v := range codes {
		codesString = append(codesString, strconv.Itoa(v))
	}

	return "\x1b[" + strings.Join(codesString, ";") + "m"
}
