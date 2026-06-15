// Package tui provides framework for creating Text User Interfaces (TUI) in Go.
//
// Библиотека go-tui позволяет легко создавать интерактивные TUI-приложения.
// Она включает в себя набор готовых компонентов: кнопки, надписи,
// прогресс-бар и другие.
//
// Быстрый старт:
//
//	package main
// import (
// 	"github.com/eiannone/keyboard"
// 	"github.com/romanSPB15/go-tui"
// )

//	func main() {
//		wnd := tui.NewWindow()
//		wnd.AddWidgets(tui.NewStaticLabel("Привет, Go!"))
//		btn := tui.NewButton("Нажми ↑", keyboard.KeyArrowUp)
//		btn.OnClick = wnd.Quit
//		wnd.AddWidgets(btn)
//		wnd.Run()
//	}
package tui
