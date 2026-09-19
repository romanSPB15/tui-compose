package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

// 8 базовых цветов — совпадает с первыми 8 у tui-compose.
var gotuiPalette = []ui.Color{
	ui.ColorBlack, ui.ColorRed, ui.ColorGreen, ui.ColorYellow,
	ui.ColorBlue, ui.ColorMagenta, ui.ColorCyan, ui.ColorWhite,
}

type StressWidget struct {
	*widgets.Paragraph
	frame     int
	w, h      int
	cells     int
	paletteSz int
}

func NewStressWidget(w, h, cells, paletteSz int) *StressWidget {
	b := widgets.NewParagraph()
	b.Border = false
	b.SetRect(0, 0, w, h)
	return &StressWidget{Paragraph: b, w: w, h: h, cells: cells, paletteSz: paletteSz}
}

func (s *StressWidget) Draw(buf *ui.Buffer) {
	s.frame++
	total := s.w * s.h
	for k := 0; k < s.cells; k++ {
		idx := (s.frame*7919 + k*104729) % total
		y, x := idx/s.w, idx%s.w
		fg := gotuiPalette[(s.frame+k)%s.paletteSz]
		buf.SetCell(
			ui.Cell{
				Rune:  rune('A' + (s.frame+k)%26),
				Style: ui.NewStyle(fg),
			},
			image.Pt(x, y),
		)
	}
}

var redraws atomic.Int64

func envInt(name string, def int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

type renderMsg struct{ done chan struct{} }

func main() {
	cells := envInt("STRESS_CELLS", 300)
	paletteSz := envInt("STRESS_PALETTE", 8)
	if paletteSz > len(gotuiPalette) {
		paletteSz = len(gotuiPalette)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("init: %v", err)
	}
	defer ui.Close()

	stress := NewStressWidget(80, 24, cells, paletteSz)
	ui.Render(stress)

	logFile, _ := os.Create("gotui_fps.log")
	defer logFile.Close()

	renderCh := make(chan renderMsg, 1)
	stopCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				msg := renderMsg{done: make(chan struct{})}
				select {
				case renderCh <- msg:
					<-msg.done
					redraws.Add(1)
				case <-stopCh:
					return
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		var last int64
		for range ticker.C {
			now := redraws.Load()
			fmt.Fprintf(logFile, "FPS: %d (cells=%d palette=%d)\n", now-last, cells, paletteSz)
			last = now
		}
	}()

	events := ui.PollEvents()

	for {
		select {
		case e := <-events:
			if e.ID == "q" || e.ID == "<C-c>" {
				close(stopCh)
				return
			}
		case msg := <-renderCh:
			ui.Render(stress)
			close(msg.done)
		}
	}
}
