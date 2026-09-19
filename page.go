package tui

// Page — страница приложения.
// Имеет собственное содержимое и опциональный заголовок.
// Добавлено в TUI 3.0.0.
type Page struct {
	title   string
	content Widget
}

func NewPage(content Widget) *Page {
	return &Page{content: content}
}

// SetTitle устанавливает заголовок страницы и возвращает её.
func (p *Page) SetTitle(title string) *Page {
	p.title = title
	return p
}

// Open открывает страницу в текущем окне.
func (p *Page) Open(w Window) {
	if w == nil {
		return
	}
	if p.title != "" {
		w.SetTitle(p.title)
	}
	w.SetContent(p.content)
	w.Redraw()
}
