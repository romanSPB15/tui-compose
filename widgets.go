//go:build !no_widgets

package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// Label — это виджет текстовой метки.
type Label struct {
	ANSI      string // Приставка ANSI escape последовательности
	Text      string // Текст виджета.
	maxLength int
	Block     bool // Отображение в блочном режиме.
}

func (l *Label) InnerText() string {
	if l.ANSI == "" {
		return l.Text
	}
	return fmt.Sprintf("%s%s\033[0m", l.ANSI, l.Text)
}

// NewStaticLabel() создаёт виджет текста.
func NewStaticLabel(txt string) *Label { return &Label{Text: txt, maxLength: len(txt)} }

// NewDynamicLabel() создаёт виджет текста с возможностью изменения содержимого в будущем.
// maxLength это место, зарезервированное под метку в символах.
func NewDynamicLabel(txt string, maxLength int) *Label {
	return &Label{Text: txt, maxLength: maxLength}
}

// ColorizeForeground() окрашивает текст в один из стандартных цветов.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeForeground(clr Color) *Label {
	lbl.ANSI += fmt.Sprintf("\033[%dm", clr)
	return lbl
}

// ColorizeBackground() окрашивает фон текста в один из стандартных цветов.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeBackground(clr Color) *Label {
	lbl.ANSI += fmt.Sprintf("\033[%dm", clr+10)
	return lbl
}

// ColorizeForegroundRGB() окрашивает текст в RGB.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeForegroundRGB(r, g, b uint8) *Label {
	lbl.ANSI += fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	return lbl
}

// ColorizeBackgroundRGB() окрашивает фон текста в RGB.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeBackgroundRGB(r, g, b uint8) *Label {
	lbl.ANSI += fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
	return lbl
}

// Bold() делает текст жирным.
func (lbl *Label) Bold() *Label {
	lbl.ANSI += "\033[1m"
	return lbl
}

// Italic() делает текст курсивом.
func (lbl *Label) Italic() *Label {
	lbl.ANSI += "\033[3m"
	return lbl
}

// Underline() подчеркивает текст.
func (lbl *Label) Underline() *Label {
	lbl.ANSI += "\033[4m"
	return lbl
}

// Reverse() реверсирует цвет текста.
func (lbl *Label) Reverse() *Label {
	lbl.ANSI += "\033[7m"
	return lbl
}

// Reset() убирает все декорации текста.
// Добавлено в TUI v1.5.0
func (lbl *Label) Reset() *Label {
	lbl.ANSI = ""
	return lbl
}

// Blink() делает текст мигающим(работает не во всем терминалах).
// Добавлено в TUI v1.1.0
func (lbl *Label) Blink() *Label {
	lbl.ANSI += "\033[7m"
	return lbl
}

// MaxLength() реализует интерфейс Widget
func (lbl *Label) MaxLength() int {
	return lbl.maxLength
}

// DisplayMode() реализует интерфейс Widget
func (lbl *Label) DisplayMode() DisplayMode {
	if lbl.Block {
		return DisplayBlock
	}
	return DisplayInline
}

// SetIndex() реализует интерфейс Widget
func (l *Label) SetIndex(int) {}

/////////////////////////////////////////////////////////////////////////////////

type spaser struct{}

func (l *spaser) InnerText() string {
	return ""
}

func (l *spaser) DisplayMode() DisplayMode {
	return DisplayBlock
}

func (l *spaser) MaxLength() int {
	return 0
}

func (l *spaser) SetIndex(int) {}

// Spaser() создаёт пустой компонент, занимающий всю строку для визуального разделения.
func Spaser() Widget { return &spaser{} }

/////////////////////////////////////////////////////////////////////////////////

type newLine struct{}

func (l *newLine) InnerText() string {
	return ""
}

func (l *newLine) DisplayMode() DisplayMode {
	return DisplayNewLine
}

func (l *newLine) MaxLength() int {
	return 0
}

func (l *newLine) SetIndex(int) {}

// NewLine() создаёт компонент, занимающий весь остаток его строки, что переносит следующие на новую строку.
func NewLine() Widget { return &newLine{} }

// Button это объект кнопки, нажимающейся от нажатия её клавиши. Обработчик в OnClick.
type Button struct {
	clicked Widget
	base    Widget
	OnClick func()
	Widget
	idx int
}

// NewButton() создаёт кнопку.
// key это её клавиша.
func NewButton(text string, key keyboard.Key) *Button {
	btn := &Button{
		clicked: NewStaticLabel(text).ColorizeForeground(Blue),
		base:    NewStaticLabel(text),
		OnClick: func() {},
	}
	btn.Widget = btn.base
	currentWindow.RegisterKeyHandler(key, func() {
		if btn.OnClick != nil {
			btn.OnClick()
		}
		btn.Widget = btn.clicked
		currentWindow.RedrawWidget(btn.idx)
		currentWindow.LogInfo("%d %s", btn.idx, btn.InnerText())
		time.Sleep(time.Millisecond * 50)
		btn.Widget = btn.base
		currentWindow.LogInfo("%d %s", btn.idx, btn.InnerText())
		currentWindow.RedrawWidget(btn.idx)
	})
	return btn
}

func (btn *Button) SetIndex(idx int) {
	btn.idx = idx
	btn.base.SetIndex(idx)
	btn.clicked.SetIndex(idx)
}

// ColorProgress — это виджет шкалы прогресса.
// Добавлено в TUI v1.2.0
type ColorProgress struct {
	text          string
	size          int
	clrOn, clrOff Color
	idx           int
}

// NewColorProgress() cоздаёт виджет шкалы прогресса.
// len — максимальная длина в пикселях
// on — цвет "включенных" пикселей
// off — цвет "выключенных" пикселей
// Добавлено в TUI v1.2.0
func NewColorProgress(len int, on, off Color) *ColorProgress {
	return &ColorProgress{
		text:   strings.Repeat(" ", len),
		size:   len,
		clrOn:  on,
		clrOff: off,
	}
}

// SetValue() устанавливает значение прогресса. Диапазон 0-1
// Добавлено в TUI v1.2.0
func (p *ColorProgress) SetValue(f float64) {
	on := int(float64(p.size) * f)
	p.text = fmt.Sprintf("\033[%dm%s\033[%dm%s\033[0m", p.clrOn+10, strings.Repeat(" ", on), p.clrOff+10, strings.Repeat(" ", p.size-on))
	if currentWindow.IsRunned() {
		currentWindow.RedrawWidget(p.idx)
	}
}

func (p *ColorProgress) SetIndex(idx int) {
	p.idx = idx
}

func (p *ColorProgress) DisplayMode() DisplayMode {
	return DisplayInline
}

func (p *ColorProgress) MaxLength() int {
	return p.size
}

func (p *ColorProgress) InnerText() string {
	return p.text
}
