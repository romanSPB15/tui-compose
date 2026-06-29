// Package tui provides framework for creating Text User Interfaces (TUI) in Go.
//
// Библиотека tui-compose позволяет легко создавать интерактивные TUI-приложения.
// Она включает в себя набор готовых компонентов: кнопки, надписи, текстовые поля и другие.
//
// Быстрый старт:
//
// package main

// import "github.com/romanSPB15/tui-compose/v3"

// func main() {
//     wnd := tui.NewWindow()
//     wnd.SetTitle("Моё приложение")

//     label := tui.NewStaticLabel("Привет, TUI!").ColorizeForeground(tui.Cyan)

//     btn := tui.NewButton("Выход", func() {
//         wnd.Quit()
//     })

//     box := tui.NewVBox(label, btn)
//     wnd.SetContent(box)

//	    wnd.Run()
//	}
package tui
