package main

import (
	"math"
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.3 - Sparkline")

	const n = 30
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = int(20 + 15*math.Sin(float64(i)*2*math.Pi/float64(n))) // от 5 до 35
	}

	spark1 := extra.NewSparkline().
		WithValues(data).
		WithHeight(5).
		WithBarStyle(func(i, v int) tui.Style {
			return tui.FrCyan
		})

	spark2 := extra.NewSparkline().
		WithValues(data).
		WithHeight(5).
		WithBarStyle(func(i, v int) tui.Style {
			return tui.FrRed
		})

	data2 := make([]int, 20)
	for i := range data2 {
		data2[i] = i * 2
	}

	spark3 := extra.NewSparkline().
		WithValues(data).
		WithHeight(5).
		WithBarStyle(func(i, v int) tui.Style {
			switch {
			case v > 30:
				return tui.FrRed
			case v > 20:
				return tui.FrYellow
			}
			return tui.FrGreen
		})

	go func() {
		ticker := time.NewTicker(time.Second / 15)
		defer ticker.Stop()
		offset := 0
		for {
			select {
			case <-wnd.OnQuit():
				return
			case <-ticker.C:
				offset++
				newData := make([]int, n)
				for i := 0; i < n; i++ {
					newData[i] = int(20 + 15*math.Sin(float64(i+offset)*2*math.Pi/float64(n)))
				}
				newData2 := make([]int, n)
				for i := 0; i < n; i++ {
					newData2[i] = int(20 + 15*math.Cos(float64(i+offset)*2*math.Pi/float64(n)))
				}
				wnd.Commit(func() {
					spark1.WithValues(newData)
					spark2.WithValues(newData2)
				})
			}
		}
	}()

	wnd.SetContent(tui.NewFrame(
		tui.NewVBox(
			spark1,
			spark2,
			spark3,
		),
	).Rounded().WithTitle(tui.Title{Text: "TUI Compose v3.3 - Sparkline", Pos: tui.TitleTopCenter, Style: tui.Bold | tui.FrBrightMagenta}))
	wnd.Run()
}
