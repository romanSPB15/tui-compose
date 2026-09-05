package tui

import "github.com/romanSPB15/tui-compose/v4/input"

type Event any

type KeyboardEvent = input.KeyboardEvent

type MouseEvent = input.MouseEvent

type FocusEvent struct {
	Focused bool
}

type CheckFocusableEvent struct {
	Result bool
}
