package extra

import (
	"strings"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/builder"
)

type Tree struct {
	nodes                                    []TreeNode
	Offset                                   int
	DepthEnd, Continue, Horizontal, Vertical rune
}

type TreeNode struct {
	Label string
	Depth int // уровень вложенности (0, 1, 2...)
	Style tui.Style
}

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

func (t *Tree) InnerText() string {
	var sbRes builder.Builder
	lastDepth := 0
	lines := make(map[int]bool)
	for i, node := range t.nodes {
		d := (node.Depth)*t.Offset - t.Offset
		d = max(d, 0)

		var sb builder.Builder

		indent := strings.Repeat(" ", d)
		sb.WriteString(indent)

		// чтобы выбрать └ или ├
		lastWithThisDepth := func() bool {
			for j := i + 1; j < len(t.nodes); j++ {
				if t.nodes[j].Depth == node.Depth {
					return false
				}
				if t.nodes[j].Depth < node.Depth {
					return true
				}
			}
			return true
		}

		end := func() {
			sb.WriteRune(t.DepthEnd)
			for range t.Offset - 1 {
				sb.WriteRune(t.Horizontal)
			}
		}

		nonend := func() {
			sb.WriteRune(t.Continue)
			for range t.Offset - 1 {
				sb.WriteRune(t.Horizontal)
			}
		}

		if node.Depth > 0 {
			last := lastWithThisDepth()
			if lastDepth != node.Depth && !lines[d] {
				if last {
					end()
					lines[d] = false
				} else {
					nonend()
				}
				lines[d] = true
			} else {
				if last {
					end()
					lines[d] = false
				} else {
					nonend()
				}
			}
		}

		style := node.Style.String()
		if style != "" {
			sb.WriteString(style)
		}

		sb.WriteString(node.Label)
		if style != "" {
			sb.WriteString("\033[0m")
		}
		sb.WriteString("\n")
		lastDepth = node.Depth

		s := []rune(sb.String())

		for k, v := range lines {
			if v && k != d {
				s[k] = t.Vertical
			}
		}
		sbRes.WriteString(string(s))

	}
	return sbRes.String()
}

func (t *Tree) MaxWidth() int {
	max := 0
	for _, node := range t.nodes {
		length := len(node.Label) + node.Depth*t.Offset
		if length > max {
			max = length
		}
	}
	return max
}

func (t *Tree) MaxHeight() int {
	return len(t.nodes)
}

func (t *Tree) Default() *Tree {
	t.DepthEnd = '└'
	t.Continue = '├'
	t.Horizontal = '─'
	t.Vertical = '│'
	return t
}

func (t *Tree) Rounded() *Tree {
	t.DepthEnd = '╰'
	t.Continue = '├'
	t.Horizontal = '─'
	t.Vertical = '│'
	return t
}

func (t *Tree) Heavy() *Tree {
	t.DepthEnd = '┗'
	t.Continue = '┣'
	t.Horizontal = '━'
	t.Vertical = '┃'
	return t
}

func (t *Tree) ASCII() *Tree {
	t.DepthEnd = '+'
	t.Continue = '+'
	t.Horizontal = '-'
	t.Vertical = '|'
	return t
}

func (t *Tree) WithNodes(n []TreeNode) *Tree {
	t.nodes = n
	return t
}
