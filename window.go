package tui

import (
	"github.com/charmbracelet/x/term"
)

type window struct{}

func (*window) Width() int {
	w, _, err := term.GetSize(currentApp.f.Fd())
	if err != nil {
		currentApp.LogFatal("tui: window size error")
	}
	return w
}

func (*window) Height() int {
	_, h, err := term.GetSize(currentApp.f.Fd())
	if err != nil {
		currentApp.LogFatal("tui: window size error")
	}
	return h
}
