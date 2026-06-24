package tui

// Page — страница приложения
// Добавлено в TUI 3.0.0
type Page struct {
	title   string
	content Widget
}

func NewPage(content Widget) *Page {
	return &Page{content: content}
}

// SetTitle устанавливает заголовок страницы и возвращает её
func (p *Page) SetTitle(title string) *Page {
	p.title = title
	return p
}

// Open открывает страницу в текущем окне
func (p *Page) Open() {
	if currentWindow == nil {
		return
	}
	if p.title != "" {
		currentWindow.SetTitle(p.title)
	}
	currentWindow.SetContent(p.content)
	currentWindow.Redraw()
}
