package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v3"
)

var tcellPalette = []tcell.Color{
	tcell.ColorBlack, tcell.ColorRed, tcell.ColorGreen, tcell.ColorYellow,
	tcell.ColorBlue, tcell.ColorPurple, tcell.ColorTeal, tcell.ColorWhite,
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

func main() {
	cells := envInt("STRESS_CELLS", 300)
	paletteSz := envInt("STRESS_PALETTE", 8)
	if paletteSz > len(tcellPalette) {
		paletteSz = len(tcellPalette)
	}

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("NewScreen: %v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("Init: %v", err)
	}
	defer s.Fini()

	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	s.Clear()
	s.Show()

	logFile, _ := os.Create("tcell_fps.log")
	defer logFile.Close()

	// Канал для сигнала "пора рендерить следующий кадр".
	// Буфер 1: продюсер не накапливает лишние запросы.
	renderCh := make(chan struct{}, 1)
	stopCh := make(chan struct{})

	// Продюсер: просит рендерить как можно быстрее.
	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				select {
				case renderCh <- struct{}{}:
				case <-stopCh:
					return
				}
			}
		}
	}()

	// Тикер FPS — отдельная горутина, только читает atomic.
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

	// ВСЁ в главной горутине: и события, и рендер.
	eventQ := s.EventQ()
	frame := 0
	total := 80 * 24

	for {
		select {
		case ev := <-eventQ:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Str() == "q" {
					close(stopCh)
					return
				}
			case *tcell.EventResize:
				s.Sync()
			}

		case <-renderCh:
			frame++
			for k := 0; k < cells; k++ {
				idx := (frame*7919 + k*104729) % total
				x, y := idx%80, idx/80
				fg := tcellPalette[(frame+k)%paletteSz]
				s.SetContent(x, y, rune('A'+(frame+k)%26), nil,
					tcell.StyleDefault.Foreground(fg))
			}
			s.Show()
			redraws.Add(1)
		}
	}
}
