package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// Color — это код цвета.
type Color int

// Обычные цвета.
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Яркие цвета(работают не во всем терминалах).
const (
	BrightBlack Color = iota + 90
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

// DisplayMode  — это режим отображения виджета.
type DisplayMode int

const (
	DisplayInline  DisplayMode = iota // В одну строку.
	DisplayBlock                      // На отдельной строке.
	DisplayNewLine                    // Перенос строки.
)

// Label — это виджет текстовой метки.
type Label struct {
	decoration string
	Text       string // Текст виджета.
	maxLength  int
	Block      bool // Отображение в блочном режиме.
}

func (l *Label) innerText() string {
	if l.decoration == "" {
		return l.Text
	}
	return fmt.Sprintf("%s%s\033[0m", l.decoration, l.Text)
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
	lbl.decoration += fmt.Sprintf("\033[%dm", clr)
	return lbl
}

// ColorizeBackground() окрашивает фон текста в один из стандартных цветов.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeBackground(clr Color) *Label {
	lbl.decoration += fmt.Sprintf("\033[%dm", clr+10)
	return lbl
}

// ColorizeForegroundRGB() окрашивает текст в RGB.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeForegroundRGB(r, g, b uint8) *Label {
	lbl.decoration += fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	return lbl
}

// ColorizeBackgroundRGB() окрашивает фон текста в RGB.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeBackgroundRGB(r, g, b uint8) *Label {
	lbl.decoration += fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
	return lbl
}

// Bold() делает текст жирным.
func (lbl *Label) Bold() *Label {
	lbl.decoration += "\033[1m"
	return lbl
}

// Italic() делает текст курсивом.
func (lbl *Label) Italic() *Label {
	lbl.decoration += "\033[3m"
	return lbl
}

// Underline() подчеркивает текст.
func (lbl *Label) Underline() *Label {
	lbl.decoration += "\033[4m"
	return lbl
}

// Reverse() реверсирует цвет текста.
func (lbl *Label) Reverse() *Label {
	lbl.decoration += "\033[7m"
	return lbl
}

// Blink() делает текст мигающим(работает не во всем терминалах).
// Добавлено в TUI v1.1.0
func (lbl *Label) Blink() *Label {
	lbl.decoration += "\033[7m"
	return lbl
}

// MaxWidth() реализует интерфейс Component
func (lbl *Label) MaxWidth() int {
	return lbl.maxLength
}

// DisplayMode() реализует интерфейс Component
func (lbl *Label) DisplayMode() DisplayMode {
	if lbl.Block {
		return DisplayBlock
	}
	return DisplayInline
}

// setIndex() реализует интерфейс Component
func (l *Label) setIndex(int) {}

/////////////////////////////////////////////////////////////////////////////////

type spaser struct{}

func (l *spaser) innerText() string {
	return ""
}

func (l *spaser) DisplayMode() DisplayMode {
	return DisplayBlock
}

func (l *spaser) MaxWidth() int {
	return 0
}

func (l *spaser) setIndex(int) {}

// Spaser() создаёт пустой компонент, занимающий всю строку для визуального разделения.
func Spaser() Component { return &spaser{} }

/////////////////////////////////////////////////////////////////////////////////

type newLine struct{}

func (l *newLine) innerText() string {
	return ""
}

func (l *newLine) DisplayMode() DisplayMode {
	return DisplayNewLine
}

func (l *newLine) MaxWidth() int {
	return 0
}

func (l *newLine) setIndex(int) {}

// NewLine() создаёт компонент, занимающий весь остаток его строки, что переносит следующие на новую строку.
func NewLine() Component { return &newLine{} }

// Button это объект кнопки, нажимающейся от нажатия её клавиши. Обработчик в OnClick.
type Button struct {
	clicked Component
	base    Component
	OnClick func()
	Component
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
	btn.Component = btn.base
	currentApp.AddKeyHandler(key, func() {
		if btn.OnClick != nil {
			btn.OnClick()
		}
		btn.Component = btn.clicked
		currentApp.RedrawComponent(btn.idx)
		currentApp.LogInfo("%d %s", btn.idx, btn.innerText())
		time.Sleep(time.Millisecond * 50)
		btn.Component = btn.base
		currentApp.LogInfo("%d %s", btn.idx, btn.innerText())
		currentApp.RedrawComponent(btn.idx)
	})
	return btn
}

func (btn *Button) setIndex(idx int) {
	btn.idx = idx
	btn.base.setIndex(idx)
	btn.clicked.setIndex(idx)
}

// ColorProgress — это виджет шкалы прогресса.
// Добавлено в TUI v1.2.0
type ColorProgress struct {
	text          string
	size          int
	clrOn, clrOff Color
	idx           int
}

// SetValue() устанавливает значение прогресса. Диапазон 0-1
// Добавлено в TUI v1.2.0
func (p *ColorProgress) SetValue(f float64) {
	on := int(float64(p.size) * f)
	p.text = fmt.Sprintf("\033[%dm%s\033[%dm%s\033[0m", p.clrOn+10, strings.Repeat(" ", on), p.clrOff+10, strings.Repeat(" ", p.size-on))
	if currentApp.IsRunned() {
		currentApp.RedrawComponent(p.idx)
	}
}

func (p *ColorProgress) setIndex(idx int) {
	p.idx = idx
}

func (p *ColorProgress) DisplayMode() DisplayMode {
	return DisplayInline
}

func (p *ColorProgress) MaxWidth() int {
	return p.size
}

func (p *ColorProgress) innerText() string {
	return p.text
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
