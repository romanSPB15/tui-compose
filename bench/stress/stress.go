package main

import (
	"os"
	"strconv"
	"sync/atomic"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

// ANSI-палитра: 8 normal + 8 bright Fg.
var stressPalette = []string{
	"30", "31", "32", "33", "34", "35", "36", "37",
}

type StressWidget struct {
	w, h      int
	frame     int
	cells     int
	paletteSz int
}

func NewStressWidget(w, h, cells, paletteSz int) *StressWidget {
	return &StressWidget{w: w, h: h, cells: cells, paletteSz: paletteSz}
}

func (s *StressWidget) Width() int  { return s.w }
func (s *StressWidget) Height() int { return s.h }

func (s *StressWidget) Render(buf [][]cell.Cell) {
	s.frame++
	total := s.w * s.h
	for k := 0; k < s.cells; k++ {
		idx := (s.frame*7919 + k*104729) % total
		y, x := idx/s.w, idx%s.w
		if y >= len(buf) || x >= len(buf[y]) {
			continue
		}
		fg := stressPalette[(s.frame+k)%s.paletteSz]
		buf[y][x] = cell.Cell{
			Char:  rune('A' + (s.frame+k)%26),
			Style: cell.Style{Fg: fg},
		}
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

func main() {
	cells := envInt("STRESS_CELLS", 300)
	paletteSz := envInt("STRESS_PALETTE", 8)

	wnd := tui.NewWindow()
	wnd.SetContent(NewStressWidget(80, 24, cells, paletteSz))

	go func() {
		for {
			select {
			case <-wnd.OnQuit():
				return
			default:
				wnd.DoAndWait(func() {
					wnd.Redraw()
				})
				redraws.Add(1)
			}
		}
	}()

	wnd.Run()
}
