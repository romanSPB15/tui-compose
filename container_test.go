package tui

import (
	"testing"

	"github.com/romanSPB15/tui-compose/v4/cell"
)

// stubWidget — минимальный Widget для тестов контейнеров.
type stubWidget struct {
	w, h int
}

func (s *stubWidget) Render([][]cell.Cell) {}
func (s *stubWidget) Width() int           { return s.w }
func (s *stubWidget) Height() int          { return s.h }

func TestVBoxLayout(t *testing.T) {
	tt := []struct {
		name     string
		children []Widget
		gap      int
		wantW    int
		wantH    int
		wantPos  []Pos
	}{
		{
			name:     "empty",
			children: nil,
			gap:      0,
			wantW:    0,
			wantH:    0,
			wantPos:  []Pos{},
		},
		{
			name:     "single",
			children: []Widget{&stubWidget{w: 5, h: 3}},
			gap:      0,
			wantW:    5,
			wantH:    3,
			wantPos:  []Pos{{0, 0}},
		},
		{
			name: "two no gap",
			children: []Widget{
				&stubWidget{w: 5, h: 2},
				&stubWidget{w: 3, h: 4},
			},
			gap:     0,
			wantW:   5,
			wantH:   6,
			wantPos: []Pos{{0, 0}, {2, 0}},
		},
		{
			name: "two with gap",
			children: []Widget{
				&stubWidget{w: 5, h: 2},
				&stubWidget{w: 3, h: 4},
			},
			gap:     1,
			wantW:   5,
			wantH:   7,
			wantPos: []Pos{{0, 0}, {3, 0}},
		},
		{
			name: "with nil child",
			children: []Widget{
				&stubWidget{w: 5, h: 2},
				nil,
				&stubWidget{w: 3, h: 4},
			},
			gap:     0,
			wantW:   5,
			wantH:   6,
			wantPos: []Pos{{0, 0}, {2, 0}, {2, 0}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			v := NewVBox(tc.children...)
			v.SetGap(tc.gap)

			if got := v.Width(); got != tc.wantW {
				t.Errorf("Width: got %d, want %d", got, tc.wantW)
			}
			if got := v.Height(); got != tc.wantH {
				t.Errorf("Height: got %d, want %d", got, tc.wantH)
			}

			got := v.Child()
			if len(got) != len(tc.children) {
				t.Fatalf("Child len: got %d, want %d", len(got), len(tc.children))
			}
			for i := range tc.wantPos {
				if p := v.Pos(i); p != tc.wantPos[i] {
					t.Errorf("Pos(%d): got %v, want %v", i, p, tc.wantPos[i])
				}
			}
		})
	}
}

func TestHBoxLayout(t *testing.T) {
	tt := []struct {
		name     string
		children []Widget
		gap      int
		wantW    int
		wantH    int
		wantPos  []Pos
	}{
		{
			name:     "empty",
			children: nil,
			gap:      1, // NewHBox ставит gap=1 по умолчанию
			wantW:    0,
			wantH:    0,
			wantPos:  []Pos{},
		},
		{
			name:     "single",
			children: []Widget{&stubWidget{w: 5, h: 3}},
			gap:      1,
			wantW:    5,
			wantH:    3,
			wantPos:  []Pos{{0, 0}},
		},
		{
			name: "two with default gap=1",
			children: []Widget{
				&stubWidget{w: 5, h: 2},
				&stubWidget{w: 3, h: 4},
			},
			gap:     1,
			wantW:   9, // 5 + 1 + 3
			wantH:   4,
			wantPos: []Pos{{0, 0}, {0, 6}},
		},
		{
			name: "two with gap=0",
			children: []Widget{
				&stubWidget{w: 5, h: 2},
				&stubWidget{w: 3, h: 4},
			},
			gap:     0,
			wantW:   8,
			wantH:   4,
			wantPos: []Pos{{0, 0}, {0, 5}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHBox(tc.children...)
			h.SetGap(tc.gap)

			if got := h.Width(); got != tc.wantW {
				t.Errorf("Width: got %d, want %d", got, tc.wantW)
			}
			if got := h.Height(); got != tc.wantH {
				t.Errorf("Height: got %d, want %d", got, tc.wantH)
			}

			for i := range tc.wantPos {
				if p := h.Pos(i); p != tc.wantPos[i] {
					t.Errorf("Pos(%d): got %v, want %v", i, p, tc.wantPos[i])
				}
			}
		})
	}
}

func TestVBoxEmptyWithGap(t *testing.T) {
	v := NewVBox().WithGap(2)
	if got := v.Height(); got != 0 {
		t.Errorf("empty VBox with gap=2: Height got %d, want 0", got)
	}
}
