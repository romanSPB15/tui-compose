package tea

import (
	"strings"
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/cell"
	"github.com/romanSPB15/tui-compose/v3/input"
)

type Msg any
type Cmd func() Msg

type quitMsg struct{}
type batchMsg struct{ cmds []Cmd }

func Quit() Msg { return quitMsg{} }

func QuitCmd() Cmd {
	return func() Msg { return quitMsg{} }
}

func Batch(cmds ...Cmd) Cmd {
	return func() Msg {
		return batchMsg{cmds: cmds}
	}
}

func After(d time.Duration, fn func(time.Time) Msg) Cmd {
	return func() Msg {
		<-time.After(d)
		return fn(time.Now())
	}
}

type Model interface {
	Init() Cmd
	Update(msg Msg) (Model, Cmd)
	View() tui.Widget
}

type Program struct {
	window tui.Window
	model  Model
	msgCh  chan Msg
}

var currentProgram *Program

func NewProgram(initialModel Model) *Program {
	pr := &Program{
		model: initialModel,
		msgCh: make(chan Msg, 64),
	}
	currentProgram = pr
	return pr
}

func (p *Program) Run() {
	p.window = tui.NewWindow()
	p.window.SetContent(p.model.View())

	p.window.RegisterClickHandler(func(ev *input.MouseEvent) {
		p.Send(*ev)
	})

	p.window.RegisterKeyHandler(func(ev *input.KeyboardEvent) {
		p.Send(*ev)
	})

	go p.eventLoop()
	p.window.Run()
}

func (p *Program) Send(msg Msg) {
	p.msgCh <- msg
}

func (p *Program) eventLoop() {
	if cmd := p.model.Init(); cmd != nil {
		p.msgCh <- cmd()
	}

	for msg := range p.msgCh {
		if batch, ok := msg.(batchMsg); ok {
			for _, cmd := range batch.cmds {
				go func(c Cmd) {
					p.msgCh <- c()
				}(cmd)
			}
			continue
		}

		newModel, cmd := p.model.Update(msg)
		p.model = newModel

		if cmd != nil {
			go func(c Cmd) {
				p.msgCh <- c()
			}(cmd)
		}

		select {
		case <-p.window.OnQuit():
			close(p.msgCh)
			return
		default:
			p.window.Commit(func() {
				idx := p.window.Focus().FocusedIndex()

				p.window.SetContent(p.model.View()) // так как SetContent сбрасывает фокус, нужно его сохранить

				p.window.Focus().SetIndex(idx)
			})

			if _, ok := msg.(quitMsg); ok {
				p.window.Quit()
				close(p.msgCh)
				return
			}
		}
	}
}

type stringView struct {
	wnd  tui.Window
	text string
}

func (sv *stringView) Render(buf [][]cell.Cell) {
	p := cell.ParseMultiline(sv.text)
	copy(buf, p)
}

func (sv *stringView) Width() int {
	max := 0
	for _, line := range strings.Split(sv.text, "\n") {
		if l := len([]rune(line)); l > max {
			max = l
		}
	}
	return max
}

func (sv *stringView) Height() int {
	return strings.Count(sv.text, "\n") + 1
}

func NewStringView(s string) tui.Widget {
	return &stringView{
		text: s,
	}
}

func Send(msg Msg) {
	if currentProgram != nil {
		currentProgram.Send(msg)
	}
}

func (p *Program) ClearFocus() {
	if p.window != nil {
		p.window.Commit(func() {
			p.window.Focus().ClearFocus()
		})
	}
}

func ClearFocusCmd() Cmd {
	return func() Msg {
		if currentProgram != nil {
			currentProgram.ClearFocus()
		}
		return nil
	}
}

func Focus() tui.FocusManager {
	if currentProgram == nil {
		return nil
	}
	return currentProgram.window.Focus()
}
