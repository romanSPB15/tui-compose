package main

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/builder"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose - Slider Demo")

	// 1. Простой слайдер [0, 1]
	basic := extra.NewSlider(25)

	// 2. Кастомные символы и стили
	styled := extra.NewSlider(25).
		WithThumbRune('◆').
		WithTrackRune('·').
		WithFillRune('━').
		WithTrackStyle(tui.FrBrightBlack).
		WithFillStyle(tui.FrGreen).
		WithThumbStyle(tui.FrBrightGreen).
		WithValue(0.7)

	// 3. Диапазон, шаг и своя подпись
	ranged := extra.NewSlider(25).
		WithRange(0, 100).
		WithStep(4).
		WithValue(42).
		WithLabelFunc(func(v float64) string {
			b := builder.Builder{}
			b.WriteInt(int(v))
			b.WriteRune(' ')
			b.WriteRune('°')
			b.WriteRune('C')
			return b.String()
		}).
		WithFillStyle(tui.FrYellow).
		WithThumbStyle(tui.FrBrightYellow)

	// 4. ASCII-пресет
	ascii := extra.NewSlider(25).ASCII().
		WithRange(0, 1).
		WithValue(0.6).
		WithFillStyle(tui.FrMagenta).
		WithThumbStyle(tui.FrMagenta)

	wnd.SetContent(tui.NewFrame(tui.NewVBox(
		tui.NewFrame(basic).
			WithTitle(tui.Title{Text: "Default", Pos: tui.TitleTopCenter}),
		tui.NewFrame(styled).
			WithTitle(tui.Title{Text: "Custom", Pos: tui.TitleTopCenter, Style: tui.FrBrightGreen}),
		tui.NewFrame(ranged).
			WithTitle(tui.Title{Text: "Range [0, 100], step=4", Pos: tui.TitleTopCenter, Style: tui.FrYellow}).
			WithTitle(tui.Title{Text: "Custom Label", Pos: tui.TitleBottomCenter, Style: tui.FrBlue}),
		tui.NewFrame(ascii).
			WithTitle(tui.Title{Text: "ASCII preset", Pos: tui.TitleTopCenter, Style: tui.FrBrightMagenta}),
	).WithGap(1)).
		Rounded().
		WithTitle(tui.Title{
			Text:  "TUI Compose v4.0 — Slider",
			Pos:   tui.TitleTopCenter,
			Style: tui.Bold | tui.FrBlue,
		}))

	wnd.Run()
}
