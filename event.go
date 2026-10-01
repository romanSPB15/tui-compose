package tui

import "github.com/romanSPB15/acell"

type Event any

type FocusEvent struct {
	Focused bool
}

type CheckFocusableEvent struct {
	Result bool
}

type MouseHoverEvent struct {
	Entered bool
	Pos     acell.Point
}

type WindowEvent struct {
	Window Window
}
