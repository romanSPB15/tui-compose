package tui

import (
	"strconv"

	"github.com/romanSPB15/tui-compose/v3/cell"
)

type TitlePosition int

const (
	TitleTopLeft TitlePosition = iota
	TitleTopCenter
	TitleTopRight
	TitleBottomLeft
	TitleBottomCenter
	TitleBottomRight
)

// Frame – контейнер, который оборачивает содержимое в рамку.
type Frame struct {
	border border
}

// Title – заголовок рамки.
type Title struct {
	Pos   TitlePosition
	Text  string
	Style Style
}

type border struct {
	content        Widget
	tl, tr, bl, br rune
	h, v           rune
	ph, pv         int

	borderStyle Style
	titles      []Title
}

// NewFrame создаёт рамку.
func NewFrame(content Widget) *Frame {
	return &Frame{
		border: border{
			content: content,
			tl:      '┌', tr: '┐', bl: '└', br: '┘',
			h: '─', v: '│',
			ph: 1,
			pv: 0,
		},
	}
}

// Custom – вручную задать символы рамки.
func (b *Frame) Custom(tl, tr, bl, br, h, v rune) *Frame {
	b.border.tl, b.border.tr, b.border.bl, b.border.br = tl, tr, bl, br
	b.border.h, b.border.v = h, v
	return b
}

// Heavy – толстая рамка.
func (b *Frame) Heavy() *Frame {
	return b.Custom('┏', '┓', '┗', '┛', '━', '┃')
}

// ASCII – ASCII-совместимая рамка (подходит для старых терминалов).
func (b *Frame) ASCII() *Frame {
	return b.Custom('+', '+', '+', '+', '-', '|')
}

// Rounded – закруглённые углы.
func (b *Frame) Rounded() *Frame {
	return b.Custom('╭', '╮', '╰', '╯', '─', '│')
}

// Double – двойная рамка.
func (b *Frame) Double() *Frame {
	return b.Custom('╔', '╗', '╚', '╝', '═', '║')
}

// Default – рамка по умолчанию.
func (b *Frame) Default() *Frame {
	return b.Custom('┌', '┐', '└', '┘', '─', '│')
}

// Dashed – рамка пунктиром.
func (b *Frame) Dashed() *Frame {
	return b.Custom('┌', '┐', '└', '┘', '┄', '┆')
}

// RoundedDashed – закруглённые углы & пунктир.
func (b *Frame) RoundedDashed() *Frame {
	return b.Custom('╭', '╮', '╰', '╯', '┄', '┆')
}

// Bevel - скошенные углы.
func (b *Frame) Bevel() *Frame {
	return b.Custom('╱', '╲', '╲', '╱', '─', '│')
}

// Bevel - ASCII-совместимые скошенные углы (подходит для старых терминалов).
func (b *Frame) BevelASCII() *Frame {
	return b.Custom('/', '\\', '\\', '/', '─', '│')
}

// WithTitle добавляет заголовок. Можно вызвать несколько раз, чтобы добавить несколько.
func (b *Frame) WithTitle(t Title) *Frame {
	b.border.titles = append(b.border.titles, t)
	return b
}

// WithBorderStyle устанавливает стиль рамки.
func (b *Frame) WithBorderStyle(s Style) *Frame {
	b.border.borderStyle = s
	return b
}

func (b *border) MaxWidth() int {
	if b.content == nil {
		return b.ph*2 + 2
	}
	return b.content.MaxWidth() + b.ph*2 + 2
}

func (b *border) MaxHeight() int {
	if b.content == nil {
		return b.pv*2 + 2
	}
	return b.content.MaxHeight() + b.pv*2 + 2
}

func (b *Frame) MaxWidth() int {
	return b.border.MaxWidth()
}

func (b *Frame) MaxHeight() int {
	return b.border.MaxHeight()
}

func convertToCellStyle(s Style) cell.Style {
	var cs cell.Style
	fg := int(s & 0x1F)
	if fg != 0 {
		if fg <= 8 {
			cs.Fg = strconv.Itoa(fg + 29)
		} else {
			cs.Fg = strconv.Itoa(fg + 81)
		}
	}
	bg := int((s >> 5) & 0x1F)
	if bg != 0 {
		if bg <= 8 {
			cs.Bg = strconv.Itoa(bg + 39)
		} else {
			cs.Bg = strconv.Itoa(bg + 91)
		}
	}

	if s&Bold != 0 {
		cs.Args |= cell.Bold
	}
	if s&Italic != 0 {
		cs.Args |= cell.Italic
	}
	if s&Underline != 0 {
		cs.Args |= cell.Underline
	}
	if s&Blink != 0 {
		cs.Args |= cell.Blink
	}
	if s&Reverse != 0 {
		cs.Args |= cell.Reverse
	}
	if s&Reset != 0 {
		cs.Args |= cell.Reset
	}
	return cs
}

func (b *border) InnerText() string {
	w := b.MaxWidth()
	h := b.MaxHeight()

	cells := make([][]cell.Cell, h)
	for i := range cells {
		cells[i] = make([]cell.Cell, w)
		for j := range cells[i] {
			cells[i][j] = cell.Cell{Char: ' '}
		}
	}

	borderStyle := convertToCellStyle(b.borderStyle)

	// углы
	cells[0][0] = cell.Cell{b.tl, borderStyle}     // верхний левый
	cells[0][w-1] = cell.Cell{b.tr, borderStyle}   // верхний правый
	cells[h-1][0] = cell.Cell{b.bl, borderStyle}   // нижний левый
	cells[h-1][w-1] = cell.Cell{b.br, borderStyle} // нижний правый

	// линии
	for i := 1; i < w-1; i++ {
		cells[0][i] = cell.Cell{b.h, borderStyle}
		cells[h-1][i] = cell.Cell{b.h, borderStyle}
	}

	for i := 1; i < h-1; i++ {
		cells[i][0] = cell.Cell{b.v, borderStyle}
		cells[i][w-1] = cell.Cell{b.v, borderStyle}
	}

	for _, t := range b.titles {
		titleStyle := convertToCellStyle(t.Style)

		titleRunes := []rune(t.Text)

		drawTitle := func(x, y int) {
			for i, r := range titleRunes {
				cells[y][x+i] = cell.Cell{r, titleStyle}
			}
		}

		switch t.Pos {
		case TitleTopLeft:
			drawTitle(2, 0)
		case TitleTopRight:
			drawTitle(w-len(titleRunes)-2, 0)
		case TitleTopCenter:
			drawTitle(w/2-len(titleRunes)/2, 0)
		case TitleBottomLeft:
			drawTitle(2, h-1)
		case TitleBottomRight:
			drawTitle(w-len(titleRunes)-2, h-1)
		case TitleBottomCenter:
			drawTitle(w/2-len(titleRunes)/2, h-1)
		}
	}

	return cell.ToString(cells)
}

// WithPaddings задаёт отступы внутри рамки.
// по умолчанию v=0 h=1
func (b *Frame) WithPaddings(v, h int) *Frame {
	b.border.ph, b.border.pv = h, v
	return b
}

func (b *Frame) InnerText() string {
	return ""
}

func (b *Frame) Child() []Widget {
	return []Widget{&b.border, b.border.content}
}

func (b *Frame) Pos(i int) Pos {
	if i == 0 {
		return Pos{0, 0}
	}
	return Pos{b.border.pv + 1, b.border.ph + 1}
}
