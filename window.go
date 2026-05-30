package tui

import (
	"log"

	"github.com/charmbracelet/x/term"
)

type window struct{}

func (*window) Width() int {
	w, _, err := term.GetSize(currentApp.f.Fd())
	if err != nil {
		log.Fatal("tui: window size error")
	}
	return w
}

func (*window) Height() int {
	_, h, err := term.GetSize(currentApp.f.Fd())
	if err != nil {
		log.Fatal("tui: window size error")
	}
	return h
}
