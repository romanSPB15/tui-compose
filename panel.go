package tui

type RightPanel struct {
	content Widget
	s       int
}

func NewRightPanel(content Widget) *RightPanel {
	return &RightPanel{content: content}
}

func (rp *RightPanel) MaxWidth() int {
	return CurrentWindow().Width()
}

func (rp *RightPanel) MaxHeight() int {
	return CurrentWindow().Height()
}

func (rp *RightPanel) InnerText() string {
	return ""
}

func (rp *RightPanel) Child() []Widget {
	return []Widget{rp.content}
}

func (rp *RightPanel) Pos(idx int) Pos {
	return Pos{
		Line: 0,
		Col:  rp.MaxWidth() - rp.panelWidth(),
	}
}

func (rp *RightPanel) panelWidth() int {
	mw := rp.content.MaxWidth()
	if rp.s > mw {
		return rp.s
	}
	return mw
}

func (rp *RightPanel) SetPanelWidth(w int) *RightPanel {
	rp.s = w
	return rp
}
