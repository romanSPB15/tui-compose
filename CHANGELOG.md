# Changelog

## [3.0.0] – 2026-06-29

### Added
- Контейнеры `HBox` и `VBox` с автоматическим расчётом позиций (`Container`).
- Виджет `InputField` – однострочное поле ввода с курсором.
- Поддержка мыши: SGR (1006), события кликов, перетаскивания.
- Виджеты `Canvas` (16 цветов) и `CanvasRGB` (True Color) с двумя режимами отображения (`PixelTwoSymbol`, `PixelOneSymbol`).
- Методы `Width()` и `Height()` для `Canvas` и `CanvasRGB`.
- Методы `SetContent`, `SetOverlay`, `ShowOverlay`, `HideOverlay` для управления содержимым окна.
- Структура `Page` с методом `Open()`.
- Методы `SetTitle`, `CopyToClipboard`, `DisableFocusChange`.
- Интерфейсы `TextInput`, `Container`, `ClickableAt`.
- Собственный парсер ANSI-последовательностей (`parseAnsiKeyboardInput`).
- `RegisterKeyHandler` – регистрация обработчиков клавиш.
- `RegisterClickHandler` – глобальные обработчики мыши.
- Методы `Do` и `DoAndWait` для выполнения задач в UI-потоке.

### Changed
- Удалена зависимость от `github.com/eiannone/keyboard` – теперь используется собственный парсер ввода.
- Система фокуса переключения по Tab/Shift+Tab.
- Перерисовка оптимизирована (перерисовываются только изменённые строки).
- Вместо глобального `App` используется `Window`.

### Removed
- Зависимость от `github.com/eiannone/keyboard`.
- Методы `AddWidgets`, `Clear` – заменены на `SetContent`.
- Старые примеры (будут возвращены в следующих версиях).

### Fixed
- Исправлены ошибки фокуса, перерисовки и обработки событий.
- Исправлено определение размеров терминала.
- Исправлен парсинг клавиш в контейнерах.

### Dependencies
- Единственная внешняя зависимость – `golang.org/x/term`.