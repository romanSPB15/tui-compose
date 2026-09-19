package extra

import (
	"unicode/utf8"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

// Tree — виджет древовидного списка.
// Узлы задаются плоско, а иерархия определяется полем Depth.
// Отступ между уровнями задаётся полем Offset.
// Добавлено в TUI v3.3.0.
type Tree struct {
	nodes                                    []TreeNode
	Offset                                   int
	DepthEnd, Continue, Horizontal, Vertical rune
}

// TreeNode — один узел дерева.
// Depth = 0 — корень, 1 — первый уровень вложенности и т.д.
// Добавлено в TUI v3.3.0.
type TreeNode struct {
	Label string
	Depth int // уровень вложенности (0, 1, 2...)
	Style tui.Style
}

// NewTree создаёт дерево с узлами по умолчанию (Unicode-символы, Offset = 3).
// Добавлено в TUI v3.3.0.
func NewTree(nodes []TreeNode) *Tree {
	return &Tree{
		nodes:      nodes,
		Offset:     3,
		DepthEnd:   '└',
		Continue:   '├',
		Horizontal: '─',
		Vertical:   '│',
	}
}

func (t *Tree) Render(buf [][]cell.Cell) {
	var verticalCols []int
	for i, node := range t.nodes {
		d := (node.Depth) * t.Offset
		if node.Depth > 0 {
			d -= t.Offset
		} else {
			d = 0
		}

		lastAtDepth := true
		for j := i + 1; j < len(t.nodes); j++ {
			if t.nodes[j].Depth == node.Depth {
				lastAtDepth = false
				break
			}
			if t.nodes[j].Depth < node.Depth {
				break
			}
		}

		if node.Depth > 0 {
			for _, col := range verticalCols {
				if col < d {
					buf[i][col] = cell.Cell{Char: t.Vertical, Style: cell.Style{}}
				}
			}
			lineEndChar := t.DepthEnd
			if !lastAtDepth {
				lineEndChar = t.Continue
			}
			buf[i][d] = cell.Cell{Char: lineEndChar, Style: cell.Style{}}
			for x := d + 1; x < d+t.Offset; x++ {
				buf[i][x] = cell.Cell{Char: t.Horizontal, Style: cell.Style{}}
			}
		}

		style := tui.ConvertToCellStyle(node.Style)
		labelRunes := []rune(node.Label)
		startX := 0
		if node.Depth > 0 {
			startX = d + t.Offset
		}
		for j, r := range labelRunes {
			if startX+j < len(buf[i]) {
				buf[i][startX+j] = cell.Cell{Char: r, Style: style}
			}
		}

		if node.Depth > 0 {
			if !lastAtDepth {
				found := false
				for _, col := range verticalCols {
					if col == d {
						found = true
						break
					}
				}
				if !found {
					verticalCols = append(verticalCols, d)
				}
			} else {
				for idx, col := range verticalCols {
					if col == d {
						verticalCols = append(verticalCols[:idx], verticalCols[idx+1:]...)
						break
					}
				}
			}
		}
	}
}

func (t *Tree) Width() int {
	max := 0
	for _, node := range t.nodes {
		length := utf8.RuneCountInString(node.Label) + node.Depth*t.Offset
		if length > max {
			max = length
		}
	}
	return max
}

func (t *Tree) Height() int {
	return len(t.nodes)
}

// Default устанавливает стандартные Unicode-символы ветвления.
// Добавлено в TUI v3.3.0.
func (t *Tree) Default() *Tree {
	t.DepthEnd = '└'
	t.Continue = '├'
	t.Horizontal = '─'
	t.Vertical = '│'
	return t
}

// Rounded устанавливает скруглённый символ для последнего элемента.
// Добавлено в TUI v3.3.0.
func (t *Tree) Rounded() *Tree {
	t.DepthEnd = '╰'
	t.Continue = '├'
	t.Horizontal = '─'
	t.Vertical = '│'
	return t
}

// Heavy устанавливает жирные Unicode-символы ветвления.
// Добавлено в TUI v3.3.0.
func (t *Tree) Heavy() *Tree {
	t.DepthEnd = '┗'
	t.Continue = '┣'
	t.Horizontal = '━'
	t.Vertical = '┃'
	return t
}

// ASCII устанавливает ASCII-совместимые символы ветвления.
// Добавлено в TUI v3.3.0.
func (t *Tree) ASCII() *Tree {
	t.DepthEnd = '+'
	t.Continue = '+'
	t.Horizontal = '-'
	t.Vertical = '|'
	return t
}

// WithNodes заменяет список узлов дерева.
// Добавлено в TUI v3.3.0.
func (t *Tree) WithNodes(n []TreeNode) *Tree {
	t.nodes = n
	return t
}
